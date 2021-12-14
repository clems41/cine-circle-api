package mailServiceMock

import (
	"cine-circle-api/internal/service/mailService"
	"cine-circle-api/pkg/logger"
	"cine-circle-api/pkg/utils/envUtils"
)

var _ mailService.Service = (*service)(nil)

type service struct {
	redirectAll     bool
	redirectAddress string
}

func New() (svc mailService.Service) {
	redirectAddress := envUtils.GetFromEnvOrDefault(mailService.EnvRedirectAddress, mailService.DefaultRedirectAddress)
	svc = &service{
		redirectAll:     redirectAddress != mailService.DefaultRedirectAddress,
		redirectAddress: redirectAddress,
	}
	return
}

func (svc *service) Send(form mailService.SendEmailForm) (err error) {
	if svc.redirectAll {
		form.To = []string{svc.redirectAddress}
		form.Cc = nil
		form.Bcc = nil
	}
	if form.From == "" {
		form.From = mailService.DefaultSender
	}
	logger.Infof("Email should be sent using mailService : form=%v", form)
	return
}
