# if role is different, it is not duplicated
resource "google_project_iam_binding" "binding_4" {
  project = var.project
  role    = "roles/storage.objectViewer"
  members = [
    "serviceAccount:dummy4@example.com",
  ]
}

