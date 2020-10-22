module github.com/kiali/demos/fraud-detection/bank_server

go 1.14

require (
	github.com/kiali/demos/fraud-detection/accounts_api v0.0.0
	github.com/kiali/demos/fraud-detection/cards_api v0.0.0
	github.com/kiali/demos/fraud-detection/bank_api v0.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	google.golang.org/grpc v1.33.0
)

replace github.com/kiali/demos/fraud-detection/accounts_api => ../accounts_api
replace github.com/kiali/demos/fraud-detection/cards_api => ../cards_api
replace github.com/kiali/demos/fraud-detection/bank_api => ../bank_api