package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	People = map[string]person{
		"alice": {name: "alice", location: "london", age: "29"},
		"bob":   {name: "bob", location: "wales", age: "42"},
		"sue":   {name: "sue", location: "manchester", age: "51"},
	}
}

func TestPersonHandler(t *testing.T) {
	tt := []struct {
		name          string //Name = name of the test.
		method        string //Switch/case http method.
		nameParameter string //Query param for name (person's name - key).
		want          string //Expected response.
		statusCode    int    //Status response.
	}{
		{
			name:          "missing name parameter", //Negative test - checks when query param name is missing.
			method:        http.MethodGet,           //Trying to 'get' empty person.
			nameParameter: "",                       //Empty name.
			want:          "missing name parameter", //Expected error message (from my code).
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "get with non existing person", //Negative test - checks when the name doesn't exist in the kvs (person with that name does not exist).
			method:        http.MethodGet,                 //Trying to get someone that does not exist.
			nameParameter: "ffion",                        //Person name.
			want:          "ffion not found",              //Person not found as does not exist in KVS.
			statusCode:    http.StatusNotFound,
		},
		{
			name:          "unsupported method",                                       //Negative test when trying to do a request with unsupported http method (switch/cases method that does not exist).
			method:        http.MethodPut,                                             //Trying to do a method that does not exist in the KVS.
			nameParameter: "alys",                                                     //Person name.
			want:          "Sorry, only GET/DELETE/PATCH/POST methods are supported.", //Expected error message (from my code).
			statusCode:    http.StatusMethodNotAllowed,
		},
		{
			name:          "update with non existing person", // Negative test - trying to update someone who isn't in the store
			method:        http.MethodPatch,                  // Trying to update non existing person.
			nameParameter: "minnie",                          //Person name.
			want:          "minnie not found",                //Person not found as does not exist in KVS.
			statusCode:    http.StatusNotFound,
		},
		{
			name:          "adding someone that exists", //Negative test - trying to add someone who is already in store.
			method:        http.MethodPost,              //Trying to add someone who exists.
			nameParameter: "bob",                        //Bob exists in people variable in func init.
			want:          "error - person exists",      //Expected error message (from my code).
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "check delete", //Positive test for seeing if someone has successfully been deleted from the kvs.
			method:        http.MethodDelete,
			nameParameter: "alice",
			want:          "deleted successfully",
			statusCode:    http.StatusOK,
		},
		{
			name:          "check get", //Positive test for successfully getting 'person' from kvs by their name.
			method:        http.MethodGet,
			nameParameter: "sue",
			want:          "[sue found, manchester found, 51 found]",
			statusCode:    http.StatusOK,
		},
		{
			name:          "check update", //Positive test for successfully updating person age from kvs by their name.
			method:        http.MethodPatch,
			nameParameter: "bob",
			want:          "update successfully",
			statusCode:    http.StatusOK,
		},
		{
			name:          "check adding new person", //Positive test for successfully adding someone new into kvs.
			method:        http.MethodPost,
			nameParameter: "ian",
			want:          "",
			statusCode:    http.StatusOK,
		},
	}

	//This goes through each of our test cases from above.
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			log.Println(People)
			endpoint := fmt.Sprintf("/people?name=%s", tc.nameParameter) //Passing in the name parameter to the endpoint (endpoint from line 21 in main).
			request := httptest.NewRequest(tc.method, endpoint, nil)     //Creates new request with endpoint.
			responseRecorder := httptest.NewRecorder()                   //Records the test response, this is from HTTP library (unsure how works).

			peopleHandler(responseRecorder, request) //Calls my code from main.

			if responseRecorder.Code != tc.statusCode { //Checking if the response code matches expected status code.
				t.Errorf("test case: '%q', Want status '%d', got '%d'", tc.name, tc.statusCode, responseRecorder.Code) //Error print out response.
			}

			if responseRecorder.Body.String() != tc.want { //Checks the body of the response is the same as my want (from test cases).
				t.Errorf("test case: '%q', Want '%s', got '%s'", tc.name, tc.want, responseRecorder.Body) //Error print out response.
			}
		})
	}
}