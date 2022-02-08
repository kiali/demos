# variables in Makefile
version := v0.0.1
image_tag := quay.io/leandroberetta/topogen
dockerfile := build/Dockerfile

# build generator binary
build-binary: clean build-ui
	cp -r ui/build pkg/controller/
	go build app/topogen.go
	rm -rf pkg/controller/build

# build ui in production
build-ui:
	cd ui && yarn && yarn build

# build docker image
build-image: clean build-ui
	docker build -t $(image_tag):$(version) -f $(dockerfile) .

# clean built ui files
clean-ui:
	rm -rf ui/build ui/node_modules
	rm -rf pkg/controller/build

# clean generated files
clean-gen:
	rm -rf topogen deploy.json

# clean all 
clean: clean-ui clean-gen
