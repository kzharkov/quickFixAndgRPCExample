package adapter

type BitstampAdapter struct {
}

func NewBitstampAdapter() *BitstampAdapter {
	return &BitstampAdapter{}
}

func (b *BitstampAdapter) StreamPriceAdapter() {
	panic("implement me")
}

func (b *BitstampAdapter) StreamCommandsAdapter() {
	panic("implement me")
}

func (b *BitstampAdapter) PushExecutionReportAdapter() {
	panic("implement me")
}

func (b *BitstampAdapter) GetConfigAdapter() {
	panic("implement me")
}

func (b *BitstampAdapter) PushBalanceAdapter() {
	panic("implement me")
}

func (b *BitstampAdapter) StreamBookAdapter() {
	panic("implement me")
}
