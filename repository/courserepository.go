package repository

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/entity"
	"log"
)

/*This package provides a set of functionality used to interact with the persistence layer
that store information about the courses*/

var dynamodbClient *dynamodb.DynamoDB

// Ad hoc errors returned by the functions
var UnknownError = errors.New("an unknown error occurred during the interaction with course data store")
var ConflictError = errors.New("a course with the provided information already exists")
var NotFoundError = errors.New("the provided course does not exist in the data store")

// Encapsulates the fields of the DynamoDB item representing a course
type CourseItem struct {
	CourseName  string
	MailingList []string
}

// As above but without MailingList field (to avoid problems about the type when the mailing list gets empty)
type CourseCreationItem struct {
	CourseName string
}

// initializeClient instantiate a DynamoDB client that will be then shared between the functions
// that interact with the data-store
func InitializeDynamoDbClient() {
	sessionInitializer := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbClient = dynamodb.New(sessionInitializer)
}

// Add a course to the data store returning a not-nil value in case of error:
// - ConflictError it the caller try to create an existent course
// - UnknownError otherwise
func CreateCourse(course entity.Course) error {
	// Convert the course in the format read by dynamodb
	courseItem := CourseCreationItem{CourseName: course.Name + "_" + course.Department + "_" + course.Year} // Name, Department and Year acts as a composite key
	marshaledCourse, err := dynamodbattribute.MarshalMap(courseItem)
	if err != nil {
		log.Println(err)
		return UnknownError
	}
	// Build the request for DynamoDB
	putItemInput := &dynamodb.PutItemInput{
		Item:                marshaledCourse,
		ConditionExpression: aws.String("attribute_not_exists(CourseName)"),
		TableName:           aws.String(config.Configuration.CoursesTableName),
	}
	_, err = dynamodbClient.PutItem(putItemInput)
	if err != nil {
		log.Println(err)
		// check if AWS DynamoDB raised an error
		awsError, ok := err.(awserr.Error)
		if ok {
			switch awsError.Code() {
			// raised when the client try to create a course that already exists
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return ConflictError
			default:
				return UnknownError
			}
		}
		return UnknownError
	}
	return nil
}

// Delete the given course from the data store returning a not-nil value in case of error:
// - NotFoundError when the caller try to delete a not existent course
// - UnknownError otherwise
func DeleteCourse(course entity.Course) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {
				S: aws.String(course.Name + "_" + course.Department + "_" + course.Year),
			},
		},
		ConditionExpression: aws.String("attribute_exists(CourseName)"),
		TableName:           aws.String(config.Configuration.CoursesTableName),
	}
	_, err := dynamodbClient.DeleteItem(deleteItemInput)
	if err != nil {
		log.Println(err)
		// check if AWS DynamoDB raised an error
		awsError, ok := err.(awserr.Error)
		if ok {
			switch awsError.Code() {
			// raised when the given course does not exist in the data store
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return NotFoundError
			default:
				return UnknownError
			}
		}
		return UnknownError
	}
	return nil
}

// Add the provided mail to the list of mail address associated to the given course. So the student will receive news
// about the course. The function returns a not-nil value in case of error:
// - NotFoundError when the caller try to update a not existent course
// - UnknownError otherwise
func AddStudent(course entity.Course, studentMail string) error {
	newMail := &dynamodb.AttributeValue{
		S: aws.String(studentMail),
	}
	var mailList []*dynamodb.AttributeValue
	// the mail address will be appended to the mailing list of the course.
	// it must be embedded inside a slice
	mailList = append(mailList, newMail)

	updateItemInput := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {S: aws.String(course.Name + "_" + course.Department + "_" + course.Year)},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":mail": {
				L: mailList,
			},
			":empty_list": {
				// this must be provided for the first mail, when the list does not exist yet
				L: []*dynamodb.AttributeValue{},
			},
		},
		ConditionExpression: aws.String("attribute_exists(CourseName)"),
		UpdateExpression:    aws.String("SET MailingList = list_append(if_not_exists(MailingList, :empty_list), :mail)"),
		TableName:           aws.String(config.Configuration.CoursesTableName),
	}

	_, err := dynamodbClient.UpdateItem(updateItemInput)
	if err != nil {
		log.Println(err)
		log.Println(err)
		// check if AWS DynamoDB raised an error
		awsError, ok := err.(awserr.Error)
		if ok {
			switch awsError.Code() {
			// raised when the given course does not exist in the data store
			case dynamodb.ErrCodeConditionalCheckFailedException:
				return NotFoundError
			default:
				return UnknownError
			}
		}
		return UnknownError
	}

	return nil
}

func removeMailFromList(mail string, list []string) []string {
	var newSlice []string
	for i := 0; i < len(list); i++ {
		if mail != list[i] {
			newSlice = append(newSlice, list[i])
		}
	}
	return newSlice
}

// Remove the provided mail from the list of mail address associated to the given course.
// The function returns a not-nil value in case of error:
// - NotFoundError when the caller try to update a not existent course
// - UnknownError otherwise
func RemoveStudent(course entity.Course, studentMail string) error {
	// Search for the provided course
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(config.Configuration.CoursesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {
				S: aws.String(course.Name + "_" + course.Department + "_" + course.Year),
			},
		},
	}
	getResult, err := dynamodbClient.GetItem(getItemInput)
	if err != nil {
		log.Println(err)
		return UnknownError
	}
	var matchingCourse CourseItem
	err = dynamodbattribute.UnmarshalMap(getResult.Item, &matchingCourse)
	if err != nil {
		log.Println(err)
		return UnknownError
	}
	if matchingCourse.CourseName == "" {
		return NotFoundError
	}
	// If the course exist, its mailing list is updated removing the given address
	matchingCourse.MailingList = removeMailFromList(studentMail, matchingCourse.MailingList)
	var marshaledCourse map[string]*dynamodb.AttributeValue
	if len(matchingCourse.MailingList) == 0 {
		// If the mailing list is empty, the corresponding attribute is removed to avoid problems concerning the type
		updatedCourse := CourseCreationItem{CourseName: matchingCourse.CourseName}
		marshaledCourse, err = dynamodbattribute.MarshalMap(updatedCourse)
	} else {
		marshaledCourse, err = dynamodbattribute.MarshalMap(matchingCourse)
	}
	if err != nil {
		log.Println(err)
		return UnknownError
	}
	putItemInput := &dynamodb.PutItemInput{
		Item:      marshaledCourse,
		TableName: aws.String(config.Configuration.CoursesTableName),
	}
	_, err = dynamodbClient.PutItem(putItemInput)
	if err != nil {
		log.Println(err)
		return UnknownError
	}
	return nil
}
