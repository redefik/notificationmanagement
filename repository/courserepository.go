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
)

/*This package provides a set of functionality used to interact with the persistence layer
that store information about the courses*/

var dynamodbClient *dynamodb.DynamoDB

// Ad hoc errors returned by the functions
var UnknownError = errors.New("an unknown error occurred during the interaction with course data store")
var ConflictError = errors.New("a course with the provided information already exists")

// Encapsulates the fields of the DynamoDB item representing a course
type CourseItem struct {
	CourseName string
}

// initializeClient instantiate a DynamoDB client that will be then shared between that functions
// that interact with the data-store
func InitializeDynamoDbClient() {
	sessionInitializer := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbClient = dynamodb.New(sessionInitializer)
}

// Add a course to the data store returning a not-nil value in case of error
func CreateCourse(course entity.Course) error {
	// Convert the course in the format read by dynamodb
	courseItem := CourseItem{CourseName:course.Name + "_" + course.Department + "_" + course.Year} // Name, Department and Year acts as a composite key
	marshaledCourse, err := dynamodbattribute.MarshalMap(courseItem)
	if err != nil {
		return UnknownError
	}
	// Build the request for DynamoDB
	putItemInput := &dynamodb.PutItemInput{
		Item: marshaledCourse,
		ConditionExpression: aws.String("attribute_not_exists(CourseName)"),
		TableName: aws.String(config.Configuration.CoursesTableName),
	}
	_, err = dynamodbClient.PutItem(putItemInput)
	if err != nil {
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
