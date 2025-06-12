package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/barnabasSol/grpcsetup/client"
	pb "github.com/barnabasSol/grpcsetup/coffeeshop_protos"
	grpc "google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(mr *pb.MenuRequest, streamSrv grpc.ServerStreamingServer[pb.Menu]) error {
	log.Print("GOT HIT")
	items := []*pb.Item{
		{Id: "1", Name: "Black Coffee"},
		{Id: "2", Name: "Americano"},
		{Id: "3", Name: "Iced Coffee"},
	}
	for i := range items {
		if streamSrv.Context().Err() != nil {
			return streamSrv.Context().Err()
		}
		streamSrv.Send(&pb.Menu{Items: items[0 : i+1]})
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Receipt, error) {
	log.Println("place order hit")
	return &pb.Receipt{
		Id: "abc123",
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	log.Println("get order status hit")
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN_PROGRESS",
	}, nil
}

func main() {
	go client.RunClient()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	grpcSrv := grpc.NewServer()
	pb.RegisterCoffeeShopServer(grpcSrv, &server{})
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("failed to serve %s", err)
	}
	log.Println("grpc is running")
}
