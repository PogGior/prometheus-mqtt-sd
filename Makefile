ifndef DOCKER_IMAGE_NAME
DOCKER_IMAGE_NAME=prometheus-mqtt-sd
endif

ifndef DOCKER_IMAGE_TAG
DOCKER_IMAGE_TAG=0.0.1
endif

docker-image:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .