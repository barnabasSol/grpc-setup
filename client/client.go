package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/barnabasSol/grpcsetup/coffeeshop_protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func RunClient() {
	time.Sleep(2 * time.Second)

	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn.Connect()
	for conn.GetState() != connectivity.Ready {
		if !conn.WaitForStateChange(ctx, conn.GetState()) {
			log.Fatalf("connection timed out, last state: %v", conn.GetState())
		}
	}
	log.Println("Connection state: Ready")

	c := pb.NewCoffeeShopClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("failed calling GetMenu: %v", err)
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
				log.Fatalf("cannot receive from stream: %v", err)
			}
			items = resp.Items
			fmt.Printf("items: %v\n", items)
		}
	}()
	<-done
	fmt.Printf("items: %v\n", items)
	receipt, err := c.PlaceOrder(ctx, &pb.Order{Items: items})
	if err != nil {
		log.Fatalf("failed calling PlaceOrder: %v", err)
	}
	fmt.Printf("receipt: %v\n", receipt)
	status, err := c.GetOrderStatus(ctx, receipt)
	if err != nil {
		log.Fatalf("failed calling GetOrderStatus: %v", err)
	}
	fmt.Printf("status: %v\n", status)
}
