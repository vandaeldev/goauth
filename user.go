package main

type User struct {
	ID        uint
	Email     string
	Password  string
	FirstName *string
	LastName  *string
}
