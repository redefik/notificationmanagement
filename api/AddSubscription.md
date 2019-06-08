**Add Subscription**
----
  Adds a student to the mailing list of a course.

* **URL**

  /course/student/:studentMail

* **Method:**

  `PUT`
  
*  **URL Params**

   **Required:**
 
   `studentMail=[string]`<br/>
   

* **Data Params**

    `{name:"Advanced Calculus",
      department: "Science",
      year: "2018-2019"
    }`

* **Success Response:**

  * **Code:** 200 OK <br />
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error : "Course Not Found"}`

  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Bad request" }`
    
  OR

  * **Code:** 400 BAD REQUEST <br />
    **Content:** `{ error : "Invalid Mail" }`
    This is returned when the mail is not valid (the format is not valid or
    the host does not exist)
    
  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
    