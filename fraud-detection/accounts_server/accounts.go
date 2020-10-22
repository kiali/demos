package main

import (
	pb "github.com/kiali/demos/fraud-detection/accounts_api"
	"google.golang.org/grpc"
	"net"
	"github.com/golang/glog"
	"os"
	"flag"
	"time"
	"fmt"
	"strconv"
)

var (
	currentService = "accounts"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50051"
	wait = 10
)

type server struct {
	pb.UnimplementedAccountServiceServer
	masterAccounts 	map[string][]string
	masterMovements map[string]map[string][]*pb.Movement
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
	accounts := make(map[string][]string)
	movements := make(map[string]map[string][]*pb.Movement)

	// Init Accounts data
	accounts["user1"] = append(accounts["user1"], "ac10001")
	accounts["user1"] = append(accounts["user1"], "ac10002")
	accounts["user1"] = append(accounts["user1"], "ac10003")
	accounts["user2"] = append(accounts["user2"], "ac20001")
	accounts["user2"] = append(accounts["user2"], "ac20002")

	// Init Movements data
	for u, acs := range accounts {
		movements[u] = make(map[string][]*pb.Movement)
		for _, ac := range acs {
			for i := 0; i < 3; i++ {
				m := pb.Movement{
					OwnerId:     u,
					AccountId:   ac,
					Timestamp:   time.Now().Unix(),
					Type:        fmt.Sprintf("Operation %d", i),
					Description: fmt.Sprintf("Payment to Company %d", i),
					Quantity:    1000.10 * float64(i),
				}
				movements[u][ac] = append(movements[u][ac], &m)
			}
		}
	}
	return &server{
		masterAccounts: accounts,
		masterMovements: movements,
	}
}

func (s *server) ListAccounts(acountOwner *pb.AccountOwner, stream pb.AccountService_ListAccountsServer) error {
	if acountOwner != nil {
		if accounts, exist := s.masterAccounts[acountOwner.OwnerId]; exist {
			glog.Infof("[%s] ListAccounts for %s", instance, acountOwner.OwnerId)
			for _, ac := range accounts {
				account := pb.Account{
					OwnerId:   acountOwner.OwnerId,
					AccountId: ac,
				}
				if err := stream.Send(&account); err != nil {
					glog.Errorf("[%s] %s", instance, err)
					return err
				}
				sleep()
			}
		} else {
			glog.Warningf("[%s] %s account owner doesn't exist", instance, acountOwner.OwnerId)
			return fmt.Errorf("%s account owner doesn't exist", acountOwner.OwnerId)
		}
	} else {
		glog.Warningf("[%s] AccountOwner must be not nil", instance)
		return fmt.Errorf("AccountOwner must be not nil")
	}
	return nil
}

func (s *server) ListMovements(account *pb.Account, stream pb.AccountService_ListMovementsServer) error {
	if account != nil {
		if movs, exist := s.masterMovements[account.OwnerId][account.AccountId]; exist {
			glog.Infof("[%s] ListMovements for %s, %s", instance, account.OwnerId, account.AccountId)
			for _, m := range movs {
				mov := pb.Movement{
					OwnerId:     m.OwnerId,
					AccountId:   m.AccountId,
					Timestamp:   m.Timestamp,
					Type:        m.Type,
					Description: m.Description,
					Quantity:    m.Quantity,
				}
				if err := stream.Send(&mov); err != nil {
					glog.Errorf("[%s] %s", instance, err)
					return err
				}
				sleep()
			}
		} else {
			glog.Warningf("[%s] %s, %s account doesn't exist", instance, account.OwnerId, account.AccountId)
			return fmt.Errorf("%s, %s account doesn't exist", account.OwnerId, account.AccountId)
		}
	} else {
		glog.Warningf("[%s] Account must be not nil", instance)
		return fmt.Errorf("Account must be not nil")
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
	pb.RegisterAccountServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}

