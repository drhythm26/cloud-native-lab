output "workload_identity_provider" {
  description = "Workload Identity Provider name for GitHub Actions"
  value       = google_iam_workload_identity_pool_provider.github.name
}

output "cicd_sa_email" {
  description = "CI/CD service account email"
  value       = google_service_account.cicd.email
}