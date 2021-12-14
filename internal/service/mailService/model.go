package mailService

type SendEmailForm struct {
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Cc          []string     `json:"cc,omitempty"`
	Bcc         []string     `json:"bcc,omitempty"`
	Subject     string       `json:"subject,omitempty"`
	Message     string       `json:"message,omitempty"`
	Password    string       `json:"password,omitempty"`
	Attachments []Attachment `json:"attachments, omitempty"`
	Embedded    []Attachment `json:"embedded, omitempty"` // inline
	Tags        []string     `json:"tags,omitempty"`      // tags mailGun only
	Html        bool         `json:"html"`                // to set proper content-type
}

type Attachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type mailApiResponse struct {
}
