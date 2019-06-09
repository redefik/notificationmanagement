package entity

/* Here you can find the structures that encapsulates the fields of the application domain objects*/

// Encapsulates the fields of the course creation request
type Course struct {
	Name       string `json:"name"`
	Department string `json:"department"`
	Year       string `json:"year"`
}

// Encapsulates the fields of the notifications read by the sender thread from the SQS queue
type Notification struct {
	Name       string `json:"name"`
	Department string `json:"department"`
	Year       string `json:"year"`
	Message    string `json:"message"`
}
