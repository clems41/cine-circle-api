package circleDom

var _ Service = (*service)(nil)

type Service interface {
	Create(creation Creation) (result Result, err error)
	Update(update Update) (result Result, err error)
	Delete(delete Delete) (err error)
	Get(get Get) (result Result, err error)
}

type service struct {
	r Repository
}

type Repository interface {
	Create(creation Creation) (result Result, err error)
	Update(update Update) (result Result, err error)
	Delete(delete Delete) (err error)
	Get(get Get) (result Result, err error)
}

func NewService(r Repository) Service {
	return &service{
		r:                              r,
	}
}

func (svc *service) Create(creation Creation) (result Result, err error) {
	err = creation.Valid()
	if err != nil {
		return
	}
	return svc.r.Create(creation)
}

func (svc *service) Update(update Update) (result Result, err error) {
	err = update.Valid()
	if err != nil {
		return
	}
	return svc.r.Update(update)
}

func (svc *service) Delete(delete Delete) (err error) {
	err = delete.Valid()
	if err != nil {
		return
	}
	return svc.r.Delete(delete)
}

func (svc *service) Get(get Get) (result Result, err error) {
	err = get.Valid()
	if err != nil {
		return
	}
	return svc.r.Get(get)
}
