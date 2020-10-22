module github.com/kiali/demos/fraud-detection/policies_server

go 1.14

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kiali/demos/fraud-detection/policies_api v0.0.0
	google.golang.org/grpc v1.33.0
)

replace github.com/kiali/demos/fraud-detection/policies_api => ../policies_api
