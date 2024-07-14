package adapter

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	s "github.com/core-go/sql"

	"go-service/internal/user/model"
)

func NewUserAdapter(db *sql.DB) (*UserAdapter, error) {
	userType := reflect.TypeOf(model.User{})
	fieldsIndex, _, jsonColumnMap, keys, _, _, buildParam, _, err := s.Init(userType, db)
	if err != nil {
		return nil, err
	}
	return &UserAdapter{DB: db, Map: fieldsIndex, Keys: keys, JsonColumnMap: jsonColumnMap, BuildParam: buildParam}, nil
}

type UserAdapter struct {
	DB            *sql.DB
	Map           map[string]int
	Keys          []string
	JsonColumnMap map[string]string
	BuildParam    func(int) string
}

func (r *UserAdapter) All(ctx context.Context) ([]model.User, error) {
	query := `
		select
			id, 
			username,
			email,
			phone,
			date_of_birth
		from users`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.DateOfBirth)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (r *UserAdapter) Load(ctx context.Context, id string) (*model.User, error) {
	query := `
		select
			id, 
			username,
			email,
			phone,
			date_of_birth
		from users where id = $1`
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.DateOfBirth)
		return &user, nil
	}
	return nil, nil
}

func (r *UserAdapter) Create(ctx context.Context, user *model.User) (int64, error) {
	query := `
		insert into users (
			id,
			username,
			email,
			phone,
			date_of_birth)
		values (
			$1,
			$2,
			$3, 
			$4,
			$5)`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx,
		user.Id,
		user.Username,
		user.Email,
		user.Phone,
		user.DateOfBirth)
	if err != nil {
		if strings.Index(err.Error(), "duplicate key") >= 0 {
			return -1, nil
		}
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Update(ctx context.Context, user *model.User) (int64, error) {
	query := `
		update users 
		set
			username = $1,
			email = $2,
			phone = $3,
			date_of_birth = $4
		where id = $5`
	stmt, err := r.DB.Prepare(query)
	if err != nil {
		return -1, nil
	}
	res, err := stmt.ExecContext(ctx,
		user.Username,
		user.Email,
		user.Phone,
		user.DateOfBirth,
		user.Id)
	if err != nil {
		return -1, err
	}
	count, err := res.RowsAffected()
	if count == 0 {
		return count, nil
	}
	return count, err
}

func (r *UserAdapter) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	colMap := s.JSONToColumns(user, r.JsonColumnMap)
	query, args := s.BuildToPatch("users", colMap, r.Keys, r.BuildParam)
	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *UserAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = $1"
	res, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
