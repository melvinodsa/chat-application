//Package server has the basic server implementation fo the socket.io server
package server

import (
	"encoding/json"
	"log"
	"strings"

	sio "github.com/googollee/go-socket.io"
	"github.com/melvinodsa/chat-application/models"
)

//server to hold the server instance
var server *sio.Server

//InitServer will init the socket.io server.
//The server instance can be accessed by the GetServer function
func InitServer() {
	/*
	 * We will first intialize a server
	 * Then we will register the on connection event'
	 * Then we will register on error event
	 * Then we will set the initialized server instance to the global server instance
	 */
	//initing the server instance
	ser, err := sio.NewServer(nil)
	if err != nil {
		//handle the error. exit. since nothing to do ahead
		log.Fatal(err)
	}

	//initializing the on connection instance
	ser.On("connection", Connection)

	//handling the error
	ser.On("error", func(so sio.Socket, err error) {
		log.Println("error:", so.Id(), err)
	})

	//setting the server instance
	server = ser
}

//GetServer is the getter function for the server instance
func GetServer() *sio.Server {
	return server
}

//Connection sets the initializations to be done after basic connection setup
func Connection(so sio.Socket) {
	/*
	 * First we will handle the user join event to get the user name
	 * Then we will handle the message event
	 * Then we will handle the disconnect event
	 */
	//user join event
	so.On("user-join", func(msg string) {
		UserJoin(msg, so)
	})

	//message event
	so.On("message", func(msg string) {
		Message(msg, so)
	})

	//disconnect event
	so.On("disconnection", func() {
		user := (models.User{ID: so.Id()}).Get()
		if user != nil {
			so.BroadcastTo("chat-room", "disconnected", *user)
		}
		log.Println("on disconnect", so.Id())
	})
}

//Message is will send the message send by a user to the chat room
func Message(msg string, so sio.Socket) {
	/*
	 * We will create the message
	 * Then will broadcast the same
	 */
	user := (models.User{ID: so.Id()}).Get()
	m := models.Message{Msg: msg}
	if user != nil {
		m.User = *user
	}
	so.Emit("message", m)
	so.BroadcastTo("chat-room", "message", m)
}

//UserJoin jons a user to the application.
//It expect the message to be in the json encoded format of the user model.
//Refer for user model https://godoc.org/github.com/melvinodsa/chat-application/models#User
func UserJoin(msg string, so sio.Socket) {
	/*
	 * We will parse the message to User model
	 * Set the id of the user as the socket id
	 * Set the user
	 * Then will join the user to the chat room
	 * Will broad the user join message to the chat room
	 */
	//parsing the user model
	//user model
	var user models.User
	dec := json.NewDecoder(strings.NewReader(msg))
	if err := dec.Decode(&user); err != nil {
		//silent handle the error
		log.Println("Error while parsing the message on user join to user model", err.Error())
		return
	}

	//setting the id of the user
	user.ID = so.Id()
	log.Println("User has joined", user.Name)

	//setting the model
	user.Set()

	//joining the user to the chat room
	so.Join("chat-room")

	//broadcast the join message
	so.Emit("user-join", user)
	so.BroadcastTo("chat-room", "user-join", user)
}
