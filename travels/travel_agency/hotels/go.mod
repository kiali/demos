module github.com/kiali/demos/travels/travel_agency/hotels

go 1.22

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/mux v1.7.4
	go.opentelemetry.io/contrib/propagators/b3 v1.31.0
	go.opentelemetry.io/otel v1.31.0
)

require go.opentelemetry.io/otel/trace v1.31.0 // indirect
