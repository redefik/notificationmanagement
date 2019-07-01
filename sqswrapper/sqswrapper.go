package sqswrapper

/*This package contains a set of functions that wrap the API of Amazon Simple Queue Service SDK for Go*/

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/redefik/notificationmanagement/config"
)

// GetSqsClient builds a *sqs.SQS object that can be used to make requests to AWS SQS
func GetSqsClient() *sqs.SQS {
	sessionInitializer := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Configuration.AwsSqsRegion),
	}))
	//sessionInitializer := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))
	sqsClient := sqs.New(sessionInitializer)
	return sqsClient
}

// GetMessageQueueUrl returns the url of the SQS queue with the provided name
func GetMessageQueueUrl(sqsClient *sqs.SQS, queueName string) (string, error) {
	getQueueUrlOutput, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", errors.New("cannot retrieve message queue url:" + err.Error())
	}
	queueUrl := *getQueueUrlOutput.QueueUrl
	return queueUrl, nil
}

// ReadMessageFromQueue try to read a message from the queue with the given url
// using the provided client. The result, in JSON format, is stored in the message
// structure provided by the caller.
// waitTime is the duration (in seconds) for which the call waits for a message to arrive in the queue before returning
// The function returns the number of read messages (0 or 1), the receipt handler of the message and an error
func ReadJsonMessageFromQueue(sqsClient *sqs.SQS, queueUrl string, message interface{}, waitTime int64) (int, *string, error) {
	receiveMessageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(waitTime),
	}
	// polling
	receiveMessageOutput, err := sqsClient.ReceiveMessage(receiveMessageInput)
	if err != nil {
		return 0, nil, errors.New("error in retrieving message from queue:" + err.Error())
	}
	receivedMessages := receiveMessageOutput.Messages
	if len(receivedMessages) == 0 {
		return 0, nil, nil
	}
	// parsing
	messageBody := *receivedMessages[0].Body
	err = json.Unmarshal([]byte(messageBody), message)
	if err != nil {
		return 1, nil, errors.New("error in parsing the received message:" + err.Error())
	}
	return 1, receivedMessages[0].ReceiptHandle, nil
}

// DeleteMessageFromQueue delete the message with the provided handler from the queue with the given url.
func DeleteMessageFromQueue(sqsClient *sqs.SQS, queueUrl string, messageHandler *string) error {
	deleteMessageInput := &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: messageHandler,
	}
	_, err := sqsClient.DeleteMessage(deleteMessageInput)
	if err != nil {
		return errors.New("error in deleting the received message from the queue:" + err.Error())
	}
	return nil
}
