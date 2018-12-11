//Package models has all the models required for the chat application
package models

/*
 * This file contains the user model
 */

var users = map[string]User{}

//User is the model representing a user
type User struct {
	//Name of the user
	Name string
	//ID of the user
	ID string
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
