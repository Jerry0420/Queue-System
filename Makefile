up_build_dev:
	docker-compose -f ./docker-compose/docker-compose.yml -f ./docker-compose/docker-compose.dev.yml up -d --build

down_dev:
	docker-compose -f ./docker-compose/docker-compose.yml -f ./docker-compose/docker-compose.dev.yml down
	docker volume rm docker-compose_db_data

exec_backend_dev:
	docker exec -it backend sh

exec_frontend_dev:
	docker exec -it frontend sh

up_build:
	docker-compose -f ./docker-compose/docker-compose.yml up -d --build

down:
	docker-compose -f ./docker-compose/docker-compose.yml down

exec_backend:
	docker exec -it backend sh

logs_backend:
	docker-compose logs -f backend

kind_create:
	kind create cluster --config ./scripts/kind.yaml

kind_delete:
	kind delete cluster

kind_loadimage:
	kind load docker-image jerry0420/queue-system:v$(ver)

docker_build:
	docker build -f Dockerfile -t jerry0420/queue-system:v$(ver) .