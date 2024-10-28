package model


type Email struct {
	To []string
	Subject string
	Body string
}

type User struct {
	Id int64
	Username string
	Email string
	Password []byte
}