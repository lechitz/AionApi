#makefile
APPLICATION_NAME := aion-api
PORT := 5001

docker-compose-up: docker-compose-down
	docker-compose rm -f -v postgres
	docker-compose up

docker-compose-down:
	docker-compose down -v

docker-compose-dev-up: # docker-compose-down
	docker-compose -f docker-compose-dev.yaml down -v
	docker-compose -f docker-compose-dev.yaml rm -f -v postgres
	docker-compose -f docker-compose-dev.yaml up

test-cover:
	go test ./... -coverprofile=coverage_tmp.out
	cat coverage_tmp.out | grep -v "Mock" > coverage.out
	rm -f coverage_tmp.out
	go tool cover -html=coverage.out

docker-build:
	docker build -t $(APPLICATION_NAME):latest .

docker-build-run:	docker-build docker-run

docker-run:
	docker run -e REDIS_ADDR='redis:6379' -p $(PORT):$(PORT) -t $(APPLICATION_NAME):latest

docker-clean-all:
	#To clear containers:
	docker rm -f $(docker ps -a -q) \
	#To clear images:
	docker rmi -f $(docker images -a -q) \
	#To clear volumes:
	docker volume rm $(docker volume ls -q) \
	#To clear networks:
	docker network rm $(docker network ls | tail -n+2 | awk '{if($2 !~ /bridge|none|host/){ print $1 }}')