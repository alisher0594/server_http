package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alisher0594/validator/pkg/validator"
	"log"
	"time"
)

type Device struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Value       int       `json:"value"`
	Description string    `json:"description"`
	Version     int32     `json:"version"`
}

func (d *Device) Validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(d.Value != 0, "value", "must be provided")
	v.Check(d.Description != "", "description", "must be provided")
}

type DeviceModel struct {
	DB *sql.DB
}

func (m DeviceModel) Insert(ctx context.Context, device *Device) error {
	query := `
        INSERT INTO devices (name, value, description)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	args := []interface{}{device.Name, device.Value, device.Description}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&device.ID, &device.CreatedAt, &device.Version)
}

func (m DeviceModel) Get(ctx context.Context, id int64) (*Device, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, name, value, description, version
        FROM devices
        WHERE id = $1`

	var device Device
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&device.ID,
		&device.CreatedAt,
		&device.Name,
		&device.Value,
		&device.Description,
		&device.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &device, nil
}

func (m DeviceModel) GetAll(ctx context.Context, name string, value int, filters Filters) ([]*Device, Metadata, error) {
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, name, value, descripton, version
        FROM devices
        WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
        AND (value >= $2 OR $2 = 0)     
        ORDER BY %s %s, id ASC
        LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(filters.limit(), filters.offset())

	args := []interface{}{name, value, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	movies := make([]*Device, 0, filters.limit())

	for rows.Next() {
		var movie Device

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Name,
			&movie.Value,
			&movie.Description,
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

func (m DeviceModel) Update(ctx context.Context, movie *Device) error {
	query := `
        UPDATE devices
        SET name = $1, value = $2, description = $3, version = version + 1
        WHERE id = $4 AND version = $5
        RETURNING version`

	args := []interface{}{
		movie.Name,
		movie.Value,
		movie.Description,
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m DeviceModel) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM devices
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
