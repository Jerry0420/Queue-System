up_db:
	docker-compose -f ./compose/docker-compose.db.yml up -d --build

exec_db:
	docker exec -it migration_tools sh

down_db:
	docker-compose -f ./compose/docker-compose.db.yml down

# ========================================================

up_test:
	docker-compose -f ./compose/docker-compose.test.yml up -d --build

exec_test:
	docker exec -it backend_test sh

down_test:
	docker-compose -f ./compose/docker-compose.test.yml down

# ========================================================

up_build_dev:
	docker-compose -f ./compose/docker-compose.dev.yml up -d --build

re_build_nginx:
	docker-compose -f ./compose/docker-compose.dev.yml up -d --force-recreate --build nginx

down_dev:
	docker-compose -f ./compose/docker-compose.dev.yml down

exec_backend:
	docker exec -it backend sh

# =========================================================

up_build:
	docker-compose -f ./compose/docker-compose.yml up -d --build

down:
	docker-compose -f ./compose/docker-compose.yml down

# ==========================================================

add_crontab:
	(crontab -l 2>/dev/null; echo "* * * * * curl --max-time 30 --connect-timeout 5 -X DELETE --url $(server)/api/v1/routine/stores") | crontab -

# ==========================================================

kind_create:
	kind create cluster --config ./k8s/kind.yaml

kind_delete:
	kind delete cluster

kind_loadimage_backend:
	kind load docker-image queue-system-backend:v$(ver)

kind_loadimage_frontend:
	kind load docker-image queue-system-frontend:v$(ver)

kind_loadimage_grpc:
	kind load docker-image queue-system-grpc:v$(ver)

# ==========================================================

docker_build_backend:
	mkdir backend_temp
	cp -r ./backend ./backend_temp/backend
	docker build -f Dockerfile.backend -t queue-system-backend:v$(ver) --no-cache ./backend_temp
	rm -r backend_temp
	docker save queue-system-backend:v$(ver) > queue-system-backend.tar

docker_build_frontend:
	mkdir frontend_temp frontend_temp/scripts frontend_temp/scripts/nginx
	cp -r ./frontend ./frontend_temp/frontend
	cp ./scripts/nginx/nginx.frontend.conf ./frontend_temp/scripts/nginx/nginx.frontend.conf
	docker build -f Dockerfile.frontend -t queue-system-frontend:v$(ver) --no-cache ./frontend_temp
	rm -r frontend_temp
	docker save queue-system-frontend:v$(ver) > queue-system-frontend.tar

docker_build_grpc:
	mkdir grpc_temp
	cp -r ./grpc ./grpc_temp/grpc
	docker build -f Dockerfile.grpc -t queue-system-grpc:v$(ver) --no-cache ./grpc_temp
	rm -r grpc_temp
	docker save queue-system-grpc:v$(ver) > queue-system-grpc.tar

# ==============================================================================

vm_create:
	multipass launch --cpus 2 --disk 15G --mem 6G --name $(name) 20.04
	multipass mount ./ $(name)
	multipass exec $(name) -- sudo apt install make
	multipass shell $(name)

microk8s_install:
	sudo snap install microk8s --classic
	sudo usermod -a -G microk8s ubuntu
	sudo snap alias microk8s.kubectl kubectl
	sudo chown -f -R ubuntu ~/.kube
	newgrp microk8s

microk8s_docker_install:
	sudo apt update
	sudo apt install docker.io
	sudo gpasswd -a $(user) docker

microk8s_load_images_backend:
	microk8s ctr image import queue-system-backend.tar 

microk8s_load_images_frontend:
	microk8s ctr image import queue-system-frontend.tar 

microk8s_load_images_grpc:
	microk8s ctr image import queue-system-grpc.tar 

microk8s_enable_addons:
	microk8s enable dns ingress hostpath-storage metallb

# microk8s config
# hostname -I | awk '{print $1}'
# multipass find
# multipass info --all