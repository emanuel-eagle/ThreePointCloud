resource "aws_ecr_repository" "three-point-cloud-player-list-container-repository" {
  name                 = var.container_repo_name
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

variable "container_repo_name" {
  type = string
  default = "three-point-cloud-player-list-container-repository"
} 