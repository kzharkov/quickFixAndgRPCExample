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
	app "quickFix/internal"
	"quickFix/internal/adapter"
)

type Application struct {
	*quickfix.MessageRouter
	adapter.Adapter
	server *app.Server
}

func NewApplication(adp adapter.Adapter) *Application {
	application := &Application{
		MessageRouter: quickfix.NewMessageRouter(),
		Adapter:       adp,
	}

	application.AddRoute(marketdatarequest.Route(application.onMarketDataRequest))
	application.AddRoute(newordersingle.Route(application.onNewOrderSingle))
	application.AddRoute(ordercancelrequest.Route(application.onOrderCancelRequest))
	application.AddRoute(ordercancelreplacerequest.Route(application.onOrderCancelReplaceRequest))
	application.AddRoute(massquote.Route(application.onMassQuote))
	application.AddRoute(quotecancel.Route(application.onQuoteCancel))
	application.AddRoute(businessmessagereject.Route(application.onBusinessMessageReject))
	application.AddRoute(securitystatusrequest.Route(application.onSecurityStatusRequest))
	application.AddRoute(securitydefinitionrequest.Route(application.onSecurityDefinitionStatus))

	return application
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

func (a *Application) onMarketDataRequest(msg marketdatarequest.MarketDataRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onNewOrderSingle(msg newordersingle.NewOrderSingle, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onOrderCancelRequest(msg ordercancelrequest.OrderCancelRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onOrderCancelReplaceRequest(msg ordercancelreplacerequest.OrderCancelReplaceRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onMassQuote(msg massquote.MassQuote, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onQuoteCancel(msg quotecancel.QuoteCancel, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onBusinessMessageReject(msg businessmessagereject.BusinessMessageReject, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onSecurityStatusRequest(msg securitystatusrequest.SecurityStatusRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func (a *Application) onSecurityDefinitionStatus(msg securitydefinitionrequest.SecurityDefinitionRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}
