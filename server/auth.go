package server

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/melvinodsa/chat-application/models"
	"golang.org/x/oauth2"
)

/*
 * This file contains the social authentication implementations
 */

type Urls struct {
	Google string
}

//GoogleAuthCallback is invoked by google oauth agent as callback when the user is authenticated.
//This is a standard HandleFunc implmentation of the http.HandleFunc
func GoogleAuthCallback(res http.ResponseWriter, req *http.Request) {
	/*
	 * We will load the template for the chatroom
	 * We will get the code param from the request
	 * Then we will get the auth token from the google auth with the code
	 * Then we will get the user profile info from google apis
	 * Will give response with the template and user info
	 */
	//paring the template
	t := template.New("chatroom.html")                   // Create a template.
	te, err := t.ParseFiles("./templates/chatroom.html") // Parse template file.
	if err != nil {
		log.Println("Error while loading the chatroom template file", err.Error())
	}

	//getting the code parameter from the request
	code, ok := req.URL.Query()["code"]
	if !ok || len(code[0]) < 1 {
		log.Println("Couldn't find the code param in the google auth callback")
		te.Execute(res, struct {
			Error error
			User  *models.User
		}{errors.New("Couldn't find the code from google auth"), nil})
		return
	}

	//get the auth token from the google auth with the code
	tok, err := models.GOOGLE.Exchange(oauth2.NoContext, code[0])
	if err != nil {
		log.Println("Couldn't get the auth token from the google auth callback", err.Error())
		te.Execute(res, struct {
			Error error
			User  *models.User
		}{errors.New("Couldn't fetch the auth token from google auth"), nil})
		return
	}

	//fetching the user profile from google apis
	user := &models.User{AccessToken: tok.AccessToken}
	user, err = GetGoogleUserInfo(user)
	if err != nil {
		log.Println(err.Error())
		te.Execute(res, struct {
			Error error
			User  *models.User
		}{errors.New("Couldn't fetch the user profile from google auth"), nil})
		return
	}

	//returning the response
	te.Execute(res, struct {
		Error error
		User  *models.User
	}{nil, user})
}

//LoginPage handles the all the request to show the social auth login page.
//We will set the social auth urls in this handler.
func LoginPage(res http.ResponseWriter, req *http.Request) {
	/*
	 * We will return the auth urls for Google
	 * Then we will parse the template for login page
	 * Will return the response as the template page
	 */
	urls := Urls{
		Google: models.GOOGLE.AuthCodeURL("state", oauth2.AccessTypeOffline),
	}

	t := template.New("login.html")                   // Create a template.
	te, err := t.ParseFiles("./templates/login.html") // Parse template file.
	if err != nil {
		log.Println("Error while loading the login template file", err.Error())
	}
	err = te.Execute(res, urls) // merge.
	if err != nil {
		log.Println("Error while executing the login template", err.Error())
	}
}

//GetGoogleUserInfo returns the returns google user's info
func GetGoogleUserInfo(u *models.User) (*models.User, error) {
	/*
		We will hit the google api for user info
		If any error happens we will return after logging it
		Then we will parse the api response
		Then will set the properties based on the google user info api's response
		We will set the user model email with the info email
	*/
	//hitting the google apis
	resp, err := http.Get(models.GOOGLEUSERURL + url.QueryEscape(u.AccessToken))
	if err != nil {
		//handling the error while getting the info from the google
		log.Println("Error while getting the userinfo from the google")
		return nil, err
	}
	defer resp.Body.Close()

	//parsing the api response
	ui := map[string]interface{}{}
	//handling the error while parsing the api response
	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		log.Println("Error while parsing the userinfo api response from google")
		return nil, err
	}

	//setting the info from the api response in the userinfo model
	info := &models.User{}
	//setting the name
	v, ok := ui[models.GOOGLEUSERINFO.Name]
	if !ok {
		log.Println("Error while getting the user's name from the google user info api")
		return nil, errors.New("Key not found " + models.GOOGLEUSERINFO.Name)
	}
	info.Name = v.(string)
	//setting the email
	v, ok = ui[models.GOOGLEUSERINFO.Email]
	if !ok {
		log.Println("Error while getting the user's email from the google user info api")
		return nil, errors.New("Key not found " + models.GOOGLEUSERINFO.Email)
	}
	info.Email = v.(string)
	//setting the picture
	v, ok = ui[models.GOOGLEUSERINFO.Picture]
	if !ok {
		log.Println("Error while getting the user's pciture from the google user info api")
		return nil, errors.New("Key not found " + models.GOOGLEUSERINFO.Picture)
	}
	info.Picture = v.(string)

	//setting the email of the user with info email
	u.Email = info.Email
	return info, nil
}
