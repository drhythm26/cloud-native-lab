resource "google_container_cluster" "primary" {
  name                     = var.cluster_name
  location                 = var.zone
  network                  = google_compute_network.vpc.id
  subnetwork               = google_compute_subnetwork.subnet.id
  remove_default_node_pool = true
  initial_node_count       = 1
  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }
  workload_identity_config {
    workload_pool = "${var.project}.svc.id.goog"
  }
  deletion_protection = false
  depends_on = [
    google_project_service.container,
    google_compute_subnetwork.subnet,
  ]
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.cluster_name}-node-pool"
  location   = var.zone
  cluster    = google_container_cluster.primary.name
  node_count = var.node_count
  node_config {
    machine_type    = var.machine_type
    service_account = google_service_account.gke_node.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
    labels = {
      environment = var.environment
    }
  }
  management {
    auto_repair  = true
    auto_upgrade = true
  }
}