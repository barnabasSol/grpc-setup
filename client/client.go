package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/barnabasSol/grpcsetup/coffeeshop_protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunClient() {
	time.Sleep(3 * time.Second)
	conn, err := grpc.NewClient(
		"localhost:5001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to gRPC server %v", err)
	}
	defer conn.Close()
	c := pb.NewCoffeeShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("failed calling GetMenu %v", err)
	}
	done := make(chan bool)
	var items []*pb.Item
	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive from stream %v", err)
			}
			items = resp.Items
		}
	}()
	<-done
	fmt.Printf("items: %v\n", items)
	receipt, _ := c.PlaceOrder(ctx, &pb.Order{Items: items})
	fmt.Printf("receipt: %v\n", receipt)
	status, _ := c.GetOrderStatus(ctx, receipt)
	fmt.Printf("status: %v\n", status)

}
