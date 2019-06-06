package entity


/* Here you can find the structures that encapsulates the fields of the application domain objects*/

// Encapsulates the fields of the course creation request
type Course struct {
	Name string `json:"name"`
	Department string `json:"department"`
	Year string `json:"year"`
}
