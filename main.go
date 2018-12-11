//Chat application is a application based on socket io
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/melvinodsa/chat-application/server"
)

//PORT for the application. One can provide the port through cmd argument flags
var PORT = flag.Int("port", 9000, "port for the chat application")

//STATIC is the static assets directory for the application.
//By default it is the one specified by the assets directory. If developer want to have
//to have something different, it can be done.
var STATIC = flag.String("static", "./assets", "static assets directory for the application")

//PATH is the path path variable used by the socket.io client.
//By default it is `/socket.io/`. Web requests will be handle using this path.
var PATH = flag.String("path", "/socket.io/", "custom path with which socket.io client has been initialized")

func main() {
	/*
	 * We will first parse the command line argument flags
	 * Then we will init the server
	 * Then we will set handler for the socket io requests to the socket io server
	 * Will serve the static contents
	 * Start the server
	 */
	//parsing the flags
	flag.Parse()

	//server initialization
	server.InitServer()

	//setting the socket io requests
	http.Handle(*PATH, server.GetServer())

	//static file contents
	http.Handle("/", http.FileServer(http.Dir(*STATIC)))

	//starting the server
	log.Println("Serving at localhost:", strconv.Itoa(*PORT), "...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*PORT), nil))
}
