package main

import (
	"context"
	"github.com/golang/glog"
	"os"
	"flag"
	pb "github.com/kiali/demos/fraud-detection/bank_api"
	pbAccounts "github.com/kiali/demos/fraud-detection/accounts_api"
	pbCards "github.com/kiali/demos/fraud-detection/cards_api"
	"google.golang.org/grpc"
	"time"
	"io"
	"fmt"
	"sync"
	"net"
	"strconv"
)

var (
	currentService = "bank"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50053"
	// Invoked services
	accountsAddress = "localhost:50051"
	cardsAddress = "localhost:50052"
	wait = 10
)

type server struct {
	pb.UnimplementedBankServiceServer
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

	aa := os.Getenv("ACCOUNTS_ADDRESS")
	if aa != "" {
		accountsAddress = aa
		glog.Infof("ACCOUNTS_ADDRESS=%s", accountsAddress)
	} else {
		glog.Warningf("ACCOUNTS_ADDRESS variable empty. Using default [%s]", accountsAddress)
	}
	ca := os.Getenv("CARDS_ADDRESS")
	if ca != "" {
		cardsAddress = ca
		glog.Infof("CARDS_ADDRESS=%s", cardsAddress)
	} else {
		glog.Warningf("CARDS_ADDRESS variable empty. Using default [%s]", cardsAddress)
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

func (s *server) FraudReport(ctx context.Context, customer *pb.Customer) (*pb.BankReport, error) {
	if customer != nil {
		var riskA float64
		var errA error
		var riskC float64
		var errC error
		wg := sync.WaitGroup{}
		wg.Add(2)

		glog.Infof("[%s] FraudReport for %s", instance, customer.CustomerId)
		
		go func() {
			defer wg.Done()
			riskA, errA = riskAccounts(customer.CustomerId)	
		}()

		go func() {
			defer wg.Done()
			riskC, errC = riskCards(customer.CustomerId)
		}()
		
		wg.Wait()
		
		if errA != nil {
			return nil, errA
		}
		if errC != nil {
			return nil, errC
		}
		
		bankReport := pb.BankReport{
			CustomerId:  customer.CustomerId,
			Risk:        riskA + riskC,
			Description: fmt.Sprintf("BankReport generated on %s", time.Now().String()),
		}
		sleep()
		return &bankReport, nil
	}
	return nil, fmt.Errorf("Customer must be not nil")
}

func riskAccounts(ownerId string) (float64, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(accountsAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
		return 0, err
	}
	defer conn.Close()

	accountsClient := pbAccounts.NewAccountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	risk := (float64)(100.10)

	ao := pbAccounts.AccountOwner{OwnerId: ownerId}

	stream, err := accountsClient.ListAccounts(ctx, &ao)
	if err != nil {
		glog.Errorf("[%s] ListAccounts error: %s", instance, err)
		return risk, err
	}
	for {
		account, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Errorf("[%s] ListAccounts. Error: %s", instance, err)
			return risk, err
		}
		glog.Infof("[%s] Account: %s, %s", instance, account.OwnerId, account.AccountId)
		ac := pbAccounts.Account{OwnerId: account.OwnerId, AccountId: account.AccountId}
		movsStream, err := accountsClient.ListMovements(ctx, &ac)
		if err != nil {
			glog.Errorf("[%s] ListMovements error: %s", instance, err)
			return risk, err
		}
		for {
			mov, err := movsStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				glog.Errorf("[%s] ListMovements. Error: %s", instance, err)
				return risk, err
			}
			glog.Infof("[%s] Account Movement (Timestamp %d) %s - %f - %s ", instance, mov.Timestamp, mov.Type, mov.Quantity, mov.Description)
			risk = risk + (mov.Quantity * 0.20)
		}
	}
	glog.Infof("[%s] Risk Account for %s: %f", instance, ownerId, risk)
	return risk, nil
}

func riskCards(ownerId string) (float64, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(cardsAddress, opts...)
	if err != nil {
		glog.Fatalf("[%s] failed to open: %v", instance, err)
		return 0, err
	}
	defer conn.Close()

	cardsClient := pbCards.NewCardsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	risk := (float64)(100.10)

	co := pbCards.CardOwner{OwnerId: ownerId}

	stream, err := cardsClient.ListCards(ctx, &co)
	if err != nil {
		glog.Errorf("[%s] ListCards error: %s", instance, err)
		return risk, err
	}
	for {
		card, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Errorf("[%s] ListCards. Error: %s", instance, err)
			return risk, err
		}
		glog.Infof("[%s] Card: %s, %s", instance, card.OwnerId, card.CardId)
		cc := pbCards.Card{OwnerId: card.OwnerId, CardId: card.CardId}
		movsStream, err := cardsClient.ListMovements(ctx, &cc)
		if err != nil {
			glog.Errorf("[%s] ListMovements error: %s", instance, err)
			return risk, err
		}
		for {
			mov, err := movsStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				glog.Errorf("[%s] ListMovements. Error: %s", instance, err)
				return risk, err
			}
			glog.Infof("[%s] Card Movement (Timestamp %d) %s - %f - %s ", instance, mov.Timestamp, mov.Type, mov.Quantity, mov.Description)
			risk = risk + (mov.Quantity * 0.10)
		}
	}
	glog.Infof("[%s] Risk Card for %s: %f", instance, ownerId, risk)
	return risk, nil
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
	pb.RegisterBankServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}
