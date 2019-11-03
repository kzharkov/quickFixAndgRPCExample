package adapter

import app "quickFix/internal"

type Adapter interface {
	StreamBookAdapter(market *app.Market) error
	StreamCommandsAdapter()
	PushExecutionReportAdapter()
	GetConfigAdapter()
	PushBalanceAdapter()
	SendMDR(command *app.Command) error
	GetClient(string) (string, bool)
}
