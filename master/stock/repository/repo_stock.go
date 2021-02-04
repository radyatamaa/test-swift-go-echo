package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/master/stock"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type stockRepository struct {
	Conn *sql.DB
}


// NewuserRepository will create an object that represent the article.repository interface
func NewStockRepository(Conn *sql.DB) stock.Repository {
	return &stockRepository{Conn}
}
func (m *stockRepository) DeleteInOutBound(ctx context.Context, inboundId string,outboundId string, deleted_by string) error {
	query := ``
	if inboundId != ""{
		query = query + `UPDATE stocks SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE inbound_id = '`+inboundId + `' `
	}else if outboundId != ""{
		query = query + `UPDATE stocks,outbounds SET stocks.deleted_by=? ,stocks. deleted_date=? , stocks.is_deleted=? , stocks.is_active=? 
				WHERE (outbounds.id = stocks.outbound_id) AND outbounds.reference_number = '`+outboundId + `' `
	}
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(), 1, 0)
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

func (m *stockRepository) fetchJoinProductInOutBound(ctx context.Context, query string, args ...interface{}) ([]*models.StockJoinProductInOutbound, error) {
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

	result := make([]*models.StockJoinProductInOutbound, 0)
	for rows.Next() {
		t := new(models.StockJoinProductInOutbound)
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
			&t.ProductId,
			&t.InboundId ,
			&t.OutboundId 	,
			&t.CurrentStock,
			&t.ProductSKU 		,
			&t.ProductName 	,
			&t.InboundDate			,
			&t.InboundQTY 		,
			&t.OutboundDate		,
			&t.OutboundQTY 			,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *stockRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Stock, error) {
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

	result := make([]*models.Stock, 0)
	for rows.Next() {
		t := new(models.Stock)
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
			&t.ProductId,
			&t.InboundId ,
			&t.OutboundId 	,
			&t.CurrentStock,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *stockRepository) GetFirst(ctx context.Context, productId string) (res *models.Stock,err error) {
	query := `SELECT * FROM stocks WHERE  is_active = 1 AND is_deleted = 0 `

	if productId != "" {
		query = query + ` AND product_id = '` + productId + `' `
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

func (m *stockRepository) Insert(ctx context.Context, a *models.Stock) error {
	query := `INSERT stocks SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? ,
	product_id=?,inbound_id=?,outbound_id=?,current_stock=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1,
		a.ProductId,a.InboundId,a.OutboundId,a.CurrentStock)
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

func (m *stockRepository) Count(ctx context.Context,productId string,bound int) (int, error) {
	query := `SELECT count(*) as count 
				FROM stocks s
				LEFT JOIN inbounds i ON i.id = s.inbound_id
				LEFT JOIN outbounds o ON i.id = s.outbound_id
				WHERE s.is_deleted = 0 and s.is_active = 1 `

	if bound == 1{
		query = query + ` AND s.inbound_id is not null `
	}else if bound == 2{
		query = query + ` AND s.outbound_id is not nul `
	}
	if productId != ""{
		query = query + ` AND s.product_id = '` + productId + `' `
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

func (m *stockRepository) List(ctx context.Context, limit, offset int,productId string,bound int) ([]*models.StockJoinProductInOutbound, error) {
	query := `SELECT s.*,
					p.sku as product_sku ,
					p.name as product_name ,
					i.created_date as inbound_date,
					i.jumlah as  inbound_qty,
					o.created_date as outbound_date,
					o.qty as outbound_qty
				FROM stocks s
				JOIN products p ON p.id = s.product_id
				LEFT JOIN inbounds i ON i.id = s.inbound_id
				LEFT JOIN outbounds o ON i.id = s.outbound_id
				WHERE s.is_deleted = 0 and s.is_active = 1 `

	if bound == 1{
		query = query + ` AND s.inbound_id is not null `
	}else if bound == 2{
		query = query + ` AND s.outbound_id is not null `
	}
	if productId != ""{
		query = query + ` AND s.product_id = '` + productId + `' `
	}
	query = query + ` order by created_date desc `
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetchJoinProductInOutBound(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}
