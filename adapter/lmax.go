package adapter

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/tag"
	"sync"
)

type LmaxAdapter struct {
	fix          quickfix.Application
	gRPC         *SCgRPC
	sessionID    quickfix.SessionID
	clients      map[string]string
	mutexClients *sync.RWMutex
}

func NewLmaxAdapter(gRPC *SCgRPC) *LmaxAdapter {
	return &LmaxAdapter{
		mutexClients: new(sync.RWMutex),
		clients:      make(map[string]string),
		gRPC:         gRPC,
	}
}

func (b *LmaxAdapter) StreamBookAdapter(market *Market) error {
	return b.gRPC.SendBookMsg(market)
}

func (b *LmaxAdapter) StreamCommandsAdapter() {
	panic("implement me")
}

func (b *LmaxAdapter) PushExecutionReportAdapter() {
	panic("implement me")
}

func (b *LmaxAdapter) GetConfigAdapter() {
	panic("implement me")
}

func (b *LmaxAdapter) PushBalanceAdapter() {
	panic("implement me")
}

func (b *LmaxAdapter) SendMDR(command *Command) error {
	switch command.Command {
	case "md":
		if command.Action == "subscribe" {
			m := quickfix.NewMessage()
			m.Body.SetInt(tag.SubscriptionRequestType, 1)
			m.Body.SetInt(tag.MarketDepth, 0)
			m.Body.SetInt(tag.MDUpdateType, 0)
			m.Body.SetString(tag.Symbol, command.Item)
			err := quickfix.Send(m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *LmaxAdapter) GetClient(client string) (string, bool) {
	b.mutexClients.RLock()
	defer b.mutexClients.RUnlock()
	result, ok := b.clients[client]
	return result, ok
}

func (b *LmaxAdapter) SetApplication(application quickfix.Application) {
	b.fix = application
}
