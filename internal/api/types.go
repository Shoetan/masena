package api

import "time"

// Response types
type AuthorsResponse struct {
	Authors []Author `json:"authors"`
}

type BooksResponse struct {
	Books []Book `json:"books"`
}

type CreateBookResponse struct {
	Message string `json:"message"`
	Book    Book   `json:"book"`
}

type GetBookResponse struct {
	Book Book `json:"book"`
}

type UpdateBookResponse struct {
	Message string `json:"message"`
	Book    Book   `json:"book"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Domain types
type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	ISBN          string    `json:"isbn"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	Author        Author    `json:"author"`
	PublishedDate string    `json:"published_date"` // Format: YYYY-MM-DD
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AuthorStats struct {
	ID                  int64          `json:"id"`
	Name                string         `json:"name"`
	TotalBooks          int            `json:"total_books"`
	AverageBookPrice    float64        `json:"average_book_price"`
	EarliestPublication string         `json:"earliest_publication"` // YYYY-MM-DD
	LatestPublication   string         `json:"latest_publication"`   // YYYY-MM-DD
	TotalRevenue        float64        `json:"total_revenue"`
	BooksByYear         map[string]int `json:"books_by_year"` // e.g. {"2023": 3, "2024": 1}
}

// Request types
type CreateBookRequest struct {
	Title         string  `json:"title"`
	ISBN          string  `json:"isbn"`
	Description   string  `json:"description,omitempty"`
	Price         float64 `json:"price"`
	AuthorID      int64   `json:"author_id"`
	PublishedDate string  `json:"published_date"` // Format: YYYY-MM-DD
}

type UpdateBookRequest struct {
	Title         string   `json:"title,omitempty"`
	Description   string   `json:"description,omitempty"`
	Price         *float64 `json:"price,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"` // Format: YYYY-MM-DD
}

// Validation methods
func (r CreateBookRequest) Validate() error {
	if r.Title == "" || len(r.Title) > 255 {
		return ValidationError{Field: "title", Message: "must be between 1 and 255 characters"}
	}

	if !isValidISBN(r.ISBN) {
		return ValidationError{Field: "isbn", Message: "must be a valid ISBN format"}
	}

	if r.Price < 0 {
		return ValidationError{Field: "price", Message: "must be greater than or equal to 0"}
	}

	if !isValidDate(r.PublishedDate) {
		return ValidationError{Field: "published_date", Message: "must be in YYYY-MM-DD format"}
	}

	return nil
}

func (r UpdateBookRequest) Validate() error {
	if r.Title != "" && len(r.Title) > 255 {
		return ValidationError{Field: "title", Message: "must be between 1 and 255 characters"}
	}

	if r.Price != nil && *r.Price < 0 {
		return ValidationError{Field: "price", Message: "must be greater than or equal to 0"}
	}

	if r.PublishedDate != "" && !isValidDate(r.PublishedDate) {
		return ValidationError{Field: "published_date", Message: "must be in YYYY-MM-DD format"}
	}

	return nil
}

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// Helper functions for validation
func isValidISBN(isbn string) bool {
	// Simple check for now - feel free to improve this
	if len(isbn) < 10 || len(isbn) > 17 {
		return false
	}
	for _, r := range isbn {
		if r != '-' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
