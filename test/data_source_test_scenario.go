package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
//	"reflect"
//	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
							Required: true,
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
					Type: schema.TypeString,
				},
			},
			"raw_result": {
				Type:	schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
		},
	}
}

func dataSourceTestScenarioRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] Begin dataSourceTestScenarioRead")
	if rawSteps, hasRawSteps := d.GetOk("step"); hasRawSteps {
		log.Printf("[TRACE] hasRawSteps true")
		var rawStepIntfs = rawSteps.([]interface{})
		steps := make([]string, len(rawStepIntfs))
		results := make([]map[string]string, len(rawStepIntfs))

		for i, stepI := range rawStepIntfs {
			rawStep := stepI.(map[string]interface{})
			programI := rawStep["program"].([]interface{})
			workingDir := rawStep["working_dir"].(string)
			query := rawStep["query"].(map[string]interface{})
			expect := rawStep["expect"].(map[string]interface{})

			mapStringExpect := make(map[string]string)

			for key, value := range expect {
			        strValue := fmt.Sprintf("%v", value)
				mapStringExpect[key] = strValue
			}

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

			expectString, _ := json.Marshal(mapStringExpect)
			resultString, _ := json.Marshal(result)
                        log.Printf("[TRACE] JSON expectString output: %+v\n", expectString)
                        log.Printf("[TRACE] JSON resultString output: %+v\n", resultString)

			if Equal(expectString, resultString) {
				steps[i] = "PASS"
                	        log.Printf("[TRACE] expectString and resultString are equal")
			} else {
				steps[i] = "FAIL"
                        	log.Printf("[TRACE] expectString and resultString are not equal")
			}
			results[i] = result

//			return nil
		}
		log.Println("[TRACE] steps", steps)
		d.Set("result", steps)
		log.Println("[TRACE] raw_results", results)
		d.Set("raw_result", results)

		d.SetId("-")
	}
	return nil
}
