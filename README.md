# Auth-fabian

Authorization service with golang.

## Development

### Install dependencies

`cd /src && go install`

#### Enviroment variables

.env_secrets files should contain:

- EMAIL_LOGIN, EMAIL_PASSWORD for the gmail account
- EMAIL_PUBLIC for the mail account which the recipient sees

### Start

`go run src/main.go`

#### Start with docker

`docker compose up -d`
