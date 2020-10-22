module github.com/kiali/demos/fraud-detection/insurance_server

go 1.14

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kiali/demos/fraud-detection/claims_api v0.0.0
	github.com/kiali/demos/fraud-detection/insurance_api v0.0.0
	github.com/kiali/demos/fraud-detection/policies_api v0.0.0
	google.golang.org/grpc v1.33.0

)

replace github.com/kiali/demos/fraud-detection/insurance_api => ../insurance_api

replace github.com/kiali/demos/fraud-detection/claims_api => ../claims_api

replace github.com/kiali/demos/fraud-detection/policies_api => ../policies_api
