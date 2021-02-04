package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/master/outbound"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type outboundRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewOutboundRepository(Conn *sql.DB) outbound.Repository {
	return &outboundRepository{Conn}
}
func (m *outboundRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Outbound, error) {
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

	result := make([]*models.Outbound, 0)
	for rows.Next() {
		t := new(models.Outbound)
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
			&t.TimeStamp ,
			&t.ProductId 		,
			&t.Total 			,
			&t.Usecase 		,
			&t.ReferenceNumber 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *outboundRepository) fetchJoinProduct(ctx context.Context, query string, args ...interface{}) ([]*models.OutboundJoinProduct, error) {
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

	result := make([]*models.OutboundJoinProduct, 0)
	for rows.Next() {
		t := new(models.OutboundJoinProduct)
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
			&t.TimeStamp ,
			&t.ProductId 		,
			&t.Total 			,
			&t.Usecase 		,
			&t.ReferenceNumber 	,
			&t.ProductSKU,
			&t.ProductName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *outboundRepository) GetByID(ctx context.Context, id string) (res *models.Outbound, err error) {
	query := `SELECT * FROM outbounds WHERE  is_active = 1 AND is_deleted = 0 `

	if id != "" {
		query = query + ` AND id = '` + id + `' `
	}

	list, err := m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *outboundRepository) Update(ctx context.Context, a *models.Outbound) error {
	query := `UPDATE outbounds set modified_by=?, modified_date=? , 
		time_stamp=?,product_id=?,total=?,usecase=?,reference_number=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(),
		a.TimeStamp,a.ProductId,a.Total,a.Usecase,a.ReferenceNumber, a.Id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

func (m *outboundRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE outbounds SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(), 1, 0, id)
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

func (m *outboundRepository) Insert(ctx context.Context, a *models.Outbound) error {
	query := `INSERT outbounds SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	time_stamp=?,product_id=?,total=?,usecase=?,reference_number=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,
		a.TimeStamp,a.ProductId,a.Total,a.Usecase,a.ReferenceNumber)
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

func (m *outboundRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) as count FROM
				outbounds i
				JOIN products p ON p.id = i.product_id
				WHERE i.is_active = 1 AND i.is_deleted = 0 AND p.is_active = 1 AND p.is_deleted = 0`

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

func (m *outboundRepository) List(ctx context.Context, limit, offset int) ([]*models.OutboundJoinProduct, error) {
	query := `SELECT i.*,p.sku as product_sku , p.name as product_name FROM
				outbounds i
				JOIN products p ON p.id = i.product_id
				WHERE i.is_active = 1 AND i.is_deleted = 0 AND p.is_active = 1 AND p.is_deleted = 0`
	query = query + ` order by created_date desc`
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetchJoinProduct(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
