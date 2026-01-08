package main

import (
	"fmt"
	"time"
)

// User represents a simple user structure
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// Greet returns a greeting message for the user
func (u *User) Greet() string {
	return fmt.Sprintf("Hello, %s! Your ID is %d", u.Name, u.ID)
}

// ValidateEmail checks if the user has a valid email
func (u *User) ValidateEmail() bool {
	if u.Email == "" {
		return false
	}
	// Simple validation - just check if @ exists
	for i := 0; i < len(u.Email); i++ {
		if u.Email[i] == '@' {
			return true
		}
	}
	return false
}

// CalculateAge calculates user age based on creation date
func (u *User) CalculateAge() int {
	diff := time.Since(u.CreatedAt)
	return int(diff.Hours() / 24 / 365)
}

func main() {
	user := User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
	}

	fmt.Println(user.Greet())
	fmt.Printf("Email valid: %v\n", user.ValidateEmail())
	fmt.Printf("Account age: %d years\n", user.CalculateAge())
}
