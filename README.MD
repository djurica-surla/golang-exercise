
## Description

The following requirements have been met.

- Postgres database
- JWT authentication
- Rest endpoints
- App config and linting config
- Unit/Integration tests
- Dockerized production environment
- Dockerized environment for integration tests

## Instructions

- Clone the repository
- Run go mod tidy to download dependencies
- **make start** is used to run the production docker compose file. Config file is setup to connect to this environment. App will listen on localhost:8080 on your machine. If you want to change the port than change it in the production docker compose file as well.
- **make integration** is used to run environment which will run integration tests with postgres. Credentials for this test environment are hardcoded and don't need to be setup.
- **make run** will run the app locally. However you will need to provide postgres credentials for the local environment.
- **make lint-fix** will run the linter and fix linting as defined in .golangci.yaml


## Endpoints
Database migrations will contain insert migration to add one mock company. You can use the ID provided bellow to get/update/delete that record.

You will need to call login with some randome username in json body, after which the app will return http only cookie that will last for 5 minutes and will give you authorization to perform crud operations.

POST - localhost:8080/login
{
    "username": "test"
}

GET - localhost:8080/companies/c03887be-e2ee-11ed-b5ea-0242ac120002

POST - localhost:8080/companies
{
    "name": "Google",
    "description": "Google description",
    "employees": 11230,
    "registered": true,
    "type": "Corporation"
}

PATCH - localhost:8080/companies/c03887be-e2ee-11ed-b5ea-0242ac120002
{
    "name": "FacebookUpdated",
    "description": "",
    "employees": 11,
    "registered": true,
    "type": "Corporation"
}

DELETE - localhost:8080/companies/c03887be-e2ee-11ed-b5ea-0242ac120002

