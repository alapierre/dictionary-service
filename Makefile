CONNECTION_STRING=postgres://app:qwedsazxc@localhost:5432/app?sslmode=disable
IMAGE_VERSION=0.0.1

modelgen:
	genna model -c $(CONNECTION_STRING) -o model/model.go -k -g 9

build:
	cd cmd/dictionary-service && CGO_ENABLED=0 go build -a -installsuffix cgo -o dictionary-service .

docker: build
	cd cmd/dictionary-service && docker build -t lapierre/dictionary-service:$(IMAGE_VERSION) .

push:
	docker push lapierre/dictionary-service:$(IMAGE_VERSION)
