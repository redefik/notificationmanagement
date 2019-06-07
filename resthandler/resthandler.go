package resthandler

import (
	"encoding/json"
	"github.com/redefik/notificationmanagement/entity"
	"github.com/redefik/notificationmanagement/repository"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func isValidCourseCreationBody(body entity.Course) bool {
	alphanumericPattern := regexp.MustCompile("^[a-zA-Z0-9]*$")
	if !alphanumericPattern.MatchString(body.Name) {
		return false
	}
	alphaPattern := regexp.MustCompile("^[a-zA-Z]*$")
	if !alphaPattern.MatchString(body.Department) {
		return false
	}
	yearParts := strings.Split(body.Year, "-")
	if len(yearParts) != 2 {
		return false
	}
	startYear, err := strconv.Atoi(yearParts[0])
	if err != nil {
		return false
	}
	endYear, err := strconv.Atoi(yearParts[1])
	if err != nil {
		return false
	}
	if startYear >= endYear {
		return false
	}
	if body.Name == "" || body.Department == "" || body.Year == "" {
		return false
	}

	return true
}

// CreateCourse creates a course in the data-store using the information provided by the body of the HTTP request
func NewCourse(w http.ResponseWriter, r *http.Request) {
	var requestBody entity.Course

	//Parse the body of the request
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Check if all the required fields have been provided in the request
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	requestBody.Department = strings.TrimSpace(requestBody.Department)
	requestBody.Year = strings.TrimSpace(requestBody.Year)

	if !isValidCourseCreationBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Try to create the course
	err = repository.CreateCourse(requestBody)
	if err != nil {
		// On error an appropriated status code is returned
		if err == repository.ConflictError {
			MakeErrorResponse(w, http.StatusConflict, "Conflict - The course already exists")
			log.Println(err)
			return
		} else if err == repository.UnknownError {
			MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			log.Println(err)
			return
		}
	}
	// On success 201 is returned and the created course is provided inside the response
	responseBody, err := json.Marshal(requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)
}

// DeleteCourse delete the provided course from the data store.
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	var requestBody entity.Course

	//Parse the body of the request
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Check if all the required fields have been provided in the request
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	requestBody.Department = strings.TrimSpace(requestBody.Department)
	requestBody.Year = strings.TrimSpace(requestBody.Year)

	if !isValidCourseCreationBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Try to create the course
	err = repository.DeleteCourse(requestBody)
	if err != nil {
		// On error an appropriated status code is returned
		if err == repository.NotFoundError {
			MakeErrorResponse(w, http.StatusNotFound, "Course Not Found")
			log.Println(err)
			return
		} else if err == repository.UnknownError {
			MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			log.Println(err)
			return
		}
	}
	// On success 200 OK is returned
	w.WriteHeader(http.StatusOK)
}
