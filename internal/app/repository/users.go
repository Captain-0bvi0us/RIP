package repository

import (
	"RIP/internal/app/ds"

	"golang.org/x/crypto/bcrypt"
)

// Методы для работы с пользователями

func (r *Repository) CreateUser(user *ds.Users) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByUsername(username string) (*ds.Users, error) {
	var user ds.Users
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*ds.Users, error) {
	var user ds.Users
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(id uint, req ds.UserUpdateRequest) error {
	updates := make(map[string]interface{})

	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Password != nil {
		// Хешируем новый пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password"] = string(hashedPassword)
	}
	if req.Moderator != nil {
		updates["moderator"] = *req.Moderator
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&ds.Users{}).Where("id = ?", id).Updates(updates).Error
}
