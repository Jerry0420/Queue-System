up_build_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

up_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

down_dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down

exec_backend_dev:
	docker exec -it backend sh

exec_frontend_dev:
	docker exec -it frontend sh

up_build:
	docker-compose build --no-cache
	docker-compose up -d

up:
	docker-compose up -d

down:
	docker-compose down

exec_backend:
	docker exec -it backend sh

logs_backend:
	docker-compose logs -f backend