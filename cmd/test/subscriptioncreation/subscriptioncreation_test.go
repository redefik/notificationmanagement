package subscriptioncreation

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/notificationmanagement/config"
	"github.com/redefik/notificationmanagement/coursehandler"
	"github.com/redefik/notificationmanagement/resthandler"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// createTestMicroserviceCourseSubscriptionCreation builds an http handler used to test functionality about student subscriptions
// to course's mailing list
func createTestMicroserviceCourseSusbscriptionCreation() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/notification_management/api/v1.0/course/student/{studentMail}", resthandler.AddCourseSubscription).Methods(http.MethodPut)
	return r
}

func setup() {
	config.SetConfiguration("../../../config/config-test.json")
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	courseItem := coursehandler.CourseCreationItem{CourseName: "testcoursesubscription_testdepartment_2018-2019"}
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

// TestCourseCourseSubscriptionSuccess tests the following scenario: the client correctly add a subscription to the mailing list
// of a course.
func TestCourseSubscriptionCreationSuccess(t *testing.T) {

	coursehandler.InitializeDynamoDbClient()

	// It is assumed that a course with the following information exists in the testing data-store
	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcoursesubscription")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPut, "/notification_management/api/v1.0/course/student/isssr.ticketing@gmail.com",
		bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseSusbscriptionCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 20O Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))

	}
}

// TestCourseSubscriptionInvalidMail tests the following scenario: the student cannot be added to the mailing list
// because the mail addressi is not valid
func TestCourseSubscriptionInvalidMail(t *testing.T) {
	coursehandler.InitializeDynamoDbClient()

	// It is assumed that a course with the following information exists in the testing data-store
	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcoursesubscription")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	// the provided mail address does not exist
	request, _ := http.NewRequest(http.MethodPut, "/notification_management/api/v1.0/course/student/invalidMail",
		bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseSusbscriptionCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad Request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))

	}
}
