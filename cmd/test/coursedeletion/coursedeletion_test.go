package coursedeletion

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

// createTestMicroserviceCourseDeletion builds an http handler used to test functionality of course deletion
func createTestMicroserviceCourseDeletion() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/notification_management/api/v1.0/course", resthandler.DeleteCourse).Methods(http.MethodDelete)
	return r
}

func setup() {
	config.SetConfiguration("../../../config/config-test.json")
	newSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	courseItem := repository.CourseItem{CourseName: "testcoursedeletion_testdepartment_2018-2019"}
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

// TestCourseDeletionSuccess tests the following scenario: the course deletion asked by the client
// is correctly done.
func TestCourseDeletionSuccess(t *testing.T) {

	repository.InitializeDynamoDbClient()

	// It is assumed that a course with the following information exists in the testing data-store
	jsonBody := simplejson.New()
	jsonBody.Set("name", "testcoursedeletion")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodDelete, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseDeletion()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 20O Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))

	}
}

// TestCourseDeletionNotFoundCourse tests the following scenario: the client tries to delete a not existent course.
// Therefore, the microservice response should be 404 Not Found
func TestCourseDeletionNotFoundCourse(t *testing.T) {
	repository.InitializeDynamoDbClient()

	jsonBody := simplejson.New()
	// it is assumed that a course with the given information does not exist in the data store
	jsonBody.Set("name", "testcoursenotexistent")
	jsonBody.Set("department", "testdepartment")
	jsonBody.Set("year", "2018-2019")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodDelete, "/notification_management/api/v1.0/course", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestMicroserviceCourseDeletion()
	// simulates a request-response interaction between client and microservice
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
