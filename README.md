# terraform-provider-test
A Terraform provider for running arbitrary tests. Intended to allow tests to be written in the native HCL2 language and run as part of a pipeline.

This provider is currently based on the terraform-provider-external, which is probably a good reference point.

The provider takes a number of step blocks and each step block accepts:
* program: List of strings of program to run. Required.
* working_dir: Working directory to run in. Defaults to current directory.
* query: Optional query to pass to program.
* expect: Map of expected results from program.

The program is expected to return JSON. The provider converts this to map[string]string and back to JSON, and converts expect input to JSON as well. This eliminates any issues with extra spacing or new line characters in the comparision. If the output of the program matches the expect, a pass is given, otherwise a fail. Those are added to a list in the order of the steps. Additionally the raw results of of the program's json output is given in the same order to raw_result.

This can definitely be cleaned up. Probably should move result and raw result back under each step. Having a single combined PASS/FAIL for the full scenario is probably the only thing needed at the top level.

Multiple tests can be run in each scenario resource, or multiple scenario resources can be described. It likely makes sense to group very similiar tests in the same resource along with any prerequisite steps that prepare for those steps (prerequisate steps still have an expect, so they need to return something success message for them being ran. Could be switched to optional to not check if that made sense).

Unlike traditional CI testing solutions, this doesn't have different environment variables to pass in and then run tests again. The idea is that in each pipeline (say a TFC workspace) would be created for each infrastructure scenario and the variables passed in a traditional way to terraform code. Then the tests run against each of those workspaces.

