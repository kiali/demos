module "github.com/kiali/demos/fraud-detection/fraud"

go 1.14

require (
	github.com/kiali/demos/fraud-detection/bank_api v0.0.0
	github.com/kiali/demos/fraud-detection/insurance_api v0.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	google.golang.org/grpc v1.33.0
)

replace github.com/kiali/demos/fraud-detection/bank_api => ../bank_api
replace github.com/kiali/demos/fraud-detection/insurance_api => ../insurance_api
