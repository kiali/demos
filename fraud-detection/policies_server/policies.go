package main

import (
	pb "github.com/kiali/demos/fraud-detection/policies_api"
	"github.com/golang/glog"
	"strconv"
	"os"
	"flag"
	"google.golang.org/grpc"
	"net"
	"fmt"
	"time"
)

var (
	currentService = "policies"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50062"
	wait = 10
)

type server struct {
	pb.UnimplementedPoliciesServiceServer
	masterPolicies map[string][]*pb.Policy
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
	policies := make(map[string][]*pb.Policy)

	policies["user1"] = append(policies["user1"], &pb.Policy{
		OwnerId: "user1",
		PolicyId: "policy1",
	})
	policies["user2"] = append(policies["user2"], &pb.Policy{
		OwnerId: "user2",
		PolicyId: "policy2",
	})
	policies["user2"] = append(policies["user2"], &pb.Policy{
		OwnerId: "user2",
		PolicyId: "policy3",
	})

	return &server{
		masterPolicies: policies,
	}
}

func (s *server) ListPolicies(policyOwner *pb.PolicyOwner, stream pb.PoliciesService_ListPoliciesServer) error {
	if policyOwner != nil {
		if policies, exist := s.masterPolicies[policyOwner.OwnerId]; exist {
			glog.Infof("[%s] ListPolicies for %s", instance, policyOwner.OwnerId)
			for _, policy := range policies {
				p := pb.Policy{
					OwnerId:  policy.OwnerId,
					PolicyId: policy.PolicyId,
				}
				if err := stream.Send(&p); err != nil {
					glog.Errorf("[%s] %s", instance, err)
					return err
				}
				sleep()
			}
		} else {
			glog.Warningf("[%s] %s  policy doesn't have claims", instance, policyOwner.OwnerId)
			return nil
		}
	} else {
		glog.Warningf("[%s] PolicyOwner must be not nil", instance)
		return fmt.Errorf("PolicyOwner must be not nil")
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
	pb.RegisterPoliciesServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}
