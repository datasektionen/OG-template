job "og-template" {
  namespace = "default"

  type = "service"

  group "og-template" {
    network {
      port "http" {}
    }

    service {
      name     = "og-template"
      port     = "http"
      provider = "nomad"
      tags = [
        "traefik.enable=true",
        "traefik.http.routers.og-template.rule=Host(`og-template.datasektionen.se`)",
        "traefik.http.routers.og-template.tls.certresolver=default",
      ]
    }

    task "og-template" {
      driver = "docker"

      config {
        image = var.image_tag
        ports = ["http"]
      }

      template {
        data        = <<ENV
SERVER_PORT={{ env "NOMAD_PORT_http" }}
SERVER_URL=https://og-template.datasektionen.se
{{ with nomadVar "nomad/jobs/og-template" }}
DB_CONN=postgres://og-template:{{ .db_password }}@postgres.dsekt.internal:5432/og_template?sslmode=disable
OIDC_CLIENT_SECRET={{ .oidc_client_secret }}
APP_SECRET_KEY={{ .app_secret_key }}
{{ end }}
OIDC_PROVIDER=https://sso.datasektionen.se/op
OIDC_CLIENT_ID=og-template
OIDC_REDIRECT_URL=https://og-template.datasektionen.se/auth/callback
ENV
        destination = "local/.env"
        env         = true
      }
    }
  }
}

variable "image_tag" {
  type    = string
  default = "ghcr.io/datasektionen/og-template:latest"
}
