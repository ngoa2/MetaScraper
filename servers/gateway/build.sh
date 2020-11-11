
GOOS=linux go build
docker build -t ngoa2/ngoa2server .
docker push ngoa2/ngoa2server

go clean