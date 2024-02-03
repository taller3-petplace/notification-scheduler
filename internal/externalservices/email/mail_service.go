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
	To      string `json:"to" binding:"required" example:"tomasfanciotti@gmail.com"`
	Subject string `json:"subject" binding:"required" example:"testing subject"`
	Body    string `json:"body" binding:"required" example:"body of the mail"`
}
