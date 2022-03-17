resource "google_project_iam_binding" "binding_2" {
  project = var.project
  role    = "roles/storage.editor"
  members = [
    "serviceAccount:dummy2@example.com",
  ]
}
