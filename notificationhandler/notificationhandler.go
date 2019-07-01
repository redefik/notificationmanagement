package notificationhandler

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/entity"
)

/*This package contains a set of functions that wrap the API of Amazon Simple Email Service SDK for Go*/

var Client *ses.SES

// InitializeSesClient instantiate a Sess client that will be used to make API requests to SES. The initialization
// is performed once because, as reported in the documentation, the client is safe to be used concurrently
func InitializeSesClient(region string) error {
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	Client = ses.New(newSession)
	return nil
}

// SendNotificationToMultipleRecipients sends the message provided to the given recipients.
// The name  of the mail template is read from config.Configuration and the message is built replacing the template
// parameters with the content of the provided notificationthread struct
func SendNotificationToMultipleRecipients(message entity.Notification, from string, to []string) error {
	var destinations []*ses.BulkEmailDestination
	// The parameters of the mail template are set according to the message provided
	defaultTemplateData := "{ \"courseName\":\"%s\", \"year\":\"%s\", \"body\":\"%s\"}"
	defaultTemplateData = fmt.Sprintf(defaultTemplateData, message.Name, message.Year, message.Message)
	for i := 0; i < len(to); i++ {
		bulkEmailDestination := &ses.BulkEmailDestination{Destination: &ses.Destination{BccAddresses: []*string{aws.String(to[i])}},
			ReplacementTemplateData: aws.String(defaultTemplateData)}
		destinations = append(destinations, bulkEmailDestination)
	}
	sendBulkTemplatedEmailInput := &ses.SendBulkTemplatedEmailInput{
		Source:              aws.String(from),
		Destinations:        destinations,
		Template:            aws.String(config.Configuration.MailTemplate),
		DefaultTemplateData: aws.String(defaultTemplateData),
	}
	// TODO sono ammessi al più 50 destinatari per volta
	_, err := Client.SendBulkTemplatedEmail(sendBulkTemplatedEmailInput)
	if err != nil {
		return errors.New("error in sending e-mails")
	}
	return nil
}
