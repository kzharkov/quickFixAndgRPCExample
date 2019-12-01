package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"quickFix/adapter"
	"time"
)

type Server struct {

}

func (s *Server) StreamPrices(adapter.MM_StreamPricesServer) error {
	panic("implement me")
}

func (s *Server) StreamCommands(streamCommand adapter.MM_StreamCommandsServer) error {
	for {
		err := streamCommand.Send(&adapter.Command{
			Command: "md",
			Action: "subscribe",
			Item: "BTC/USD",
		})
		if err != nil {
			log.Println(err)
		}
		time.Sleep(time.Minute*5)
	}
}

func (s *Server) PushExecutionReport(context.Context, *adapter.ExecutionReport) (*adapter.Confirmation, error) {
	panic("implement me")
}

func (s *Server) GetConfig(context.Context, *adapter.Setup) (*adapter.Config, error) {
	panic("implement me")
}

func (s *Server) StreamBook(streamBook adapter.MM_StreamBookServer) error {
	for {
		market, err := streamBook.Recv()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
		}
		log.Println(market)
	}
}

func (s *Server) PushBalance(adapter.MM_PushBalanceServer) error {
	panic("implement me")
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:5000"))
	if err != nil {
		log.Println(err)
		return
	}
	server := &Server{}
	gRpcServer := grpc.NewServer()
	adapter.RegisterMMServer(gRpcServer, server)
	if err = gRpcServer.Serve(lis); err != nil {
		log.Println(err)
		return
	}
}