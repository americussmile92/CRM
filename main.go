package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

// Getting a list of all customers
// Getting data for a single customer
// Adding a customer
// Updating a customer's information
// Removing a customer

type Customer struct {
	Id int
    Name string
    Role string
	Email string
	Phone int
	Contacted bool
}

type CustomerUpdate struct {
	Name string
    Role string
	Email string
	Phone int
	Contacted bool
}

var customerList = []Customer{
	{Id: 1, Name: "John Doe", Role: "Admin", Email: "john.doe@example.com", Phone: 1234567890, Contacted: true},
	{Id: 2, Name: "Jane Smith", Role: "User", Email: "jane.smith@example.com", Phone: 9876543210, Contacted: false},
	{Id: 3, Name: "Bob Johnson", Role: "User", Email: "bob.johnson@example.com", Phone: 5551234567, Contacted: true},
	{Id: 4, Name: "Alice Williams", Role: "Admin", Email: "alice.williams@example.com", Phone: 9998887777, Contacted: false},
	{Id: 5, Name: "Charlie Brown", Role: "User", Email: "charlie.brown@example.com", Phone: 1112223333, Contacted: true},
	{Id: 6, Name: "Eva Miller", Role: "User", Email: "eva.miller@example.com", Phone: 4445556666, Contacted: false},
	{Id: 7, Name: "David Davis", Role: "Admin", Email: "david.davis@example.com", Phone: 7778889999, Contacted: true},
	{Id: 8, Name: "Grace Taylor", Role: "User", Email: "grace.taylor@example.com", Phone: 2223334444, Contacted: false},
	{Id: 9, Name: "Frank Anderson", Role: "User", Email: "frank.anderson@example.com", Phone: 6667778888, Contacted: true},
	{Id: 10, Name: "Helen White", Role: "Admin", Email: "helen.white@example.com", Phone: 3334445555, Contacted: false},
}

func getCustomers(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")  
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customerList)
}

func getCustomerById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'id' parameter in the URL")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid 'id' parameter: %v", err)
		return
	}

	for _, c := range customerList {
		if c.Id == id {
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(c); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error encoding customer data: %v", err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Customer with ID %d not found", id)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'id' parameter in the URL")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid 'id' parameter: %v", err)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updatedCustomer CustomerUpdate
	if err := json.Unmarshal(reqBody, &updatedCustomer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var index = -1
	for i, existingCustomer := range customerList {
		if existingCustomer.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	customerList[index].Name = updatedCustomer.Name
	customerList[index].Role = updatedCustomer.Role
	customerList[index].Email = updatedCustomer.Email
	customerList[index].Contacted = updatedCustomer.Contacted

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customerList[index])
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newEntry Customer

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(reqBody, &newEntry); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, existingCustomer := range customerList {
		if existingCustomer.Id == newEntry.Id {
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	customerList = append(customerList, newEntry)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customerList)
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing 'id' parameter in the URL")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid 'id' parameter: %v", err)
		return
	}

	var index = -1
	for i, existingCustomer := range customerList {
		if existingCustomer.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	customerList = append(customerList[:index], customerList[index+1:]...)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(customerList)
}

func main() {
	// Instantiate a new router
	router := mux.NewRouter()
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customer/{id}", getCustomerById).Methods("GET")
	router.HandleFunc("/customer", createCustomer).Methods("POST")
	router.HandleFunc("/customer/{id}", updateCustomer).Methods("PATCH")
	router.HandleFunc("/customer/{id}", deleteCustomer).Methods("DELETE")
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/", fileServer)
	fmt.Println("Server is starting on port 3000...")
	// Pass the customer router into ListenAndServe
	http.ListenAndServe(":3000", router)
}