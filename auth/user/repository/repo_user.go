package repository

import (
	"context"
	"database/sql"
	"github.com/auth/user"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type userRepository struct {
	Conn *sql.DB
}



// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewUserRepository(Conn *sql.DB) user.Repository {
	return &userRepository{Conn}
}

func (m *userRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.User, error) {
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

	result := make([]*models.User, 0)
	for rows.Next() {
		t := new(models.User)
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
			&t.UserEmail,
			&t.Password,
			&t.Phone,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *userRepository) Create(ctx context.Context, a models.User) (*string, error) {
	query := `INSERT users SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , user_email=?,password=?,phone=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.UserEmail, a.Password, a.Phone)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}
func (m *userRepository) ValidateUser(ctx context.Context, email string) (res *models.User,err error) {
	query := `SELECT * FROM users WHERE is_active = 1 AND is_deleted = 0 AND user_email = ?`

	list, err := m.fetch(ctx, query, email)
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

