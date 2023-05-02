package main

import (	
"errors"	
"fmt"	
"log"	
"net/http")

type person struct {
	name     string
	location string
	age      string
}

var People = make(map[string]person)

func main() {
	http.HandleFunc("/people", peopleHandler)
	http.ListenAndServe(":5000", nil)
}

func peopleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		addPerson(w, r)
	case "GET":
		getPerson(w, r)
	case "DELETE":
		deletePerson(w, r)
	case "PATCH":
		updatePerson(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Sorry, only GET/DELETE/PATCH/POST methods are supported.")
	}
}

func nameQuery(w http.ResponseWriter, r *http.Request) (string, error) {
	name := r.URL.Query().Get("name")
	if len(name) == 0 {
		return "", errors.New("missing name parameter")
	}
	return name, nil
}

func locationQuery(w http.ResponseWriter, r *http.Request) (string, error) {
	location := r.URL.Query().Get("location")
	if len(location) == 0 {
		return "", errors.New("missing location parameter")
	}
	return location, nil
}

func ageQuery(w http.ResponseWriter, r *http.Request) (string, error) {
	age := r.URL.Query().Get("age")
	if len(age) == 0 {
		return "", errors.New("missing age parameter")
	}
	return age, nil
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	inputtedName, err := nameQuery(w, r)
	inputtedLocation, err1 := locationQuery(w, r)
	inputtedAge, err2 := ageQuery(w, r)
	if err != nil && err1 != nil && err2 != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	_, exists := People[inputtedName]
	if !exists {
		newPerson := person{name: inputtedName, location: inputtedLocation, age: inputtedAge}
		People[newPerson.name] = newPerson
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("error - person exists", inputtedName)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "error - person exists")
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	name, err := nameQuery(w, r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	personGet, exists := People[name]
	if exists {
		log.Println("person exists", name)
		message := fmt.Sprintf("[%s found, %s found, %s found]", personGet.name, personGet.location, personGet.age)
		fmt.Fprint(w, message)
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Printf("%v", People)
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s not found", name)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	name, err := nameQuery(w, r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	delete(People, name)
	log.Printf("%v", People)
	fmt.Fprint(w, "deleted successfully")
	w.WriteHeader(http.StatusOK)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	name, err := nameQuery(w, r)
	location, err1 := locationQuery(w, r)
	age, err2 := ageQuery(w, r)
	if err != nil && err1 != nil && err2 != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	personUpdate, exists := People[name]
	if exists {
		personUpdate.age = age
		personUpdate.location = location
		People[name] = personUpdate
		log.Println("update successfully", name, location, age)
		fmt.Fprint(w, "update successfully")
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println("failed to update")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s not found", name)
}