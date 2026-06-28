resource "google_container_cluster" "primary" {
  name                     = "cloud-native-lab-gke"
  location                 = var.zone
  remove_default_node_pool = true
  initial_node_count       = 1
  network                  = google_compute_network.vpc.name
  subnetwork               = google_compute_subnetwork.subnet.name
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }
  depends_on = [
    google_project_service.container,
    google_compute_subnetwork.subnet,
  ]
}

resource "google_container_node_pool" "primary" {
  name               = "cloud-native-lab-node-pool"
  location           = var.zone
  cluster            = google_container_cluster.primary.name
  initial_node_count = 1
  node_config {
    machine_type = "e2-standard-2"
    disk_size_gb = 30
    disk_type    = "pd-standard"
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
  depends_on = [
    google_project_service.container,
    google_container_cluster.primary,
  ]
}