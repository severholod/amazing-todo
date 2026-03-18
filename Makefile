include .env
export

export PROJECT_ROOT=$(CURDIR)

env-up:
	@docker compose up -d app-postgres

env-down:
	@docker compose down app-postgres

env-cleanup:
	@read -p "Are you sure? [y/n]:" ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down app-postgres && \
	  rm -rf out/pgdata && \
	  echo "Files has been deleted"; \
  	else \
  	  echo "Cancel"; \
  	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Pass the parameter 'seq'. Example: make migrate-create seq=init"; \
  		exit 1; \
	fi; \
	docker compose run --rm \
	 	app-postgres-migrate \
		create \
		-ext sql \
		-dir . \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
      		echo "Pass the parameter 'action'. Example: make migrate-action action=up"; \
      		exit 1; \
	fi; \
	docker compose run --rm \
		app-postgres-migrate \
		-path . \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@app-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
		"$(action)"

env-port-forward:
	@docker compose up -d port-forwarder
env-port-close:
	@docker compose down port-forwarder