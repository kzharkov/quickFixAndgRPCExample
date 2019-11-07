package adapter

/*
Интерфейс адаптера, ему должны соответствовать наши объекты для конкретных бирж
*/

type Adapter interface {
	StreamBookAdapter(market *Market) error
	StreamCommandsAdapter()
	PushExecutionReportAdapter()
	GetConfigAdapter()
	PushBalanceAdapter()
	SendMDR(command *Command) error
	GetClient(string) (string, bool)
}
