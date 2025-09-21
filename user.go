package main

import "time"

type User struct {
	// TODO: UUID
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
