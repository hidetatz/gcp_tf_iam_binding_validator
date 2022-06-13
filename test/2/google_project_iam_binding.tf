# binding_1 is duplicated with binding_2 in google_project_iam_binding_2.tf
resource "google_project_iam_binding" "binding_1" {
  project = var.project
  role    = "roles/storage.admin"
  members = [
    "serviceAccount:dummy@example.com",
  ]
}

