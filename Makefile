.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration sqlc api proto
MIGRATE = migrate -database "$(DB_URL)" -verbose
DB_URL_user = "postgres://root:secret@localhost:5432/user?sslmode=disable"
DB_URL_hub = "postgres://root:secret@localhost:5432/hub?sslmode=disable"

MIGRATE = migrate -database
MIGRATE_CREATE = migrate create -ext sql -seq
# List of services
SERVICES := user hub

# List of PostgreSQL instances
POSTGRES_SERVICES = postgres-core
POSTGRES_IMAGE = postgres:16-alpine
POSTGRES_USER = root
POSTGRES_PASSWORD = secret
NETWORK = core-network

# Migrate Up (Run all services)
migrateup: $(patsubst %, migrateup-%, $(SERVICES))

# Migrate Down (Run all services)
migratedown: $(patsubst %, migratedown-%, $(SERVICES))

# Migration up for specific service
migrateup-%:
	$(MIGRATE) "$(DB_URL_$*)" -verbose -path ./$(*)/migrations up

# Migration down for specific service
migratedown-%:
	$(MIGRATE) "$(DB_URL_$*)" -verbose -path ./$(*)/migrations down

# Create New Migration for a specific service
new-migration-%:
	@$(if $(name), , $(error Usage: make new-migration-<service> name=<migration-name>))
	$(MIGRATE_CREATE) -dir ./$*/migrations $(name)

postgres:
	docker run --name $(POSTGRES_SERVICES) --network $(NETWORK) -p 5432:5432 \
		-e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d $(POSTGRES_IMAGE)

# Create database inside the 'postgres-user' container
createdb-%:
	docker exec -it $(POSTGRES_SERVICES) createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $*

# Drop database inside the 'postgres-user' container
dropdb-%:
	docker exec -it $(POSTGRES_SERVICES) dropdb $*

# Stop and remove containers
clean:
	docker rm -f $(POSTGRES_SERVICES)

network:
	docker network create $(NETWORK)

redis:
	docker run --name redis --network $(NETWORK) -p 6379:6379 -d redis:7.4.1

scylla-mac:
	docker run --name scylla -d scylladb/scylla -p 9042:9042 -p 9180:9180 --network $(NETWORK)

sqlc-user:
	cd user && sqlc generate

sqlc-hub:
	cd hub && sqlc generate


auth-dev:
	go run ./auth -config=configs/config.dev.yaml
user-dev:
	go run ./user -config=configs/config.dev.yaml
hub-dev:
	go run ./hub -config=configs/config.dev.yaml

proto:
	protoc proto/*/*.proto --go-grpc_out=. --go_out=.
	