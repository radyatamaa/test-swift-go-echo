package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/master/inbound"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type inboundRepository struct {
	Conn *sql.DB
}

// NewuserRepository will create an object that represent the article.repository interface
func NewInboundRepository(Conn *sql.DB) inbound.Repository {
	return &inboundRepository{Conn}
}
func (m *inboundRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.InboundJoinProduct, error) {
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

	result := make([]*models.InboundJoinProduct, 0)
	for rows.Next() {
		t := new(models.InboundJoinProduct)
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
			&t.InboundTime ,
			&t.ExpiredDate,
			&t.ProductId 	,
			&t.Jumlah 		,
			&t.	HargaBeli 	,
			&t.Total 		,
			&t.	NoPO 			,
			&t.ProductSKU 	,
			&t.ProductName 	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *inboundRepository) GetByID(ctx context.Context, id string) (res *models.InboundJoinProduct, err error) {
	query := `SELECT i.*,p.sku as product_sku , p.name as product_name FROM
				inbounds i
				JOIN products p ON p.id = i.product_id
				WHERE i.is_active = 1 AND i.is_deleted = 0 AND p.is_active = 1 AND p.is_deleted = 0`

	if id != "" {
		query = query + ` AND i.id = '` + id + `' `
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

func (m *inboundRepository) Update(ctx context.Context, a *models.Inbound) error {
	query := `UPDATE inbounds set modified_by=?, modified_date=? , 
			inbound_time=?,product_id=?,jumlah=?,harga_beli,total=?,no_po=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.InboundTime,a.ProductId,a.Jumlah,a.HargaBeli,a.Total,a.NoPO, a.Id)
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

func (m *inboundRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE inbounds SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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

func (m *inboundRepository) Insert(ctx context.Context, a *models.Inbound) error {
	query := `INSERT inbounds SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	inbound_time=?,expired_date=?,product_id=?,jumlah=?,harga_beli=?,total=?,no_po=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,
		a.InboundTime,a.ExpiredDate,a.ProductId,a.Jumlah,a.HargaBeli,a.Total,a.NoPO)
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

func (m *inboundRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*)  as count FROM
				inbounds i
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

func (m *inboundRepository) List(ctx context.Context, limit, offset int) ([]*models.InboundJoinProduct, error) {
	query := `SELECT i.*,p.sku as product_sku , p.name as product_name FROM
				inbounds i
				JOIN products p ON p.id = i.product_id
				WHERE i.is_active = 1 AND i.is_deleted = 0 AND p.is_active = 1 AND p.is_deleted = 0`
	query = query + ` order by created_date desc`
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
