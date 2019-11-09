package adapter

import (
	"context"
	"google.golang.org/grpc"
)

type SCgRPC struct {
	client  *mMClient
	server  *mMStreamBookClient
	adapter Adapter
}

func NewClientGRPC(address string) (*SCgRPC, error) {
	conn, err := grpc.Dial(address)
	if err != nil {
		return nil, err
	}
	_mMClient := &mMClient{conn}
	streamBook, err := _mMClient.StreamBook(context.Background())
	if err != nil {
		return nil, err
	}
	return &SCgRPC{
		client: _mMClient,
		server: &mMStreamBookClient{streamBook},
	}, nil
}

func (c *SCgRPC) DoStreamCommand() error {
	streamCommand, err := c.client.StreamCommands(context.Background())
	if err != nil {
		return err
	}
	for {
		command, err := streamCommand.Recv()
		if err != nil {
			return err
		}
		err = c.adapter.SendMDR(command)
		if err != nil {
			return err
		}
	}
}

func (c *SCgRPC) SendBookMsg(market *Market) error {
	return c.server.Send(market)
}
