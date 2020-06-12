package test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testDataSourceConfig_basic = `
data "test_scenario" "test" {
  step {
    program = ["echo", "{\"foo\": \"bar\"}"]
    expect  = {"foo" = "bar"}
  }
}

output "result" {
  value = "${data.test_scenario.test.result}"
}
`

func TestDataSource_basic(t *testing.T) {
//	programPath, err := buildDataSourceTestProgram()
//	if err != nil {
//		t.Fatal(err)
//		return
//	}

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfig_basic,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.test_scenario.test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if outputs["result"] == nil {
						return fmt.Errorf("missing 'result' output")
					}

					//if outputs["result"].Value != [true]  {
					//	return fmt.Errorf(
					//		"result is false, want true",
					//		outputs["result"].Value,
					//	)
					//}
					return nil
				},
			},
		},
	})
}

const testDataSourceConfig_error = `
data "test_scenario" "test" {
  step {
    program = ["%s"]

    query = {
      fail = "true"
    }
  }
}
`

func TestDataSource_error(t *testing.T) {
	programPath, err := buildDataSourceTestProgram()
	if err != nil {
		t.Fatal(err)
		return
	}

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmt.Sprintf(testDataSourceConfig_error, programPath),
				ExpectError: regexp.MustCompile("I was asked to fail"),
			},
		},
	})
}

func buildDataSourceTestProgram() (string, error) {
	// We have a simple Go program that we use as a stub for testing.
	cmd := exec.Command(
		"go", "install",
		"github.com/chrismatteson/terraform-provider-test/test/test-check/tf-test-data-source",
	)
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("failed to build test stub program: %s", err)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME") + "/go")
	}

	programPath := path.Join(
		filepath.SplitList(gopath)[0], "bin", "tf-test-data-source",
	)
	return programPath, nil
}
