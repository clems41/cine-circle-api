package mailServiceMock

import (
	"cine-circle-api/internal/service/mailService"
	"cine-circle-api/pkg/logger"
)

var _ Mock = (*mock)(nil)

type mock struct{}

type Mock interface {
	Send(form mailService.SendEmailForm) (err error)
}

func New() *mock {
	return &mock{}
}

func (mock *mock) Send(form mailService.SendEmailForm) (err error) {
	logger.Infof("ServiceMailMock send email with form : %v", form)
	return
}
