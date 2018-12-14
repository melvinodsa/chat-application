//Chat application is a application based on socket io
package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/melvinodsa/chat-application/config"
	"github.com/melvinodsa/chat-application/server"
)

func main() {
	/*
	 * We will first load the configs
	 * Then we will init the server
	 * Then we will set handler for the socket io requests to the socket io server
	 * Social authentication google callback url
	 * Login page as index page will be served
	 * Will serve the static contents
	 * Start the server
	 */
	//Loading the configs
	config.LoadConfig()

	//server initialization
	server.InitServer()

	//setting the socket io requests
	http.Handle(*config.PATH, server.GetServer())

	//Login page
	http.HandleFunc("/", server.LoginPage)

	//Google social authentication callback
	http.HandleFunc("/googleauthcallback", server.GoogleAuthCallback)

	//static file contents
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(*config.STATIC))))

	//starting the server
	log.Println("Serving at localhost:", strconv.Itoa(*config.PORT), "...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*config.PORT), nil))
}
