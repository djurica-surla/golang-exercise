# Runs the app locally if you change config to connect to your local instances.
run:
	go run cmd/server/main.go

# Should fetch the golangci-lint binary and place in in the bin of your gopath.
get-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.50.1

# Runs the linter to check for potential errors.
lint:
	golangci-lint run ./...

# Attempts to automatically fix linting issues.
lint-fix:
	golangci-lint run --fix

# Builds the necessary environment and runs integration tests for our app.
# Postgres volume will be fresh for every run.
test-integration:
	docker compose -f test/docker-compose.integration.yaml up --build --abort-on-container-exit --force-recreate
	docker-compose -f test/docker-compose.integration.yaml down -v

# Runs the docker compose to create the app and postgres instance. 
# Config is setup to connect to this docker-compose. App will be opened on port 8080 locally.
# For WSL you might need to run sudo apt install docker-compose and than run the command as docker-compose -f ...
start:
	docker compose -f docker-compose.production.yaml up --build	

# Commands to delete all containers and volumes for the fresh start if necessary.
# docker rm -f $(docker ps -a -q) 
# docker volume rm $(docker volume ls -q)
