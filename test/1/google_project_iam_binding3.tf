# if role is different, it is not duplicated
resource "google_project_iam_binding" "binding_3" {
  project = var.project
  role    = "roles/storage.objectViewer"
  members = [
    "serviceAccount:dummy3@example.com",
  ]
}

