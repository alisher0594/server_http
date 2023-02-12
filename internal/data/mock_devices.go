package data

import "context"

type MockDeviceModel struct{}

func (m MockDeviceModel) Insert(ctx context.Context, device *Device) error {
	panic("implement me")
}

func (m MockDeviceModel) Get(ctx context.Context, id int64) (*Device, error) {
	panic("implement me")

}

func (m MockDeviceModel) GetAll(ctx context.Context, name string, value int, filters Filters) ([]*Device, Metadata, error) {
	panic("implement me")
}

func (m MockDeviceModel) Update(ctx context.Context, device *Device) error {
	panic("implement me")

}

func (m MockDeviceModel) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func NewMockModels() Models {
	return Models{
		Devices: MockDeviceModel{},
	}
}
