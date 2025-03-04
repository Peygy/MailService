BIN = backend

.PHONY: build
build:
	$(MAKE) docker
	$(MAKE) clean
	go build -o $(BIN) ./main.go 
	go run $(BIN)

.PHONY: docker
docker:
	sudo docker-compose up -d

clean:
	rm -f $(BIN)

.PHONY: docker_clean
docker_clean:
	sudo docker-compose down -v