package main

import (
	pb "github.com/kiali/demos/fraud-detection/cards_api"
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
	currentService = "cards"
	currentVersion = "no-version"
	instance = currentService + "/" + currentVersion
	listenAddress = ":50052"
	wait = 10
)

type server struct {
	pb.UnimplementedCardsServiceServer
	masterCards 	map[string][]string
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
	cards := make(map[string][]string)
	movements := make(map[string]map[string][]*pb.Movement)

	// Init Accounts data
	cards["user1"] = append(cards["user1"], "cd10001")
	cards["user1"] = append(cards["user1"], "cd10002")
	cards["user1"] = append(cards["user1"], "cd10003")
	cards["user2"] = append(cards["user2"], "cd20001")
	cards["user2"] = append(cards["user2"], "cd20002")

	// Init Movements data
	for u, cds := range cards {
		movements[u] = make(map[string][]*pb.Movement)
		for _, cd := range cds {
			for i := 0; i < 3; i++ {
				m := pb.Movement{
					OwnerId:     u,
					CardId:   cd,
					Timestamp:   time.Now().Unix(),
					Type:        fmt.Sprintf("Operation %d", i),
					Description: fmt.Sprintf("Payment to Company %d", i),
					Quantity:    1000.10 * float64(i),
				}
				movements[u][cd] = append(movements[u][cd], &m)
			}
		}
	}
	return &server{
		masterCards:  cards,
		masterMovements: movements,
	}
}

func (s *server) ListCards(cardOwner *pb.CardOwner, stream pb.CardsService_ListCardsServer) error {
	if cardOwner != nil {
		if cards, exist := s.masterCards[cardOwner.OwnerId]; exist {
			glog.Infof("[%s] ListAccounts for %s", instance, cardOwner.OwnerId)
			for _, cd := range cards {
				card := pb.Card{
					OwnerId:   cardOwner.OwnerId,
					CardId: cd,
				}
				if err := stream.Send(&card); err != nil {
					glog.Errorf("[%s] %s", instance, err)
					return err
				}
				sleep()
			}
		} else {
			glog.Warningf("[%s] %s card owner doesn't exist", instance, cardOwner.OwnerId)
			return fmt.Errorf("%s card owner doesn't exist", cardOwner.OwnerId)
		}
	} else {
		glog.Warningf("[%s] CardOwner must be not nil", instance)
		return fmt.Errorf("CardOwner must be not nil")
	}
	return nil
}

func (s *server) ListMovements(card *pb.Card, stream pb.CardsService_ListMovementsServer) error {
	if card != nil {
		if movs, exist := s.masterMovements[card.OwnerId][card.CardId]; exist {
			glog.Infof("[%s] ListMovements for %s, %s", instance, card.OwnerId, card.CardId)
			for _, m := range movs {
				mov := pb.Movement{
					OwnerId:     m.OwnerId,
					CardId:   m.CardId,
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
			glog.Warningf("[%s] %s, %s card doesn't exist", instance, card.OwnerId, card.CardId)
			return fmt.Errorf("%s, %s card doesn't exist", card.OwnerId, card.CardId)
		}
	} else {
		glog.Warningf("[%s] Card must be not nil", instance)
		return fmt.Errorf("Card must be not nil")
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
	pb.RegisterCardsServiceServer(s, newServer())
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("[%s] failed to serve: %v", instance, err)
	}
}

