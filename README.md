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

## Using the authentication server

Users will create their account on this website.
Directing users to auth.fabianoden.com/login?redirect=yoursite.com.
When they login they will be redirected to your site with the params token, expire_seconds and id.
You can then make an GET call to auth.fabianoden.com/validate_token with params id and token.
It will return json{valid:boolean, username: string}, valid determening if the user has correct token.
