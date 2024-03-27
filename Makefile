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
MOSQUITTO_PORT=1883
endif



docker-image:
	@echo "Building docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

docker-save:
	@echo "Saving docker image..."
	docker save $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -o $(DOCKER_IMAGE_NAME)_$(DOCKER_IMAGE_TAG).tar

build:
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o prometheus-mqtt-sd ./cmd/prometheus-mqtt-sd/main.go

run:
	@echo "Running binary..."
	./prometheus-mqtt-sd --config.file ./fixtures/config-simple.yaml --output.file ./fixtures/output.json

test:
	@echo "Generate certs for testing..."
	./fixtures/generate-certs.sh
	@echo "Running tests..."
	go test -v ./...

example-message:
	@echo "Publishing example message..."
	mosquitto_pub -h ${MOSQUITTO_HOST} -p ${MOSQUITTO_PORT} -t topic.test -f ./fixtures/example-message.json