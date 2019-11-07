package adapter

import (
	"context"
)

type SCgRPC struct {
	client  *mMClient
	server  *mMStreamBookClient
	adapter Adapter
}

func NewClientGRPC() *SCgRPC {
	return &SCgRPC{client: &mMClient{}}
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
