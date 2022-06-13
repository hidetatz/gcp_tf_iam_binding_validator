# binding_4 is duplicated with binding_5 in google_project_iam_binding_5.tf
resource "google_project_iam_binding" "binding_4" {
  project = var.project
  role    = "roles/storage.admin"
  members = [
    "serviceAccount:dummy@example.com",
  ]

  condition {
    title       = "expires_after_2019_12_31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}

