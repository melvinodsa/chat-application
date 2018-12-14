//Package config has the implementation for loading the config required by the application
//from command line and config file
package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/melvinodsa/chat-application/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

/*
 * This file has the implementation for loading the config from config file and from command line flags
 */

//PORT for the application. One can provide the port through cmd argument flags
var PORT = flag.Int("port", 9000, "port for the chat application")

//STATIC is the static assets directory for the application.
//By default it is the one specified by the assets directory. If developer want to have
//to have something different, it can be done.
var STATIC = flag.String("static", "assets", "static assets directory for the application")

//PATH is the path path variable used by the socket.io client.
//By default it is `/socket.io/`. Web requests will be handle using this path.
var PATH = flag.String("path", "/socket.io/", "custom path with which socket.io client has been initialized")

//GOOGLEAUTHCALLBACK is the google social authentication callback url.
//By default it is `/googleauthcallback`. After authentication google oauth redirects to the url specified here
var GOOGLEAUTHCALLBACK = flag.String("google-auth-callback", "http://localhost:9000", "callback url for google authentication")

//GOOGLEEMAILSCOPE is the email scope url of google auth
var GOOGLEEMAILSCOPE = flag.String("google-email-scope", "https://www.googleapis.com/auth/userinfo.email", "email scope url for google auth")

//GOOGLEPROFILESCOPE is the email scope url of
var GOOGLEPROFILESCOPE = flag.String("google-profile-scope", "https://www.googleapis.com/auth/userinfo.profile", "profile scope url for google auth")

//GOOGLEAUTHCLIENTID is the google auth client id
var GOOGLEAUTHCLIENTID = flag.String("google-auth-clientid", "<It's secret>", "email scope url for google auth")

//GOOGLEAUTHSECRETID is the google auth secret id
var GOOGLEAUTHSECRETID = flag.String("google-auth-secretid", "<It's secret>", "profile scope url for google auth")

//CONFIGFILEPATH is the path to the config file for thr application
var CONFIGFILEPATH = flag.String("config-file", "./config.json", "config file location")

//Config has the basic config structure required by the application.
//config file however will have priority over Command line args
type Config struct {
	Port               int    //Port of the application
	Static             string //Static files location
	Path               string //Path to be used as url for socket.io
	GoogleAuthCallback string //Google authentication callback url
	GoogleEmailScope   string //Google authentication email scope url
	GoogleProfileScope string //Google authentication profile scope url
	GoogleAuthClientid string //Google authentication client id
	GoogleAuthSecretid string //Google authentication secret id
}

//LoadConfig will load the configuratoions required by the application
func LoadConfig() {
	/*
	 * We will parse the flags
	 * Then we will load the config file
	 * Then we will init the google oauth
	 */

	//parsing the flags
	flag.Parse()

	//loading the config file
	loadConfigFile()

	//initing the google auth
	models.GOOGLE = &oauth2.Config{
		RedirectURL:  *GOOGLEAUTHCALLBACK + "/googleauthcallback",
		ClientID:     *GOOGLEAUTHCLIENTID,
		ClientSecret: *GOOGLEAUTHSECRETID,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			*GOOGLEEMAILSCOPE,
			*GOOGLEPROFILESCOPE,
		},
	}
}

func loadConfigFile() {
	/*
	 * Opening the config file
	 * Then we will decode the config file
	 * If the command line flags are not empty update them with config
	 */
	//Loading the config file
	file, err := os.Open(*CONFIGFILEPATH)
	if err != nil {
		log.Println("No config file found", err.Error())
		return
	}
	dec := json.NewDecoder(file)

	//decode the config
	config := Config{}
	err = dec.Decode(&config)
	if err != nil {
		log.Fatal("Error while reading the config file", err.Error())
	}

	//updating the command line args based on the config file
	//port
	if config.Port != 0 {
		*PORT = config.Port
	}

	//static
	if len(config.Static) != 0 {
		*STATIC = config.Static
	}

	//path
	if len(config.Path) != 0 {
		*PATH = config.Path
	}

	//google auth callback
	if len(config.GoogleAuthCallback) != 0 {
		*GOOGLEAUTHCALLBACK = config.GoogleAuthCallback
	}

	//google email scope
	if len(config.GoogleEmailScope) != 0 {
		*GOOGLEEMAILSCOPE = config.GoogleEmailScope
	}

	//google profile scope
	if len(config.GoogleProfileScope) != 0 {
		*GOOGLEPROFILESCOPE = config.GoogleProfileScope
	}

	//google auth client id
	if len(config.GoogleAuthClientid) != 0 {
		*GOOGLEAUTHCLIENTID = config.GoogleAuthClientid
	}

	//google auth secret id
	if len(config.GoogleAuthSecretid) != 0 {
		*GOOGLEAUTHSECRETID = config.GoogleAuthSecretid
	}
}
