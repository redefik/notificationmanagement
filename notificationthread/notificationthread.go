package notificationthread

import (
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/coursehandler"
	"github.com/redefik/notificationmanagement/entity"
	"github.com/redefik/notificationmanagement/notificationhandler"
	"github.com/redefik/notificationmanagement/sqswrapper"
	"log"
	"math/rand"
	"time"
)

// polls an Amazon SQS message queue searching for e-mail to send to the subscribers of a course
func Run() {
	log.Println("Notification Thread launched...")
	sqsClient := sqswrapper.GetSqsClient()
	queueUrl, err := sqswrapper.GetMessageQueueUrl(sqsClient, config.Configuration.MessageQueueName)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Queue url", queueUrl)
	// Queue polling starts
	// When there aren't messages in the queue, the thread sleep for n*PollingWaitTime seconds before retrying, where
	// n is a random number in the range [0,i] and i is the number of times the polling resulted in no message read.
	// The maximum value of i is 6. (The algorithm is similar to the binary exponential backoff used in Ethernet to avoid
	// collisions).
	i := 0 // number of empty responses
	for {
		var message entity.Notification
		log.Println("Polling...")
		readMessages, err := sqswrapper.ReadJsonMessageFromQueue(sqsClient, queueUrl, &message, config.Configuration.PollingWaitTime)
		if err != nil {
			log.Println(err)
			continue
		}
		if readMessages == 0 {
			if i == 6 {
				i = 1
			} else {
				i++
			}
			n := rand.Intn(i)
			sleepingTime := time.Duration(n * int(config.Configuration.PollingWaitTime))
			log.Println("Sleep before retrying...")
			time.Sleep(sleepingTime * time.Second)
			continue
		}
		log.Println("Sending notification requests...")
		// send a mail containing the notification to the mailing list of the course
		course := entity.Course{Name: message.Name, Department: message.Department, Year: message.Year}
		mailingList, err := coursehandler.GetCourseMailingList(course)
		if err != nil {
			log.Println("error in getting course mailing list", err)
			continue
		}
		err = notificationhandler.SendNotificationToMultipleRecipients(message, config.Configuration.MailAddress, mailingList)
		if err != nil {
			log.Println("error in sending notification", err)
			continue
		}
	}
}
