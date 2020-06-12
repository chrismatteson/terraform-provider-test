package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


type Step struct {
	Steps []*Step
}

type StepAttr struct {
	program            []string
	working_dir        string
	query              map[string][]string
	expect            string
}


func dataSourceTestScenario() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTestScenarioRead,

		Schema: map[string]*schema.Schema{
			"step": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"program": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"working_dir": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},

						"query": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},


						"expect": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
		},
	}
}

func dataSourceTestScenarioRead(d *schema.ResourceData, meta interface{}) error {
	if rawSteps, hasRawSteps := d.GetOk("step"); hasRawSteps {
		var rawStepIntfs = rawSteps.([]interface{})
//		steps := make([]*StepAttr, len(rawStepIntfs))
		steps := make([]bool, len(rawStepIntfs))

		for i, stepI := range rawStepIntfs {
			rawStep := stepI.(map[string]interface{})
			programI := rawStep["program"].([]interface{})
			workingDir := rawStep["working_dir"].(string)
			query := rawStep["query"].(map[string]interface{})
			expect := rawStep["expect"].(map[string]string)

			// This would be a ValidateFunc if helper/schema allowed these
			// to be applied to lists.
			if err := validateProgramAttr(programI); err != nil {
				return err
			}

			program := make([]string, len(programI))
			for i, vI := range programI {
				program[i] = vI.(string)
			}

			cmd := exec.Command(program[0], program[1:]...)

			cmd.Dir = workingDir

			queryJson, err := json.Marshal(query)
			if err != nil {
				// Should never happen, since we know query will always be a map
				// from string to string, as guaranteed by d.Get and our schema.
				return err
			}

			cmd.Stdin = bytes.NewReader(queryJson)

			resultJson, err := cmd.Output()
			log.Printf("[TRACE] JSON output: %+v\n", resultJson)
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.Stderr != nil && len(exitErr.Stderr) > 0 {
						return fmt.Errorf("failed to execute %q: %s", program[0], string(exitErr.Stderr))
					}
					return fmt.Errorf("command %q failed with no error message", program[0])
				} else {
					return fmt.Errorf("failed to execute %q: %s", program[0], err)
				}
			}

			result := map[string]string{}
			err = json.Unmarshal(resultJson, &result)
			if err != nil {
				return fmt.Errorf("command %q produced invalid JSON: %s", program[0], err)
			}

			steps[i] = reflect.DeepEqual(result, expect)
//			if equal {
//				steps[i] = true
//			} else {
//				steps[i] = false
//			}

			return nil
		}
		d.Set("result", steps)

		d.SetId("-")
	}
	return nil
}
