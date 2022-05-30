package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository{
	return &repository{db}
}

func (r *repository) Save(user User) (User, error){
	err := r.db.Create(&user).Error

	if err != nil{
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error){
	// tampung user ke variabel
	var user User
	// cari user berdasarkan email
	err := r.db.Where("email = ?",email).Find(&user).Error
	if err != nil{
		return user, err
	}

	return user, nil
}