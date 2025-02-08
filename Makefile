UID ?= $(shell bash -c 'id -u')
GID ?= $(shell bash -c 'id -g')
CWD ?= $(shell pwd)

DOCKER_COMPOSE_DEV = docker compose -f docker-compose.dev.yml

all: npm-install webpack tailwind compress-static-files

npm-install:
	@echo "Installing npm packages..."
	@cd frontend && npm install
	@echo "npm packages installed!"

webpack:
	@echo "Running Webpack..."
	@cd frontend && npm run prod:webpack
	@echo "Webpack finished!"
	
tailwind:
	@echo "Compiling Tailwind CSS..."
	@cd frontend && npx tailwindcss build -i src/css/style.css -o src/css/tailwind.css
	@echo "Tailwind CSS compiled!"

compress-static-files:
	@echo "Compressing static files..."
	@for file in frontend/src/js/build/*.js; do \
		gzip -c $$file > frontend/public/$$(basename $$file).gz; \
	done
	@gzip -k -f frontend/src/css/tailwind.css -c > frontend/public/tailwind.css.gz
	@gzip -k -f frontend/src/js/htmx/htmx.min.js -c > frontend/public/htmx.min.js.gz
	@find ./frontend/public -type f ! -path "./frontend/public/fonts/*" ! -name "*.gz" -exec sh -c 'gzip -c "$1" > "$1.gz"' _ {} \;
	@echo "Static files compressed!"

.PHONY:connect-db-dev
connect-db:
	$(DOCKER_COMPOSE_DEV) exec -it postgres psql -U admin -d backend

dev:
	@echo "Running Tailwind CSS..."
	@cd frontend && npx tailwindcss build -i src/css/style.css -o src/css/tailwind.css --watch
	@echo "Tailwind CSS finished"