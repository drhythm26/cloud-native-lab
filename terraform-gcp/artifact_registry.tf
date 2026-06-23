resource "google_artifact_registry_repository" "docker" {
  location      = var.region
  repository_id = var.repo_id
  format        = "DOCKER"
  description   = "Docker images for Release Tracker"
  depends_on    = [google_project_service.artifactregistry]
}
