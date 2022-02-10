# variables in Makefile
version := v0.0.1
image := quay.io/leandroberetta/topogen
dockerfile := build/Dockerfile

# build generator binary
build-binary: clean build-ui
	cp -r ui/build pkg/controller/
	go build app/topogen.go
	rm -rf pkg/controller/build

# build generator binary in image
build-binary-image: 
	cp -r ui/build pkg/controller/
	go build app/topogen.go
	rm -rf pkg/controller/build

# build ui in production
build-ui:
	cd ui && yarn && yarn build

# build docker image
build-image: clean build-ui
	docker build -t $(image):$(version) -f $(dockerfile) .

# push docker image
push-image:
	docker push $(image):$(version)

# clean built ui files
clean-ui:
	rm -rf ui/build
	rm -rf pkg/controller/build

# clean generated files
clean-gen:
	rm -rf topogen deploy.json

# clean all 
clean: clean-ui clean-gen
