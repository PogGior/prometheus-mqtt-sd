ifndef DOCKER_IMAGE_NAME
DOCKER_IMAGE_NAME=prometheus-mqtt-sd
endif

ifndef DOCKER_IMAGE_TAG
DOCKER_IMAGE_TAG=0.0.1
endif

ifndef
MOSQUITTO_HOST=0.0.0.0
endif

ifndef
MOSQUITTO_PORT=31883
endif



docker-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

test:
	@echo "Generate certs for testing..."
	./fixtures/generate-certs.sh
	@echo "Running tests..."
	go test -v ./...

example-message:
	@echo "Publishing example message..."
	mosquitto_pub -h ${MOSQUITTO_HOST} -p ${MOSQUITTO_PORT} -t topic.test -f ./fixtures/example-message.json