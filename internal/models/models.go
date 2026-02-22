package models

import (
	"time"

	"github.com/google/uuid"
)

type ResponseResource struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	Favourite   bool      `json:"favourite"`
	ReadLater   bool      `json:"readLater"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ResponseUser struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Username      string    `json:"username"`
	IsAdmin       bool      `json:"isAdmin"`
	IsPro         bool      `json:"isPro"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
