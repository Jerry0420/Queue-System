#!/bin/sh

cat <<EOF > nginx.conf
events {

}

http {
    upstream  upstream_backend {
        server   $BACKEND_HOST:8000;
        keepalive_timeout 75s;
    }

    upstream  upstream_frontend {
        server   $FRONTEND_HOST:3000;
        keepalive_timeout 75s;
    }

    server {
        listen 80 default_server;
        listen [::]:80 default_server;
        server_name _;

        proxy_connect_timeout 75s;
        proxy_read_timeout 86400s;
        proxy_send_timeout 75s;

        location /api {
                proxy_pass http://upstream_backend;
                add_header Access-Control-Allow-Origin '*';
                add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS, PUT, DELETE, PATCH';
                add_header Access-Control-Allow-Headers 'DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Connection,Authorization';
                add_header Access-Control-Allow-Credentials 'true';
        }

        location / {
                proxy_pass http://upstream_frontend;
        }
    }
}
EOF