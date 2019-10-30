package app

import (
	"context"
	"io"
)

type Server struct {
}

func (s *Server) StreamPrices(stream MM_StreamPricesServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		in.GetCcypair()
	}
}

func (s *Server) StreamCommands(MM_StreamCommandsServer) error {
	panic("implement me")
}

func (s *Server) PushExecutionReport(context.Context, *ExecutionReport) (*Confirmation, error) {
	panic("implement me")
}

func (s *Server) GetConfig(context.Context, *Setup) (*Config, error) {
	panic("implement me")
}

func (s *Server) StreamBook(MM_StreamBookServer) error {
	panic("implement me")
}

func (s *Server) PushBalance(MM_PushBalanceServer) error {
	panic("implement me")
}
