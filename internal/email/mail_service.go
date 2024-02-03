package email

type EmailService struct {
	client AwsClient
}

type EmailConfig struct {
	Endpoint  string
	Port      string
	Protocol  string
	Host      string
	Region    string
	AccessKey string
	SecretKey string
	From      string
}

type Mail struct {
	To      string `json:"to" example:"tomasfanciotti@gmail.com"`
	Subject string `json:"subject" example:"testing subject"`
	Body    string `json:"body" example:"body of the mail"`
}

func NewEmailService(emailConfig EmailConfig) (EmailService, error) {

	session := NewAwsSession(&emailConfig)
	err := session.Connect()
	if err != nil {
		return EmailService{}, err
	}

	return EmailService{client: session}, nil
}
