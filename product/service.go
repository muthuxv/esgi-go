package product

type Service interface {
	FetchAll() ([]Product, error)
	FetchByID(id int) (Product, error)
	Create(inputProduct Product) (Product, error)
	Update(id int, product Product) (Product, error)
	Delete(id int) error
}

type service struct {
	repository Repository
}

func NewProductService(r Repository) *service {
	return &service{r}
}

func (s *service) FetchAll() ([]Product, error) {
	products, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *service) FetchByID(id int) (Product, error) {
	product, err := s.repository.GetByID(id)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (s *service) Create(inputProduct Product) (Product, error) {
	var product Product
	product.Name = inputProduct.Name
	product.Price = inputProduct.Price

	newProduct, err := s.repository.Create(product)
	if err != nil {
		return product, err
	}

	return newProduct, nil
}

func (s *service) Update(id int, inputProduct Product) (Product, error) {
	uProduct, err := s.repository.Update(id, inputProduct)
	if err != nil {
		return uProduct, err
	}
	return uProduct, nil
}

func (s *service) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
