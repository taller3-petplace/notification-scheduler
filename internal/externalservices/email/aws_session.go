package email

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type AwsClient struct {
	session *session.Session
	client  *ses.SES
	config  *EmailConfig
}

func NewAwsSession(mailConfig *EmailConfig) AwsClient {
	return AwsClient{
		config: mailConfig,
	}
}

func (c *AwsClient) Connect() error {

	newSession, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(c.config.Region),
			Credentials: credentials.NewStaticCredentials(c.config.AccessKey, c.config.SecretKey, ""),
		})

	if err != nil {
		return err
	}

	c.session = newSession

	// Crea un cliente SES
	svc := ses.New(c.session)
	c.client = svc

	return nil
}

func (c *AwsClient) SendEmail(mail Mail) error {

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(mail.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(mail.Body),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(mail.Subject),
			},
		},
		Source: aws.String(c.config.From),
	}

	// Envía el correo electrónico
	output, err := c.client.SendEmail(input)

	if err != nil {
		fmt.Println("Error al enviar el correo electrónico:", err)
		return err
	}

	fmt.Printf("Correo electrónico enviado! Output: %v\n", output.MessageId)
	return nil
}
