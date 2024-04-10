provider "google" {
  project = "tic-tac-toe-kubernet-deploy"
  region  = "us-east1"
}

resource "google_container_cluster" "cluster" {
  name               = "my-cluster"
  location           = "your-region"
  initial_node_count = 1
  # Add other cluster configuration options as needed
}