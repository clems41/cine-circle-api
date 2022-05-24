package mailServiceMock

import (
	mailService2 "cine-circle-api/external/mailService"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/envUtils"
)

var _ mailService2.Service = (*service)(nil)

type service struct {
	redirectAll     bool
	redirectAddress string
}

func New() (svc mailService2.Service) {
	redirectAddress := envUtils.GetFromEnvOrDefault(mailService2.EnvRedirectAddress, mailService2.DefaultRedirectAddress)
	svc = &service{
		redirectAll:     redirectAddress != mailService2.DefaultRedirectAddress,
		redirectAddress: redirectAddress,
	}
	return
}

func (svc *service) Send(form mailService2.SendEmailForm) (err error) {
	if svc.redirectAll {
		form.To = []string{svc.redirectAddress}
		form.Cc = nil
		form.Bcc = nil
	}
	if form.From == "" {
		form.From = mailService2.DefaultSender
	}
	logger.Infof("Email should be sent using mailService : form=%v", form)
	return
}
