package internal

import "fmt"

type IntImpl struct {
	data int
}

func (i IntImpl) NotInInterface() {
	fmt.Printf("THis method is not in the interface\n")
}

func (i IntImpl) Update() {
	fmt.Println("Update")
}
