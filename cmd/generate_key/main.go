package main

import (
	"fmt"

	"github.com/rinnothing/simple-jwt/utils/jwt"
)

func main() {
	key := jwt.GenerateKey()
	fmt.Print("[]byte {")

	for i := 0; i < len(key); i++ {
		fmt.Print(key[i])
		if i != len(key)-1 {
			fmt.Print(", ")
		}
	}

	fmt.Println("}")
}
