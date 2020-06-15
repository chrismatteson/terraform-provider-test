resource "tfe_workspace" "fargate" {
  name              = "fargate"
  organization      = var.organization
  working_directory = "examples/fargate_in_tfc/fargate"

  vcs_repo {
    identifier     = "chrismatteson/terraform-provider-test"
    oauth_token_id = var.oauth_token_id
  }
}

resource "tfe_variable" "aws_access_key_id" {
  key          = "AWS_ACCESS_KEY_ID"
  value        = var.aws_access_key_id
  category     = "env"
  workspace_id = tfe_workspace.fargate.id
}

resource "tfe_variable" "aws_secret_access_key" {
  key          = "AWS_SECRET_ACCESS_KEY"
  value        = var.aws_secret_access_key
  category     = "env"
  workspace_id = tfe_workspace.fargate.id
}

resource "tfe_workspace" "fargate_test" {
  name              = "fargate_test"
  organization      = var.organization
  working_directory = "examples/fargate_in_tfc/fargate_test"

  vcs_repo {
    identifier     = "chrismatteson/terraform-provider-test"
    oauth_token_id = var.oauth_token_id
  }
}

resource "tfe_variable" "organization" {
  key          = "organization"
  value        = var.organization
  category     = "terraform"
  workspace_id = tfe_workspace.fargate_test.id
}

resource "tfe_variable" "workspace" {
  key          = "workspace"
  value        = tfe_workspace.fargate.id
  category     = "terraform"
  workspace_id = tfe_workspace.fargate_test.id
}
