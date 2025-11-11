package repository

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kanban_api/internal/model"
	"time"
)

type sqliteUserRep struct {
	db *gorm.DB
}

type userRow struct {
	ID           string `gorm:"primary_key"`
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func NewSQLiteUserRepo(path string) (UserRepository, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&userRow{}); err != nil {
		return nil, err
	}
	return &sqliteUserRep{db: db}, nil
}

func (r *sqliteUserRep) toModel(row userRow) model.User {
	return model.User{
		ID:           row.ID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		CreatedAt:    row.CreatedAt,
	}
}

func (r *sqliteUserRep) Create(email, passwordHash string) (model.User, error) {
	now := time.Now()
	rw := userRow{
		ID:           generateID(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
	}
	if err := r.db.Create(&rw).Error; err != nil {
		return model.User{}, err
	}
	return r.toModel(rw), nil
}

func (r *sqliteUserRep) GetByEmail(email string) (model.User, error) {
	var rw userRow
	if err := r.db.First(&rw, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, nil
	}
	return r.toModel(rw), nil
}

func (r *sqliteUserRep) GetByID(id string) (model.User, error) {
	var rw userRow
	if err := r.db.First(&rw, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return r.toModel(rw), nil
}
