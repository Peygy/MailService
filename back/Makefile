BIN = backend

.PHONY: build
build:
	$(MAKE) docker pass=$(pass)
	$(MAKE) clean
	go build -o $(BIN) ./main.go 
	go run $(BIN) -dev

.PHONY: docker
docker:
	echo $(pass) | sudo -S docker-compose up -d

clean:
	rm -f $(BIN)

#mkwindd
.PHONY: docker_clean
docker_clean:
	echo $(pass) | sudo -S docker-compose down -v