package main

import (
	"net/http"
	"time"
)

type DiaryEntry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type statusCodeWriter struct {
	http.ResponseWriter
	statusCode int
}
