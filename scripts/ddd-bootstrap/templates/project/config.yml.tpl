appName: {{.ProjectNameYAML}}
httpPort: 8080
env: "dev"

redis:
  default:
    host: localhost
    port: 6379
    password: ""
    db: 0
    tracing: true

postgres:
  default:
    host: {{.DBHostYAML}}
    port: {{.DBPort}}
    user: {{.DBUserYAML}}
    password: {{.DBPasswordYAML}}
    sslMode: {{.DBSSLModeYAML}}
    db: {{.DBNameYAML}}
    tracing: true
