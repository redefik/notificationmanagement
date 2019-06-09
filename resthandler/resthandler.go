package resthandler

import (
	"encoding/json"
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"github.com/redefik/notificationmanagement/coursehandler"
	"github.com/redefik/notificationmanagement/entity"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func isValidBody(body entity.Course) bool {
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

	if !isValidBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Try to create the course
	err = coursehandler.CreateCourse(requestBody)
	if err != nil {
		// On error an appropriated status code is returned
		if err == coursehandler.ConflictError {
			MakeErrorResponse(w, http.StatusConflict, "Conflict - The course already exists")
			log.Println(err)
			return
		} else if err == coursehandler.UnknownError {
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

	if !isValidBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Try to create the course
	err = coursehandler.DeleteCourse(requestBody)
	if err != nil {
		// On error an appropriated status code is returned
		if err == coursehandler.NotFoundError {
			MakeErrorResponse(w, http.StatusNotFound, "Course Not Found")
			log.Println(err)
			return
		} else if err == coursehandler.UnknownError {
			MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			log.Println(err)
			return
		}
	}
	// On success 200 OK is returned
	w.WriteHeader(http.StatusOK)
}

func isValidMail(mail string) bool {
	// Format Validation
	err := checkmail.ValidateFormat(mail)
	if err != nil {
		return false
	}
	// Host Validation
	//err = checkmail.ValidateHost(mail)
	//if err != nil {
	//	return false
	//}
	return true
}

// AddCourseSubscription subscribes a student to the mailing list of the course
func AddCourseSubscription(w http.ResponseWriter, r *http.Request) {

	var requestBody entity.Course

	//Parse the body of the request
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Check if all the required fields have been provided in the body
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	requestBody.Department = strings.TrimSpace(requestBody.Department)
	requestBody.Year = strings.TrimSpace(requestBody.Year)

	if !isValidBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	urlParameters := mux.Vars(r)
	studentMail := urlParameters["studentMail"]

	// Check if the provided mail is valid
	if !isValidMail(studentMail) {
		MakeErrorResponse(w, http.StatusBadRequest, "Invalid Mail")
		log.Println("Invalid Mail")
		return
	}

	// Try to add the subscription
	err = coursehandler.AddStudent(requestBody, studentMail)
	if err != nil {
		// On error an appropriated status code is returned
		if err == coursehandler.NotFoundError {
			MakeErrorResponse(w, http.StatusNotFound, "Course Not Found")
			log.Println(err)
			return
		} else if err == coursehandler.UnknownError {
			MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			log.Println(err)
			return
		}
	}
	// On success 200 OK is returned
	w.WriteHeader(http.StatusOK)

}

// RemoveCourseSubscription removes a student from the mailing list of the provided course
func RemoveCourseSubscription(w http.ResponseWriter, r *http.Request) {
	var requestBody entity.Course

	//Parse the body of the request
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	// Check if all the required fields have been provided in the body
	requestBody.Name = strings.TrimSpace(requestBody.Name)
	requestBody.Department = strings.TrimSpace(requestBody.Department)
	requestBody.Year = strings.TrimSpace(requestBody.Year)

	if !isValidBody(requestBody) {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad request")
		log.Println("Bad Request")
		return
	}

	urlParameters := mux.Vars(r)
	studentMail := urlParameters["studentMail"]

	// Try to remove the subscription
	err = coursehandler.RemoveStudent(requestBody, studentMail)
	if err != nil {
		// On error an appropriated status code is returned
		if err == coursehandler.NotFoundError {
			MakeErrorResponse(w, http.StatusNotFound, "Course Not Found")
			log.Println(err)
			return
		} else if err == coursehandler.UnknownError {
			MakeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			log.Println(err)
			return
		}
	}
	// On success 200 OK is returned
	w.WriteHeader(http.StatusOK)
}
