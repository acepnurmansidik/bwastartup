package transaction

import (
	"gorm.io/gorm"
)

type Repository interface{
	GetCampaignByID(CampaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	GetByID(ID int) (Transaction, error)
	Save(transaction Transaction)(Transaction, error)
	Update(transaction Transaction) (Transaction, error)

}

type repository struct{
	db *gorm.DB
}

func NewRepository (db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetCampaignByID(CampaignID int) ([]Transaction, error){
	var transactions []Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", CampaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error){
	var transactions []Transaction

	// karena transaction tidak punya relasi ke campaign images, dan hanya campaign yang punya relasinya
	// utk itu load Campaign beserta campaign images
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error

	if err != nil {
		return []Transaction{}, err
	}

	return transactions, nil
}

func (r *repository) Save(transaction Transaction)(Transaction, error){
	err := r.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error){
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) GetByID(ID int) (Transaction, error) {
	var transaction Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}