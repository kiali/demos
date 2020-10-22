package main

import (
	pb "github.com/kiali/demos/fraud-detection/claims_api"
	"google.golang.org/grpc"
	"github.com/golang/glog"
	"net"
	"flag"
	"os"
	"strconv"
	"time"
	"fmt"
)

var (
	currentService = "claims"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50061"
	wait = 10
)

type server struct {
	pb.UnimplementedClaimsServiceServer
	masterClaims map[string][]*pb.Claim
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

	w := os.Getenv("WAIT")
	if w != "" {
		if waitTime, err := strconv.Atoi(w); err == nil {
			wait = waitTime
			glog.Infof("WAIT=%s", w)
		}
	}
}

func newServer() *server {
	claims := make(map[string][]*pb.Claim)

	// Init policies
	policies := []string{"policy1", "policy2", "policy3"}

	// Init claims
	for _, p := range policies {
		claims[p] = []*pb.Claim{}
		for i := 0; i < 3; i++ {
			c := pb.Claim{
				PolicyId:    p,
				Timestamp:   time.Now().Unix(),
				Type:        fmt.Sprintf("Claim %d", i),
				Description: fmt.Sprintf("Claim for topic %d", i),
				Value:       (float64)(1000 + i*200),
			}
			claims[p] = append(claims[p], &c)
		}
	}

	return &server{
		masterClaims: claims,
	}
}

func (s *server) ListClaims(policy *pb.Policy, stream pb.ClaimsService_ListClaimsServer) error {
	if policy != nil {
		if claims, exist := s.masterClaims[policy.PolicyId]; exist {
			glog.Infof("[%s] ListClaims for %s", instance, policy.PolicyId)
			for _, c := range claims {
				c := pb.Claim{
					PolicyId:    c.PolicyId,
					Timestamp:   c.Timestamp,
					Type:        c.Type,
					Description: c.Description,
					Value:       c.Value,
				}
				if err := stream.Send(&c); err != nil {
					glog.Errorf("[%s] %s", instance, err)
					return err
				}
				sleep()
			}
		} else {
			glog.Warningf("[%s] %s  policy doesn't have claims", instance, policy.PolicyId)
			return nil
		}
	} else {
		glog.Warningf("[%s] Policy must be not nil", instance)
		return fmt.Errorf("Policy must be not nil")
	}
	return nil
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
	pb.RegisterClaimsServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}
