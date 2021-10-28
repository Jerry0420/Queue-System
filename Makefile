up_db:
	docker-compose -f ./docker-compose/docker-compose.db.yml up -d --build

exec_db:
	docker exec -it migration_tools sh

down_db:
	docker-compose -f ./docker-compose/docker-compose.db.yml down

dowm_migration_tools:
	docker stop migration_tools
	docker rm migration_tools

# ========================================================

up_build_dev:
	docker-compose -f ./docker-compose/docker-compose.dev.yml up -d --build

down_dev:
	docker-compose -f ./docker-compose/docker-compose.dev.yml down

exec_backend:
	docker exec -it backend sh

exec_frontend:
	docker exec -it frontend sh

logs_backend:
	docker-compose logs -f backend

# =========================================================

up_build:
	docker-compose -f ./docker-compose/docker-compose.yml up -d --build

down:
	docker-compose -f ./docker-compose/docker-compose.yml down

# ==========================================================

kind_create:
	kind create cluster --config ./scripts/kind.yaml

kind_delete:
	kind delete cluster

kind_loadimage:
	kind load docker-image jerry0420/queue-system:v$(ver)

# ==========================================================

docker_build_backend:
	docker build -f Dockerfile.backend -t jerry0420/queue-system-backend:v$(ver) .

docker_build_frontend:
	docker build -f Dockerfile.frontend -t jerry0420/queue-system-frontend:v$(ver) .