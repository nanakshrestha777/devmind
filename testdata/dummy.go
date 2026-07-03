package testdata

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u *User) Greet() {
	fmt.Printf("Hello, %s\n", u.Name)
}
