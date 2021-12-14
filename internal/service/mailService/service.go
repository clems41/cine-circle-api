package mailService

import (
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/envUtils"
)

var _ Service = (*service)(nil)

type Service interface {
	Send(form SendEmailForm) (err error)
}

type service struct {
	redirectAll     bool
	redirectAddress string
}

func New() (svc Service, err error) {
	redirectAddress := envUtils.GetFromEnvOrDefault(envRedirectAddress, defaultRedirectAddress)
	svc = &service{
		redirectAll:     redirectAddress != defaultRedirectAddress,
		redirectAddress: redirectAddress,
	}
	return
}

func (svc *service) Send(form SendEmailForm) (err error) {
	if svc.redirectAll {
		form.To = []string{svc.redirectAddress}
		form.Cc = nil
		form.Bcc = nil
	}
	logger.Infof("Email should be sent using mailService : form=%v", form)
	return
}
