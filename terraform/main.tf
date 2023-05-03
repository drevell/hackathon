module "cloud_run" {
  source                    = "git::https://github.com/abcxyz/terraform-modules.git//modules/cloud_run?ref=SHA_OR_TAG"
  project_id                = "my-project-id"
  name                      = "project-name"
  secrets                   = ["app-secret", "app-file-secret"]
  service_account_email = google_service_account.run_service_account.email
  envvars = {
    "PROJECT_ID" : "my-project-id",
  }

  secret_envvars = {
    "APP_SECRET" : {
      name : "app-secret",
      version : "latest",
    }
  }

  secret_volumes = {
    "/var/secrets" : {
      name : "app-file-secret",
      version : "latest",
    }
  }

  invokers = ["user:test-account-group@google.com"]
  developers = [module.github_ci.service_account_email]
}
