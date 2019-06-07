**Delete Course**
----
  Delete a course with the information provided in the JSON body of the request.

* **URL**

  /course/

* **Method:**

  `DELETE`
  
*  **URL Params**

   **Required:**
 
   None
   

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

  * **Code:** 500 INTERNAL SERVER ERROR <br />
    **Content:** `{ error : "Internal Server Error" }`
    