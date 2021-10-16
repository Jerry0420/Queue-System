up_build_dev:
	docker-compose -f ./docker-compose/docker-compose.yml -f ./docker-compose/docker-compose.dev.yml up -d --build

down_dev:
	docker-compose -f ./docker-compose/docker-compose.yml -f ./docker-compose/docker-compose.dev.yml down

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

up_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

up:
	docker-compose -f ./docker-compose/docker-compose.yml up -d

logs_backend:
	docker-compose logs -f backend