package main

import (
"fmt"
"os"

"golang.org/x/crypto/bcrypt"
)

func main() {
if len(os.Args) < 2 {
fmt.Println("Usage: go run main.go <password>")
os.Exit(1)
}

password := os.Args[1]
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
fmt.Printf("Error generating hash: %v\n", err)
os.Exit(1)
}

fmt.Printf("Password: %s\n", password)
fmt.Printf("Hash:     %s\n", string(hash))

// Verify it works
err = bcrypt.CompareHashAndPassword(hash, []byte(password))
if err == nil {
fmt.Println("✅ Hash verified successfully!")
} else {
fmt.Printf("❌ Hash verification failed: %v\n", err)
}
}
