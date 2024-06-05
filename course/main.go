package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// struct represents a collection of fields
type User struct {
	Name     string
	Email    string
	Password string
}

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-5 + 12i)
)

func init() {
	// the init() function is executed before the main
	// and it is usually used to setup configurations or packages

	fmt.Println("Init...")
}

func main() {
	// main entrypoint

	// lowercase variables cannot be exported across packages
	// uppercase variables can be exported

	// use struct
	var user User = User{Name: "leo", Email: "email", Password: "password"}

	fmt.Printf("user has: %s, %s, %s\n", user.Name, user.Email, user.Password)

	// pointers refer to a variable in memory
	i := 42
	var p *int = &i // assign memory pointer of i to p
	fmt.Println(*p)
	*p = 0
	fmt.Println(i) // it is updated

	// rune represents int32
	var number rune = 100
	fmt.Println(number)

	// rune represents uint8
	var bs byte = 1
	fmt.Println(bs)

	// array has a static and defined length
	// slice has a dynamic size
	names := []string{"leo", "foo", "jack"}

	// loop in slices
	for key, value := range names {
		fmt.Printf("%d %s\n", key, value)
	}

	// for {
	// infinite loop
	// }

	a := 10
	var b int = 20
	var c = 5

	// if with short syntax
	if v := math.Pow(float64(a), float64(c)); v < 2000 {
		fmt.Println("first number is high")
	}

	if a > 5 {
		fmt.Println("first number is high")
	} else {
		fmt.Println("first number is low")
	}

	fmt.Println(a + b + c)

	var name string = "leo"
	var surname string = "pol"

	name, surname = swap(name, surname)
	fmt.Printf("swapped name %s and surname %s\n", name, surname)

	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	fmt.Println("hello world!")
}

// functions can return multiple values
func swap(x, y string) (string, string) {
	defer fmt.Println("deferring function!")
	return y, x
}
