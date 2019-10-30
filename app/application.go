package app

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/fix43/securitystatusrequest"
	"github.com/quickfixgo/quickfix/fix44/businessmessagereject"
	"github.com/quickfixgo/quickfix/fix44/marketdatarequest"
	"github.com/quickfixgo/quickfix/fix44/massquote"
	"github.com/quickfixgo/quickfix/fix44/newordersingle"
	"github.com/quickfixgo/quickfix/fix44/ordercancelreplacerequest"
	"github.com/quickfixgo/quickfix/fix44/ordercancelrequest"
	"github.com/quickfixgo/quickfix/fix44/quotecancel"
	"github.com/quickfixgo/quickfix/fix44/securitydefinitionrequest"
	"quickFix/internal/adapter"
)

type Application struct {
	*quickfix.MessageRouter
	adapter.Adapter
}

func NewApplication(adp adapter.Adapter) *Application {
	app := &Application{
		MessageRouter: quickfix.NewMessageRouter(),
		Adapter:       adp,
	}

	app.AddRoute(marketdatarequest.Route(app.onMarketDataRequest))
	app.AddRoute(newordersingle.Route(app.onNewOrderSingle))
	app.AddRoute(ordercancelrequest.Route())
	app.AddRoute(ordercancelreplacerequest.Route())
	app.AddRoute(massquote.Route())
	app.AddRoute(quotecancel.Route())
	app.AddRoute(businessmessagereject.Route())
	app.AddRoute(securitystatusrequest.Route())
	app.AddRoute(securitydefinitionrequest.Route())

	return app
}

//OnCreate implemented as part of Application interface
func (a Application) OnCreate(sessionID quickfix.SessionID) { return }

//OnLogon implemented as part of Application interface
func (a Application) OnLogon(sessionID quickfix.SessionID) { return }

//OnLogout implemented as part of Application interface
func (a Application) OnLogout(sessionID quickfix.SessionID) { return }

//ToAdmin implemented as part of Application interface
func (a Application) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) { return }

//ToApp implemented as part of Application interface
func (a Application) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}

//FromAdmin implemented as part of Application interface
func (a Application) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

//FromApp implemented as part of Application interface, uses Router on incoming application messages
func (a *Application) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return a.Route(msg, sessionID)
}

func (a *Application) onMarketDataRequest(msg marketdatarequest.MarketDataRequest, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {
	// TODO: Определяем тип сообщения
	return
}

func (a *Application) onNewOrderSingle(msg newordersingle.NewOrderSingle, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

}
