up_db:
	docker-compose -f ./docker-compose/docker-compose.db.yml up -d --build

exec_db:
	docker exec -it migration_tools sh

down_db:
	docker-compose -f ./docker-compose/docker-compose.db.yml down

# ========================================================

up_build_dev:
	docker-compose -f ./docker-compose/docker-compose.dev.yml up -d --build

down_dev:
	docker-compose -f ./docker-compose/docker-compose.dev.yml down

exec_backend:
	docker exec -it backend sh

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

kind_loadimage_backend:
	kind load docker-image jerry0420/queue-system-backend:v$(ver)

kind_loadimage_frontend:
	kind load docker-image jerry0420/queue-system-frontend:v$(ver)

kind_loadimage_grpc:
	kind load docker-image jerry0420/queue-system-grpc:v$(ver)

# ==========================================================

docker_build_backend:
	mkdir backend_temp
	cp -r ./backend ./backend_temp/backend
	docker build -f Dockerfile.backend -t jerry0420/queue-system-backend:v$(ver) --no-cache ./backend_temp
	rm -r backend_temp

docker_build_frontend:
	mkdir frontend_temp frontend_temp/scripts frontend_temp/scripts/nginx
	cp -r ./frontend ./frontend_temp/frontend
	cp ./scripts/nginx/nginx.conf ./frontend_temp/scripts/nginx/nginx.conf
	docker build -f Dockerfile.frontend -t jerry0420/queue-system-frontend:v$(ver) --no-cache ./frontend_temp
	rm -r frontend_temp

docker_build_grpc:
	mkdir grpc_temp
	cp -r ./grpc ./grpc_temp/grpc
	docker build -f Dockerfile.grpc -t jerry0420/queue-system-grpc:v$(ver) --no-cache ./grpc_temp
	rm -r grpc_temp