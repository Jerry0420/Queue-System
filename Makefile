up_build_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

up_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

down_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

exec_backend_dev:
	docker exec -it backend sh

up_build:
	docker-compose up -d --build

up:
	docker-compose up -d

down:
	docker-compose down

exec_backend:
	docker exec -it backend sh

logs_backend:
	docker-compose logs -f backend