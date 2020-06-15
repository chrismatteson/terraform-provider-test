resource "tfe_workspace" "fargate" {
  name         = "fargate"
  organization = var.organization
  working_directory = "examples/fargate_in_tfc/fargate"

  vcs_repo {
    identifier = "chrismatteson/terraform-provider-test"
    oauth_token_id = var.oauth_token_id
  }
}

resource "tfe_workspace" "fargate_test" {
  name         = "fargate_test"
  organization = var.organization
  working_directory = "examples/fargate_in_tfc/fargate_test"

  vcs_repo {
    identifier = "chrismatteson/terraform-provider-test"
    oauth_token_id = var.oauth_token_id
  }
}
