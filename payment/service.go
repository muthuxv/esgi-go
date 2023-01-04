package payment

type Service interface {
	FetchAll() ([]Payment, error)
	FetchByID(id int) (Payment, error)
	Create(input InputPayment) (Payment, error)
	Update(id int, inputPayment InputPayment) (Payment, error)
	Delete(id int) error
}

type service struct {
	repository Repository
}

func NewPaymentService(r Repository) *service {
	return &service{r}
}

func (s *service) FetchAll() ([]Payment, error) {
	payments, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *service) FetchByID(id int) (Payment, error) {
	payment, err := s.repository.GetByID(id)
	if err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (s *service) Create(input InputPayment) (Payment, error) {
	var payment Payment
	payment.ProductID = input.ProductID
	payment.PricePaid = input.PricePaid

	newPayment, err := s.repository.Create(payment)
	if err != nil {
		return payment, err
	}

	return newPayment, nil
}

func (s *service) Update(id int, input InputPayment) (Payment, error) {
	uPayment, err := s.repository.Update(id, input)
	if err != nil {
		return uPayment, err
	}
	return uPayment, nil
}

func (s *service) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
