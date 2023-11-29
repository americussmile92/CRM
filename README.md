# Customer Management System

This project implements a simple Customer Management System with basic CRUD operations using a RESTful API. It is built in Go and utilizes the Gorilla Mux router for handling HTTP requests.

## Features

1. **Get All Customers**: Retrieve a list of all customers.
2. **Get Customer by ID**: Get data for a single customer by providing their ID.
3. **Add Customer**: Add a new customer to the system.
4. **Update Customer Information**: Update an existing customer's information.
5. **Remove Customer**: Remove a customer from the system.

## Getting Started

### Prerequisites

- Go installed on your machine.
- Gorilla Mux library: `go get -u github.com/gorilla/mux`

### Installation
How to run the project:

   ```bash
   git clone https://github.com/americussmile92/CRM
   cd CRM
   go get -u github.com/gorilla/mux
   go run main.go
