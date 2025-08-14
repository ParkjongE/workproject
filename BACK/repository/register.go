package repository

import (
	"BACK/config"
	"BACK/model"
	"database/sql"
)

// RegisterRepository 인터페이스
type RegisterRepository interface {
	FindUserByUserid(userid string) (*model.User, error)
	CreateUser(user *model.User) error
}

// RegisterRepositoryImpl 구현체
type RegisterRepositoryImpl struct {
	db config.Database
}

// NewRegisterRepository 생성자
func NewRegisterRepository(db config.Database) RegisterRepository {
	return &RegisterRepositoryImpl{db: db}
}

// FindUserByUserid 사용자명으로 사용자 찾기
func (r *RegisterRepositoryImpl) FindUserByUserid(userid string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow("SELECT id, userid, name, role, department FROM users WHERE userid = ?", userid).
		Scan(&user.ID, &user.Userid, &user.Name, &user.Role, &user.Department)

	if err == sql.ErrNoRows {
		return nil, nil // 사용자 없음
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser 새 사용자 생성
func (r *RegisterRepositoryImpl) CreateUser(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (userid, password, salt, name, role, department, hire_date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Userid, user.Password, user.Salt, user.Name, user.Role, user.Department, user.HireDate)

	return err
}
