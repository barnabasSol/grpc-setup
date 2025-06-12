package main

import (
	"context"
	"log"
	"net"

	pb "github.com/barnabasSol/grpcsetup/coffeeshop_protos"
	grpc "google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

func (s *server) GetMenu(mr *pb.MenuRequest, streamSrv grpc.ServerStreamingServer[pb.Menu]) error {
	items := []*pb.Item{
		{Id: "1", Name: "Black Coffee"},
		{Id: "2", Name: "Americano"},
		{Id: "3", Name: "Iced Coffee"},
	}
	for i := range items {
		streamSrv.Send(&pb.Menu{
			Items: items[0 : i+1],
		})
	}
	return nil
}
func (s *server) PlaceOrder(ctx context.Context, order *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "abc123",
	}, nil
}
func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN_PROGRESS",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "9001")
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
