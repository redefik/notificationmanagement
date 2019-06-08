package subscriptiondeletion

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/repository"
	"github.com/redefik/notificationmanagement/resthandler"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// createTestMicroserviceCourseSubscriptionDeletion builds an http handler used to test functionality about student subscriptions
// to course's mailing list
func createTestMicroserviceCourseSusbscriptionDeletion() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/notification_management/api/v1.0/course/student/{studentMail}", resthandler.RemoveCourseSubscription).Methods(http.MethodDelete)
	return r
}

func setup() {
	config.SetConfiguration("../../../config/config-test.json")
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	mailingList := []string{"isssr.ticketing@gmail.com", "other@mail.com"}
	courseItem := repository.CourseItem{CourseName: "testcoursesubscription_testdepartment_2018-2019", MailingList: mailingList}
	marshaledCourse, err := dynamodbattribute.MarshalMap(courseItem)
	if err != nil {
		log.Println(err)
	}
	client := dynamodb.New(newSession)
	input := &dynamodb.PutItemInput{
		Item:      marshaledCourse,
		TableName: aws.String(config.Configuration.CoursesTableName),
	}
	_, err = client.PutItem(input)
	if err != nil {
		log.Println(err)
	}
}

func tearDown() {
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := dynamodb.New(newSession)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {
				S: aws.String("testcoursesubscription_testdepartment_2018-2019"),
			},
		},
		TableName: aws.String(config.Configuration.CoursesTableName),
	}

	_, err := client.DeleteItem(input)
	if err != nil {
		log.Println(err)
	}
}

// TestMain perform courseDeletionTest setup and tear-down needed by the test
func TestMain(m *testing.M) {
	// The data-store is populated with the item needed for testing purpose
	setup()
	// Run test
	code := m.Run()
	// The data-store is cleaned up after the tests
	tearDown()
	os.Exit(code)
}

// TestCourseSubscriptionDeletionSuccess tests the following scenario: the client correctly remove a subscription from the mailing list
// of a course.
func TestCourseSubscriptionDeletionSuccess(t *testing.T) {

	repository.InitializeDynamoDbClient()

	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcoursesubscription")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodDelete, "/notification_management/api/v1.0/course/student/isssr.ticketing@gmail.com",
		bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseSusbscriptionDeletion()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 20O Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

	// check if the mailing list is up-to-date
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(config.Configuration.CoursesTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {
				S: aws.String("testcoursesubscription_testdepartment_2018-2019"),
			},
		},
	}
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	dynamodbClient := dynamodb.New(newSession)
	getResult, err := dynamodbClient.GetItem(getItemInput)
	if err != nil {
		t.Error("Error in retrieving the updated course")
	}
	var updatedCourse repository.CourseItem
	err = dynamodbattribute.UnmarshalMap(getResult.Item, &updatedCourse)
	if err != nil {
		t.Error("Error in retrieving the updated course")
	}
	for i := 0; i < len(updatedCourse.MailingList); i++ {
		if updatedCourse.MailingList[i] == "isssr.ticketing@gmail.com" {
			t.Error("Deletion not done")
		}
	}
}
