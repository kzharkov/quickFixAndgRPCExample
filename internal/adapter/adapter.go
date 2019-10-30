package adapter

type Adapter interface {
	StreamPriceAdapter()
	StreamCommandsAdapter()
	PushExecutionReportAdapter()
	GetConfigAdapter()
	PushBalanceAdapter()
	StreamBookAdapter()
}
