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

kind_backend_loadimage:
	kind load docker-image jerry0420/queue-system-backend:v$(ver)

kind_frontend_loadimage:
	kind load docker-image jerry0420/queue-system-frontend:v$(ver)

# ==========================================================

docker_build_backend:
	mkdir backend_temp
	cp -r ./backend ./backend_temp/backend
	cp ./main.go ./backend_temp/main.go
	cp ./go.mod ./backend_temp/go.mod
	cp ./go.sum ./backend_temp/go.sum
	docker build -f Dockerfile.backend -t jerry0420/queue-system-backend:v$(ver) --no-cache ./backend_temp
	rm -r backend_temp

docker_build_frontend:
	mkdir frontend_temp frontend_temp/scripts frontend_temp/scripts/nginx
	cp -r ./src ./frontend_temp/src
	cp -r ./public ./frontend_temp/public
	cp ./package.json ./frontend_temp/package.json
	cp ./package-lock.json ./frontend_temp/package-lock.json
	cp ./scripts/nginx/nginx.conf ./frontend_temp/scripts/nginx/nginx.conf
	docker build -f Dockerfile.frontend -t jerry0420/queue-system-frontend:v$(ver) --no-cache ./frontend_temp
	rm -r frontend_temp