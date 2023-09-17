package main

import "fmt"

// When
// 1. When we need to update state
// 2. When we want to optimize the memory for Large objects that are getting called a lot

type User struct {
	email    string
	username string
	age      int
	file     []byte // ?? Large ??
}

// use pointer here so we don't have to return an empty initialized user object
func getUser() (*User, error) {
	// return User{}, fmt.Errorf("foo")
	return nil, fmt.Errorf("foo")
}

// syntactic sugar
func (u User) Email() string {
	return u.email
}

// same as above
// user is x amount of bytes => sizeOf(user)
// copying over user to this function will copy everything in there
// for example file can be 1gb in size and its copied
// so use a pointer so it just 8 bytes
func Email(user User) string {
	return user.email
}

// syntactic sugar
// need memory address/reference to update the state of user
func (u *User) updateEmail(email string) {
	u.email = email
}

// same as above
// user is 8 bytes because of pointer
func UpdateEmail(u *User, email string) {
	u.email = email
}

func main() {
	user := User{
		email: "agg@foo.com",
	}
	fmt.Println(Email(user))

	user.updateEmail("goo@bar.com")
	fmt.Println(user.Email())

	UpdateEmail(&user, "foo@bar.com")
	fmt.Println(user.Email())
}
