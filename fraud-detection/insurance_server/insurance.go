package main

import (
	"github.com/golang/glog"
	"os"
	"flag"
	"strconv"
	pb "github.com/kiali/demos/fraud-detection/insurance_api"
	pbClaims "github.com/kiali/demos/fraud-detection/claims_api"
	pbPolicies "github.com/kiali/demos/fraud-detection/policies_api"
	"google.golang.org/grpc"
	"net"
	"context"
	"fmt"
	"sync"
	"strings"
	"time"
	"io"
)

var (
	currentService = "insurance"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50063"
	// Invoked services
	claimsAddress = "localhost:50061"
	policiesAddress = "localhost:50062"
	wait = 10
)

type server struct {
	pb.UnimplementedInsuranceServiceServer
}

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

	la := os.Getenv("LISTEN_ADDRESS")
	if la != "" {
		listenAddress = la
		glog.Infof("LISTEN_ADDRESS=%s", listenAddress)
	} else {
		glog.Warningf("LISTEN_ADDRESS variable empty. Using default [%s]", listenAddress)
	}

	aa := os.Getenv("CLAIMS_ADDRESS")
	if aa != "" {
		claimsAddress = aa
		glog.Infof("CLAIMS_ADDRESS=%s", claimsAddress)
	} else {
		glog.Warningf("CLAIMS_ADDRESS variable empty. Using default [%s]", claimsAddress)
	}
	ca := os.Getenv("POLICIES_ADDRESS")
	if ca != "" {
		policiesAddress = ca
		glog.Infof("POLICIES_ADDRESS=%s", policiesAddress)
	} else {
		glog.Warningf("POLICIES_ADDRESS variable empty. Using default [%s]", policiesAddress)
	}

	w := os.Getenv("WAIT")
	if w != "" {
		if waitTime, err := strconv.Atoi(w); err == nil {
			wait = waitTime
			glog.Infof("WAIT=%s", w)
		}
	}
}

func newServer() *server {
	return &server{}
}

func (s *server) FraudReport(ctx context.Context, customer *pb.Customer) (*pb.InsuranceReport, error) {
	if customer != nil {
		glog.Infof("[%s] FraudReport for %s", instance, customer.CustomerId)
		policies, err := getPolicies(customer.CustomerId)
		if err != nil {
			return nil, err
		}

		wg := sync.WaitGroup{}
		wg.Add(len(policies))

		claimsChan := make(chan []*pbClaims.Claim, len(policies))
		errChan := make(chan error, len(policies))

		for _, p := range policies {
			go func(policyId string, claimsChan chan []*pbClaims.Claim, errChan chan error) {
				defer wg.Done()
				claims, err := getClaims(policyId)
				if err != nil {
					errChan <- err
				} else {
					claimsChan <- claims
				}
			}(p, claimsChan, errChan)
		}
		wg.Wait()
		close(errChan)
		close(claimsChan)

		if len(errChan) > 0 {
			policiesError := []string{}
			for e := range errChan {
				policiesError = append(policiesError, fmt.Sprintf("%s ", e.Error()))
			}
			return nil, fmt.Errorf("%s", strings.Join(policiesError, ":"))
		}
		
		var risk float64
		risk = 1000.10

		if len(claimsChan) > 0 {
			for claims := range claimsChan {
				for _, c := range claims {
					risk = risk + c.Value * 10.10
				}
			}
		}
		ir := pb.InsuranceReport{
			CustomerId:  customer.CustomerId,
			Risk:        risk,
			Description: fmt.Sprintf("InsuranceReport generated on %s", time.Now().String()),
		}
		sleep()
		return &ir, nil
	}
	return nil, fmt.Errorf("Customer must be not nil")
}

func getPolicies(customerId string) ([]string, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(policiesAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
		return nil, err
	}
	defer conn.Close()

	policiesClient := pbPolicies.NewPoliciesServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	po := pbPolicies.PolicyOwner{OwnerId:  customerId}
	stream, err := policiesClient.ListPolicies(ctx, &po)
	if err != nil {
		glog.Errorf("[%s] ListPolicies error: %s", instance, err)
		return nil, err
	}
	policies := []string{}
	for {
		policy, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Errorf("[%s] ListPolicies. Error: %s", instance, err)
			return nil, err
		}
		glog.Infof("[%s] Policy: %s, %s", instance, policy.OwnerId, policy.PolicyId)
		policies = append(policies, policy.PolicyId)
	}
	return policies, nil
}

func getClaims(policyId string) ([]*pbClaims.Claim, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(claimsAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
		return nil, err
	}
	defer conn.Close()

	claimsClient := pbClaims.NewClaimsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p := pbClaims.Policy{PolicyId: policyId}
	stream, err := claimsClient.ListClaims(ctx, &p)
	if err != nil {
		return nil, err
	}
	claims := []*pbClaims.Claim{}
	for {
		claim, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Errorf("[%s] ListClaims. Error: %s", instance, err)
			return nil, err
		}
		c := pbClaims.Claim{
			PolicyId:   claim.PolicyId,
			Timestamp:   claim.Timestamp,
			Type:        claim.Type,
			Description: claim.Description,
			Value:       claim.Value,
		}
		glog.Infof("[%s] Claim: %s - %d - %s - %s - %f", instance, claim.PolicyId, claim.Timestamp, claim.Type, claim.Description, claim.Value)
		claims = append(claims, &c)
	}
	return claims, nil
}

func sleep() {
	requestSleep := time.Duration(wait) * time.Millisecond
	time.Sleep(requestSleep)
}

func main() {
	setup()
	glog.Infof("Starting %s \n", instance)
	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		glog.Fatalf("[%s] failed to listen: %v", instance, err)
	}
	s := grpc.NewServer()
	pb.RegisterInsuranceServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}
