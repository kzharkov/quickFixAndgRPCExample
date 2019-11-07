package app

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/quickfixgo/quickfix/fix43/securitystatusrequest"
	"github.com/quickfixgo/quickfix/fix44/businessmessagereject"
	"github.com/quickfixgo/quickfix/fix44/marketdatarequest"
	"github.com/quickfixgo/quickfix/fix44/massquote"
	"github.com/quickfixgo/quickfix/fix44/newordersingle"
	"github.com/quickfixgo/quickfix/fix44/ordercancelreplacerequest"
	"github.com/quickfixgo/quickfix/fix44/ordercancelrequest"
	"github.com/quickfixgo/quickfix/fix44/quotecancel"
	"github.com/quickfixgo/quickfix/fix44/securitydefinitionrequest"
	"github.com/quickfixgo/quickfix/tag"
	"quickFix/adapter"
)

/*
В этом файле основная логика для FIX
*/

type Application struct {
	*quickfix.MessageRouter
	adapter.Adapter
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
	msgType, rejErr := msg.GetMsgType()
	if rejErr != nil {
		return rejErr
	}
	switch msgType {
	case enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH:
		CCYPair, rejErr := msg.Body.GetString(tag.Symbol)
		if rejErr != nil {
			return rejErr
		}
		lenNoMDEntryTypes, rejErr := msg.Body.GetInt(tag.NoMDEntryTypes)
		if rejErr != nil {
			return rejErr
		}

		var bids []*adapter.BookEntry
		var asks []*adapter.BookEntry

		var totalVolAsk float64
		var totalVolBid float64
		var oldVol float64
		var volLst float64
		var lastPrice float64

		noMDEntryTypesRepeatingGroup, rejErr := msg.GetNoMDEntryTypes()
		if rejErr != nil {
			return rejErr
		}
		for i := 0; i < lenNoMDEntryTypes; i++ {
			noMDEntryTypes := noMDEntryTypesRepeatingGroup.Get(i)
			MDEntryType, rejErr := noMDEntryTypes.GetMDEntryType()
			if rejErr != nil {
				return rejErr
			}
			switch MDEntryType {
			case enum.MDEntryType_BID:
				qtyTypeField := field.QuantityField{}
				rejErr := noMDEntryTypes.GetField(tag.MDEntrySize, &qtyTypeField)
				if rejErr != nil {
					return rejErr
				}
				mDEntrySize, ok := qtyTypeField.Float64()
				if !ok {
					continue
				}
				totalVolBid += mDEntrySize

				rejErr = noMDEntryTypes.GetField(tag.MDEntryPx, &qtyTypeField)
				if rejErr != nil {
					return rejErr
				}
				mDEntryPx, ok := qtyTypeField.Float64()
				if !ok {
					continue
				}
				lastPrice = mDEntryPx

				bids = append(bids, &adapter.BookEntry{
					Type:   "bid",
					Price:  lastPrice,
					Amount: volLst,
					Total:  oldVol,
				})
				volLst = mDEntrySize
				oldVol = totalVolBid

			case enum.MDEntryType_OFFER:
				qtyTypeField := field.QuantityField{}
				rejErr := noMDEntryTypes.GetField(tag.MDEntrySize, &qtyTypeField)
				if rejErr != nil {
					return rejErr
				}
				mDEntrySize, ok := qtyTypeField.Float64()
				if !ok {
					continue
				}
				totalVolAsk += mDEntrySize

				rejErr = noMDEntryTypes.GetField(tag.MDEntryPx, &qtyTypeField)
				if rejErr != nil {
					return rejErr
				}
				mDEntryPx, ok := qtyTypeField.Float64()
				if !ok {
					continue
				}

				asks = append(asks, &adapter.BookEntry{
					Type:   "ask",
					Price:  mDEntryPx,
					Amount: mDEntrySize,
					Total:  totalVolAsk,
				})
			}
		}

		if lastPrice != 0 {
			bids = append(bids, &adapter.BookEntry{
				Type:   "bid",
				Price:  lastPrice * (1 - 0.002),
				Amount: volLst,
				Total:  oldVol,
			})
		}

		mDReqID, rejErr := msg.GetMDReqID()
		if rejErr != nil {
			return rejErr
		}

		client, ok := a.Adapter.GetClient(mDReqID)
		if !ok {
			refTagID := field.NewRefTagID(0).Tag()
			return quickfix.NewMessageRejectError("Not get client", 0, &refTagID)
		}

		market := &adapter.Market{
			Ccypair: CCYPair,
			Client:  client,
			Asks:    asks,
			Bids:    bids,
		}

		err := a.Adapter.StreamBookAdapter(market)
		if err != nil {
			refTagID := field.NewRefTagID(0).Tag()
			return quickfix.NewMessageRejectError(err.Error(), 0, &refTagID)
		}
	}
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
