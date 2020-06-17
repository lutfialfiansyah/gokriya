package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/domain"
)

type mysqlUsersRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewMysqlUsersRepository(Conn *sql.DB) *mysqlUsersRepository {
	return &mysqlUsersRepository{Conn}
}

func (m *mysqlUsersRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Users, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Users, 0)
	for rows.Next() {
		t := domain.Users{}
		err = rows.Scan(
			&t.ID,
			&t.Data,
			&t.RoleId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlUsersRepository) fetchJoinRole(ctx context.Context, query string, args ...interface{}) (result []domain.UsersJoinRole, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.UsersJoinRole, 0)
	for rows.Next() {
		t := domain.UsersJoinRole{}
		err = rows.Scan(
			&t.Id,
			&t.DataUser,
			&t.DataRole,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlUsersRepository) Fetch(ctx context.Context, page int,size int) (res []domain.Users, err error) {
	query := `SELECT id,data, role_id FROM users`

	if size != 0 {
		query = query + ` LIMIT ` + strconv.Itoa(size) + ` OFFSET ` + strconv.Itoa(page)
	}
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil,  err
	}

	//if len(res) == int(num) {
	//	nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	//}

	return res,nil
}
func (m *mysqlUsersRepository) GetByID(ctx context.Context, id string) (res domain.UsersJoinRole, err error) {
	query := `select u.id,u.data as data_user , r.data as data_role from users u 
				join roles r on u.role_id = r.id
				where u.id = '` + id + `'`
	list, err := m.fetchJoinRole(ctx, query)
	if err != nil {
		return domain.UsersJoinRole{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlUsersRepository) GetByData(ctx context.Context, title string) (res domain.Users, err error) {
	query := `SELECT id,data,role_id FROM users WHERE data = ?`
	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlUsersRepository) Store(ctx context.Context, a *domain.Users) (err error) {
	query := `INSERT  users SET data=? , role_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Data, a.RoleId)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = string(lastID)
	return
}

func (m *mysqlUsersRepository) Delete(ctx context.Context, id string) (err error) {
	query := `DELETE FROM users` +
		` where id = '` + id + `'`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return
	}

	//rowsAfected, err := res.RowsAffected()
	//if err != nil {
	//	return
	//}
	//
	//if rowsAfected != 1 {
	//	err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
	//	return
	//}

	return
}
func (m *mysqlUsersRepository) Update(ctx context.Context, ar *domain.Users) (err error) {
	query := `UPDATE users set data=?, role=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Data, ar.RoleId, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
