package main

import (
	"os"
	"github.com/golang/glog"
	"flag"
	"strconv"
	"google.golang.org/grpc"
	pbBank "github.com/kiali/demos/fraud-detection/bank_api"
	pbInsurance "github.com/kiali/demos/fraud-detection/insurance_api"
	"time"
	"context"
)

var (
	currentService = "fraud"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	// Invoked services
	bankAddress = "localhost:50053"
	insuranceAddress = "localhost:50063"
	wait = 500
)

func setup() {
	flag.Set("logtostderr", "true")
	flag.Parse()
	ss := os.Getenv("CURRENT_SERVICE")
	if ss != "" {
		currentService = ss
	}
	sv := os.Getenv("CURRENT_VERSION")
	if sv != "" {
		currentVersion = sv
	}
	instance = currentService + "/" + currentVersion

	ba := os.Getenv("BANK_ADDRESS")
	if ba != "" {
		bankAddress = ba
		glog.Infof("BANK_ADDRESS=%s", bankAddress)
	} else {
		glog.Warningf("BANK_ADDRESS variable empty. Using default [%s]", bankAddress)
	}

	ia := os.Getenv("INSURANCE_ADDRESS")
	if ia != "" {
		insuranceAddress = ia
		glog.Infof("INSURANCE_ADDRESS=%s", insuranceAddress)
	} else {
		glog.Warningf("INSURANCE_ADDRESS variable empty. Using default [%s]", insuranceAddress)
	}

	w := os.Getenv("WAIT")
	if w != "" {
		if waitTime, err := strconv.Atoi(w); err == nil {
			wait = waitTime
			glog.Infof("WAIT=%s", w)
		}
	}
}

func bankReport(user string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(bankAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
	}
	defer conn.Close()

	bankClient := pbBank.NewBankServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	glog.Infof("[%s] Bank FraudReport Request for %s", instance, user)
	c := pbBank.Customer{CustomerId: user}
	report, err := bankClient.FraudReport(ctx, &c)
	if err != nil {
		glog.Errorf("[%s] Error in Bank FraudReport Request. Error: %s ", instance, err)
	} else {
		glog.Infof("[%s] Bank Fraud Report. \n" +
			"Customer: %s \n" +
			"Risk: %f \n" +
			"Description: %s \n", instance, report.CustomerId, report.Risk, report.Description)
	}
}

func insuranceReport(user string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(insuranceAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
	}
	defer conn.Close()

	insuranceClient := pbInsurance.NewInsuranceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	glog.Infof("[%s] Insurance FraudReport Request for %s", instance, user)
	c := pbInsurance.Customer{CustomerId: user}
	report, err := insuranceClient.FraudReport(ctx, &c)
	if err != nil {
		glog.Errorf("[%s] Error in Insurance FraudReport Request. Error: %s ", instance, err)
	} else {
		glog.Infof("[%s] Insurance Fraud Report. \n" +
			"Customer: %s \n" +
			"Risk: %f \n" +
			"Description: %s \n", instance, report.CustomerId, report.Risk, report.Description)
	}
}

func sleep() {
	requestSleep := time.Duration(wait) * time.Millisecond
	time.Sleep(requestSleep)
}

func main() {
	setup()
	glog.Infof("Starting %s \n", instance)
	i := 0
	for {
		user := "user1"
		if i % 2 == 1 {
			user = "user2"
		}
		i++
		go bankReport(user)
		go insuranceReport(user)
		sleep()
	}
}
