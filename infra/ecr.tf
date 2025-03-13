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

resource "aws_ecr_repository" "three-point-cloud-career-stats-container-repository" {
  name                 = var.container_repo_name_career_stats
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

variable "container_repo_name_career_stats" {
  type = string
  default = "three-point-cloud-career-stats-container-repository"
} 

resource "aws_ecr_repository" "three-point-cloud-career-stats-coordinator-container-repository" {
  name                 = var.container_repo_name_career_stats_coordinator
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}

variable "container_repo_name_career_stats_coordinator" {
  type = string
  default = "three-point-cloud-career-stats-coordinator-container-repository"
} 