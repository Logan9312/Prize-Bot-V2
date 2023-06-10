docker build --tag prize-bot .
docker run --env-file .env -p 8080:8080 prize-bot