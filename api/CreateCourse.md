**Create Course**
----
  Create a course with the information provided in the JSON body of the request.

* **URL**

  /course/

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   None
   

* **Data Params**

    `{name:"Advanced Calculus",
      department: "Science",
      year: "2018-2019"
    }`

* **Success Response:**

  * **Code:** 201 CREATED <br />
    **Content:** `{name: "Advanced Calculus",
                   department: "Science",
                   year: "2018-2019"
                     }`
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error : "Conflict - The course already exists"}`
    This is returned when a course with the given information already exists

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
    