# variables in Makefile
version := v0.0.1
image := quay.io/leandroberetta/topogen
dockerfile := build/Dockerfile

# build multiarch topogen binary
build-multiarch-binary: build-binary-linux build-binary-darwin

# build topogen binary
build-binary: clean build-ui
	cp -r ui/build pkg/controller/
	go build app/topogen.go -o topogen
	rm -rf pkg/controller/build

# build linux topogen binary
build-binary-linux: clean build-ui
	cp -r ui/build pkg/controller/
	CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -o topogen.linux app/topogen.go 
	rm -rf pkg/controller/build

# build darwin topogen binary
build-binary-darwin: clean build-ui
	cp -r ui/build pkg/controller/
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o topogen.darwin app/topogen.go 
	rm -rf pkg/controller/build

# build topogen docker image
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
	rm -rf topogen deploy.json topogen.linux topogen.darwin

# clean all 
clean: clean-ui clean-gen
