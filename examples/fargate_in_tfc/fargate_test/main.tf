data terraform_remote_state "fargate" {
  backend = "remote"

  config = {
    organization = var.organization
    workspaces = {
      name = var.workspace
    }
  }
}

data test_scenario "fargate" {
  step {
    program = ["curl", "-o", "/dev/null", "-s", "-w", "{\"response_code\":\"%%{http_code}\"}", "http://54.241.103.180/"]
    expect = { "response_code" = "200" }
  }
}

resource null_resource "trigger" {
  triggers = {
    trigger = join(",", data.test_scenario.fargate.result)
  }
}
