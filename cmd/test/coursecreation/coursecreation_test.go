package coursecreation

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

// createTestMicroserviceCourseCreation builds an http handler used to test functionality of course creation
func createTestMicroserviceCourseCreation() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/notification_management/api/v1.0/course", resthandler.NewCourse).Methods(http.MethodPost)
	return r
}

func setup() {
	config.SetConfiguration("../../../config/config-test.json")
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	courseItem := repository.CourseItem{CourseName: "testexistentcourse_testdepartment_2018-2019"}
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
				S: aws.String("testcourse_testdepartment_2018-2019"),
			},
		},
		TableName: aws.String(config.Configuration.CoursesTableName),
	}

	_, err := client.DeleteItem(input)
	if err != nil {
		log.Println(err)
	}
	input = &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"CourseName": {
				S: aws.String("testexistentcourse_testdepartment_2018-2019"),
			},
		},
		TableName: aws.String(config.Configuration.CoursesTableName),
	}

	_, err = client.DeleteItem(input)
	if err != nil {
		log.Println(err)
	}
}

// TestMain perform setup and tear-down needed by the test
func TestMain(m *testing.M) {
	// The data-store is populated with the item needed for testing purpose
	setup()
	// Run test
	code := m.Run()
	// The data-store is cleaned up after the tests
	tearDown()
	os.Exit(code)
}

// TestCourseCreationSuccess tests the following scenario: the course creation asked by the client
// is correctly done.
func TestCourseCreationSuccess(t *testing.T) {

	repository.InitializeDynamoDbClient()

	// It is assumed that a course with the following information does not exist in the testing data-store
	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcourse")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Error("Expected 201 Created but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))

	}
}

// TestCourseCreationMissingField tests the following scenario: the course creation request lacks a field.
// Therefore, the microservice response should be 400 Bad Request
func TestCourseCreationMissingFields(t *testing.T) {
	repository.InitializeDynamoDbClient()

	jsonBody := simplejson.New()
	// the Name field lacks
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad Request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestCourseCreationConflictCourse tests the following scenario: the client asks for creating a course that already exists
// The response should be 409 Conflict
func TestCourseCreationConflictCourse(t *testing.T) {
	repository.InitializeDynamoDbClient()

	jsonBody := simplejson.New()
	jsonBody.Set("name", "testexistentcourse")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Error("Expected 409 Conflict but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestCourseCreationInvalidField tests the attempt to create a course providing an invalid field
func TestCourseCreationInvalidField(t *testing.T) {
	repository.InitializeDynamoDbClient()

	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcourse")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2020-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseCreation()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad Request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
