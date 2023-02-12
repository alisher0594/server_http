package data

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Devices interface {
	Insert(ctx context.Context, device *Device) error
	Get(ctx context.Context, id int64) (*Device, error)
	GetAll(ctx context.Context, name string, value int, filters Filters) ([]*Device, Metadata, error)
	Update(ctx context.Context, device *Device) error
	Delete(ctx context.Context, id int64) error
}

type Models struct {
	Devices
}

func NewModels(db *sql.DB) Models {
	return Models{
		Devices: DeviceModel{DB: db},
	}
}
