package mailService

type Service interface {
	Send(form SendEmailForm) (err error)
}
