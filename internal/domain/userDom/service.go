package userDom

var _ Service = (*service)(nil)

type Service interface {
	CreateUser(creation Creation) (result Result, err error)
	GetUser(get Get) (result Result, err error)
}

type service struct {
	r Repository
}

type Repository interface {
	CreateUser(creation Creation) (result Result, err error)
	GetUser(get Get) (result Result, err error)
}

func NewService(r Repository) Service {
	return &service{
		r:                              r,
	}
}

func (svc *service) CreateUser(creation Creation) (result Result, err error) {
	err = creation.Valid()
	if err != nil {
		return
	}
	return svc.r.CreateUser(creation)
}

func (svc *service) GetUser(get Get) (result Result, err error) {
	return svc.r.GetUser(get)
}
