resource "google_service_account" "gke_node" {
  account_id   = "gke-node-sa"
  display_name = "GKE Node Service Account"
  depends_on   = [google_project_service.iam]
}

resource "google_project_iam_member" "gke_node_log_write" {
  project = var.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.gke_node.email}"
}

resource "google_project_iam_member" "gke_node_metrics_write" {
  project = var.project
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.gke_node.email}"
}

resource "google_project_iam_member" "gke_node_monitoring_viewer" {
  project = var.project
  role    = "roles/monitoring.viewer"
  member  = "serviceAccount:${google_service_account.gke_node.email}"
}

resource "google_artifact_registry_repository_iam_member" "gke_node_reader" {
  location   = google_artifact_registry_repository.docker.location
  repository = google_artifact_registry_repository.docker.name
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${google_service_account.gke_node.email}"
}