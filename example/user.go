package main

import "errors"

type userStatus string

func (us userStatus) Enum() []interface{} {
	return []interface{}{
		"new",
		"approved",
		"active",
		"deleted",
	}
}

// A demo app that receives data from http and stores it in memory.

type User struct {
	FirstName string     `json:"firstName" required:"true" title:"First name" minLength:"3"`
	LastName  string     `json:"lastName" required:"true" title:"Last name" minLength:"3"`
	Locale    string     `json:"locale" title:"User locale" enum:"ru-RU,en-US"`
	Age       int        `json:"age" title:"Age" minimum:"1"`
	Status    userStatus `json:"status" title:"Status"`
	Bio       string     `json:"bio" title:"Bio" description:"A brief description of the person." formType:"textarea"`
}

func (User) Title() string {
	return "User"
}

func (User) Description() string {
	return "User is a sample entity."
}

type userRepo struct {
	st         []User
	schemaName string
}

func (r *userRepo) create(u User) {
	r.st = append(r.st, u)
}

func (r *userRepo) update(i int, u User) {
	r.st[i-1] = u
}

func (r userRepo) list() []User {
	return r.st
}

func (r userRepo) get(i int) (User, error) {
	if i-1 > len(r.st) {
		return User{}, errors.New("user not found")
	}

	return r.st[i-1], nil
}
