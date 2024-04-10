resource "kubernetes_deployment" "backend" {
  metadata {
    name = "backend"
    labels = {
      "io.kompose.service" = "backend"
    }
    annotations = {
      "kompose.cmd"         = "kompose convert -o ./kubernetes/"
      "kompose.version"     = "1.32.0 (HEAD)"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "io.kompose.service" = "backend"
      }
    }

    template {
      metadata {
        labels = {
          "io.kompose.network/tic-tac-toe-default" = "true"
          "io.kompose.service"                     = "backend"
        }
        annotations = {
          "kompose.cmd"         = "kompose convert -o ./kubernetes/"
          "kompose.version"     = "1.32.0 (HEAD)"
        }
      }

      spec {
        container {
          image = "reidelkins/tic-tac-toe-backend"
          name  = "backend"

          port {
            container_port = 8080
            host_port      = 8080
            protocol       = "TCP"
          }
        }

        restart_policy = "Always"
      }
    }
  }
}

resource "kubernetes_deployment" "frontend" {
  metadata {
    name = "frontend"
    labels = {
      "io.kompose.service" = "frontend"
    }
    annotations = {
      "kompose.cmd"         = "kompose convert -o ./kubernetes/"
      "kompose.version"     = "1.32.0 (HEAD)"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "io.kompose.service" = "frontend"
      }
    }

    template {
      metadata {
        labels = {
          "io.kompose.network/tic-tac-toe-default" = "true"
          "io.kompose.service"                     = "frontend"
        }
        annotations = {
          "kompose.cmd"         = "kompose convert -o ./kubernetes/"
          "kompose.version"     = "1.32.0 (HEAD)"
        }
      }

      spec {
        container {
          image = "reidelkins/tic-tac-toe-frontend"
          name  = "frontend"

          port {
            container_port = 80
            host_port      = 80
            protocol       = "TCP"
          }
        }

        restart_policy = "Always"
      }
    }
  }
}

resource "kubernetes_config_map" "init_sql" {
  metadata {
    name = "init-sql"
  }

  data = {
    "init.sql" = file("/Users/reidelkins/Tic-Tac-Toe/init.sql")
  }
}

resource "kubernetes_deployment" "db" {
  metadata {
    name = "db"
    labels = {
      "io.kompose.service" = "db"
    }
    annotations = {
      "kompose.cmd"     = "kompose convert -o ./kubernetes/"
      "kompose.version" = "1.32.0 (HEAD)"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        "io.kompose.service" = "db"
      }
    }

    strategy {
      type = "Recreate"
    }

    template {
      metadata {
        labels = {
          "io.kompose.network/tic-tac-toe-default" = "true"
          "io.kompose.service"                     = "db"
        }
        annotations = {
          "kompose.cmd"     = "kompose convert -o ./kubernetes/"
          "kompose.version" = "1.32.0 (HEAD)"
        }
      }

      spec {
        container {
          name  = "db"
          image = "postgres:latest"

          env {
            name  = "POSTGRES_DB"
            value = "tic-tac-toe"
          }

          env {
            name  = "POSTGRES_PASSWORD"
            value = "password"
          }

          env {
            name  = "POSTGRES_USER"
            value = "user"
          }

          port {
            container_port = 5432
            host_port      = 5432
            protocol       = "TCP"
          }

          volume_mount {
            name       = "init-sql"
            mount_path = "/docker-entrypoint-initdb.d/init.sql"
            sub_path   = "init.sql"
          }
        }

        restart_policy = "Always"

        volume {
          name = "init-sql"
          config_map {
            name = "init-sql"
          }
        }
      }
    }
  }
}