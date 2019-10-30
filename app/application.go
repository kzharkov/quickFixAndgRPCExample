package app

import "github.com/quickfixgo/quickfix"

type Application struct {
	*quickfix.MessageRouter
}

func newApplication() *Application {
	app := &Application{
		MessageRouter: quickfix.NewMessageRouter(),
	}

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
