package payment

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Payment, error)
	GetByID(id int) (Payment, error)
	Create(payment Payment) (Payment, error)
	Update(id int, payment Payment) (Payment, error)
	Delete(id int) error
	Stream() (<-chan Payment, error)
}

type repository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAll() ([]Payment, error) {
	var payments []Payment
	err := r.db.Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *repository) GetByID(id int) (Payment, error) {
	var payment Payment
	err := r.db.Where(&Payment{ID: id}).First(&payment).Error
	if err != nil {
		return payment, err
	}
	return payment, nil
}

func (r *repository) Create(payment Payment) (Payment, error) {
	err := r.db.Create(&payment).Error
	if err != nil {
		return payment, err
	}
	return payment, nil
}

func (r *repository) Update(id int, inputPayment Payment) (Payment, error) {
	payment, err := r.GetByID(id)
	if err != nil {
		return Payment{}, err
	}
	payment.ProductID = inputPayment.ProductID
	payment.PricePaid = inputPayment.PricePaid

	r.db.Save(&payment)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (r *repository) Delete(id int) error {
	payment := &Payment{ID: id}
	tx := r.db.Delete(&payment)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Payment not found")
	}
	return nil
}

func (r *repository) Stream() (<-chan Payment, error) {
	return nil, nil
}
