package transaction

import "gorm.io/gorm"

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	Save(transaction Transaction) (Transaction, error)   //save transaksi `amount`
	Update(transaction Transaction) (Transaction, error) //update paymentURL
	GetByID(ID int) (Transaction, error)                 //mengambil data id yg bersangkutan
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transactions []Transaction
	// order by id desc untuk mengurutkan dari besar ke kecil
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error) {

	var transactions []Transaction

	// load CampaignImages tanpa relasi langsung (perantara dari Campaign) dan memberi batasan field : (is_primary = 1)
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil

}

// todo save new transaction
func (r *repository) Save(transaction Transaction) (Transaction, error) {

	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil

}

// todo update paymentURL
func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// todo data ID
func (r *repository) GetByID(ID int) (Transaction, error) {

	var transaction Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
