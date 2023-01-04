package product

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll() ([]Product, error)
	GetByID(id int) (Product, error)
	Create(product Product) (Product, error)
	Update(id int, inputProduct InputProduct) (Product, error)
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAll() ([]Product, error) {
	var products []Product
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) GetByID(id int) (Product, error) {
	var product Product
	err := r.db.Where(&Product{ID: id}).First(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *repository) Create(product Product) (Product, error) {
	uniq := r.db.Where(&Product{Name: product.Name}).First(&product).Error
	if uniq != nil {
		err := r.db.Create(&product).Error
		if err != nil {
			return product, err
		}
	} else {
		return product, errors.New("name product already exist")
	}

	return product, nil
}

func (r *repository) Update(id int, inputProduct InputProduct) (Product, error) {
	product, err := r.GetByID(id)
	if err != nil {
		return Product{}, err
	}
	product.Name = inputProduct.Name
	product.Price = inputProduct.Price

	r.db.Save(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(id int) error {
	product := &Product{ID: id}
	tx := r.db.Delete(&product)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
