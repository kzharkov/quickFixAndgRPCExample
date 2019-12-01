package adapter

import (
	"context"
	"google.golang.org/grpc"
	"log"
)

type SCgRPC struct {
	client  MMClient
	server  MM_StreamBookClient
	adapter Adapter
}

func NewClientGRPC(address string) (*SCgRPC, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	_mMClient := NewMMClient(conn)
	streamBook, err := _mMClient.StreamBook(context.Background())
	if err != nil {
		return nil, err
	}
	return &SCgRPC{
		client: _mMClient,
		server: streamBook,
	}, nil
}

func (c *SCgRPC) DoStreamCommand() error {
	streamCommand, err := c.client.StreamCommands(context.Background())
	if err != nil {
		return err
	}
	err = streamCommand.Send(&Status{
		Status:  "",
		Client:  "",
		Message: "",
	})
	for {
		command, err := streamCommand.Recv()
		if err != nil {
			return err
		}
		log.Println(command)
		err = c.adapter.SendMDR(command)
		if err != nil {
			return err
		}
	}
}

func (c *SCgRPC) SendBookMsg(market *Market) error {
	return c.server.Send(market)
}

func (c *SCgRPC) SetAdapter(adapter Adapter) {
	c.adapter = adapter
}

func (c *SCgRPC) GetConfig() (*Config, error) {
	config, err := c.client.GetConfig(context.Background(), &Setup{Name: "LMAX"})
	if err != nil {
		return nil, err
	}

	return config, nil
}
