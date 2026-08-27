up:
	docker compose up --build
	start http://localhost:8080

start:
	docker start postgres_db go_app

intro_bd:
	docker exec -it postgres_db psql -U postgres -d app_db

start only db:
	docker compose up -d --build db

start migration:
docker compose up -d --build db migrate