package models

/*
 * This file contains the model definition for
 */

//Message wraps the message send by the user
type Message struct {
	//Msg is the message content
	Msg string
	//User hwo has sent the message
	User User
}
