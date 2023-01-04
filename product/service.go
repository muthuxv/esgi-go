package product

type Service interface {
	FetchAll() ([]Product, error)
	FetchByID(id int) (Product, error)
	Create(input InputProduct) (Product, error)
	Update(id int, inputProduct InputProduct) (Product, error)
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

func (s *service) Create(input InputProduct) (Product, error) {
	var product Product
	product.Name = input.Name
	product.Price = input.Price

	newProduct, err := s.repository.Create(product)
	if err != nil {
		return product, err
	}

	return newProduct, nil
}

func (s *service) Update(id int, input InputProduct) (Product, error) {
	uProduct, err := s.repository.Update(id, input)
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
