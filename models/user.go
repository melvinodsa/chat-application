//Package models has all the models required for the chat application
package models

import "golang.org/x/oauth2"

/*
 * This file contains the user model
 */

var users = map[string]User{}

//UserInfoMap contains the string of the keywords to be used to retrieve the
//user info from the api respose of the auth agent like Google, Facebook etc.
type UserInfoMap struct {
	Email   string //Email key of the user info model
	Name    string //Name key of the user info model
	Picture string //Picture url key of the user info model
}

//GOOGLEUSERINFO has the map of the key words to be used to parse the google user info api
var GOOGLEUSERINFO = UserInfoMap{"email", "name", "picture"}

//GOOGLEUSERURL has the url to be hit to fetch the user info from Google
var GOOGLEUSERURL = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="

//GOOGLE has the oauth configuration for google authentication
var GOOGLE *oauth2.Config

//User is the model representing a user
type User struct {
	//Name of the user
	Name string
	//ID of the user
	ID string
	//Email id of the user
	Email string
	//Picture is the profile image of the user
	Picture string
	//AccessToken is the access token for accessing te user information from the google apis
	AccessToken string
}

//Get is the getter for the user instance for a given User ID
func (u User) Get() *User {
	u, ok := users[u.ID]
	if !ok {
		return nil
	}
	return &u
}

//Set is the setter for the given User ID
func (u User) Set() {
	users[u.ID] = u
}

//Delete the user with given user ID
func (u User) Delete() {
	delete(users, u.ID)
}
