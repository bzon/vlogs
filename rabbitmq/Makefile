DOCKER_REPO=bzon/demo-rabbitmq

clean:
	rm -fr bin/**

dep:
	dep ensure

build-docker: clean dep
	GO111MODULE=off
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/publisher ./publisher/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/subscriber ./subscriber/
	docker build -t $(DOCKER_REPO)-publisher -f ./Dockerfile.publisher .
	docker build -t $(DOCKER_REPO)-subscriber -f ./Dockerfile.subscriber .

push:
	docker push $(DOCKER_REPO)-publisher:latest
	docker push $(DOCKER_REPO)-subscriber:latest
