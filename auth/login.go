package auth

import (
	"fmt"
	
	m "github.com/darylhjd/mangodex"
)

func Login(username, password string) (*m.DexClient, error) {
	// Create new client.
	// Without logging in, you may not be able to access 
	// all API functionality.
	c := m.NewDexClient()

	// Login using your username and password.
	err := c.Auth.Login(username, password)
	if err != nil {
		fmt.Println("Could not login!")
	}
	return c, err
}

