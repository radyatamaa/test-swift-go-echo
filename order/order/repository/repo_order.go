package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/models"
	"github.com/order/order"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type orderRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewOrderRepository(Conn *sql.DB) order.Repository {
	return &orderRepository{Conn}
}
func (m *orderRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Order, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Order, 0)
	for rows.Next() {
		t := new(models.Order)
		err = rows.Scan(
			&t.Id,
			&t.CreatedBy,
			&t.CreatedDate,
			&t.ModifiedBy,
			&t.ModifiedDate,
			&t.DeletedBy,
			&t.DeletedDate,
			&t.IsDeleted,
			&t.IsActive,
			&t.ReferenceNumber ,
			&t.CustomerName 	,
			&t.SourceAddress 	,
			&t.DestAddress 	,
			&t.Status 			,
			&t.TotalPrice 		,
			&t.CustomerReceived ,
			&t.Remarks,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *orderRepository) Insert(ctx context.Context, a *models.Order) error {
	query := `INSERT orders SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	reference_number=?,customer_name=?,source_address=?,dest_address=?,status=?,total_price=?,customer_received=?,remarks=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,
		a.ReferenceNumber, a.CustomerName, a.SourceAddress, a.DestAddress, a.Status, a.TotalPrice,a.CustomerReceived,a.Remarks)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//a.Id = lastID
	return nil
}

func (m *orderRepository) Count(ctx context.Context,referenceId string) (int, error) {
	query := `SELECT count(*) AS count FROM orders WHERE is_deleted = 0 and is_active = 1`
	if referenceId != "" {
		query = query + ` AND reference_number = '` + referenceId + `' `
	}
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (m *orderRepository) List(ctx context.Context, limit, offset int,referenceId string) ([]*models.Order, error) {
	query := `SELECT * FROM orders WHERE is_deleted = 0 and is_active = 1 `
	if referenceId != "" {
		query = query + ` AND reference_number = '` + referenceId + `' `
	}
	query = query + ` ORDER BY created_date DESC `
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
