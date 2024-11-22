package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/masena-dev/bookstore-api/internal/db"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Server struct {
	queries *db.Queries
	logger  *slog.Logger
}

func NewServer(logger *slog.Logger, db *db.Queries) *Server {
	return &Server{
		queries: db,
		logger:  logger,
	}
}

func (s *Server) NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Authors endpoints
	mux.HandleFunc("GET /api/v1/authors", s.handleListAuthors)
	mux.HandleFunc("GET /api/v1/authors/{id}/stats", s.handleGetAuthorStats) //Do this

	// Books endpoints
	mux.HandleFunc("GET /api/v1/books", s.handleListBooks)
	mux.HandleFunc("POST /api/v1/books", s.handleCreateBook)
	mux.HandleFunc("GET /api/v1/books/{id}", s.handleGetBook) 
	mux.HandleFunc("PUT /api/v1/books/{id}", s.handleUpdateBook)
	mux.HandleFunc("DELETE /api/v1/books/{id}", s.handleDeleteBook)

	return mux
}

func (s *Server) respond(w http.ResponseWriter, status int, data interface{}) {
	if data == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("failed to encode response", "error", err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.logger.Error("error", "error_message", message)
	errorResponse := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}
	s.respond(w, status, errorResponse)
}

// Authors handlers
func (s *Server) handleListAuthors(w http.ResponseWriter, r *http.Request) {
  authors, err := s.queries.ListAuthors(r.Context())
  if err != nil {
    s.respondError(w, http.StatusInternalServerError, "Failed fetching list of authors")
    return
  } 

  response := map[string]interface{}{
    "authors":authors,
  }

  s.respond(w, http.StatusOK, response)
}

func (s *Server) handleGetAuthorStats(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid author ID")
		return
	}
	s.respond(w, http.StatusNoContent, nil)
}

// Books handlers
func (s *Server) handleListBooks(w http.ResponseWriter, r *http.Request) {
  books, err := s.queries.ListBooks(r.Context())
  if err != nil {
    s.respondError(w, http.StatusInternalServerError, "Failed fetching list of books")
    return
  }

  response := map[string]interface{}{
    "books": books,
  }

  s.respond(w, http.StatusOK, response)
}

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := req.Validate(); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

  var description pgtype.Text
  description.String = req.Description
  description.Valid = req.Description != ""

  // Convert float64 to pgtype.Numeric for Price
  var price pgtype.Numeric
  // Create a decimal from the float with proper precision
  decimalStr := fmt.Sprintf("%.2f", req.Price)  // Format to 2 decimal places
  err := price.Scan(decimalStr)  // Scan the string representation
  if err != nil {
      s.respondError(w, http.StatusInternalServerError, "Failed to convert price")
      return
  }


  // Convert string to pgtype.Date for PublishedDate
  var publishedDate pgtype.Date
  date, err := time.Parse("2006-01-02", req.PublishedDate)
  if err != nil {
      s.respondError(w, http.StatusBadRequest, "Invalid date format")
      return
  }
  publishedDate.Time = date
  publishedDate.Valid = true

  book, err := s.queries.CreateBook(r.Context(), db.CreateBookParams{
    Title: req.Title,
    Isbn: req.ISBN,
    Description: description,
    Price: price,
    AuthorID: req.AuthorID,
    PublishedDate: publishedDate,
  })

  if err != nil {
    s.respondError(w, http.StatusInternalServerError, err.Error())
    return
  }

  response := map[string]interface{}{
    "message":"Book created successfully",
    "book": book,
  }

	s.respond(w, http.StatusOK, response)
}

func (s *Server) handleGetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := s.queries.GetBook(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.respondError(w, http.StatusNotFound, "Failed to find book")
			return
		}
		s.respondError(w, http.StatusInternalServerError, err.Error())
	}

  response := map[string]interface{}{
    "book": book,
  }

	s.respond(w, http.StatusOK, response)
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

  var req UpdateBookRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := req.Validate(); err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

  var description pgtype.Text
  description.String = req.Description
  description.Valid = req.Description != ""

  // Convert float64 to pgtype.Numeric for Price
  var price pgtype.Numeric
  // Create a decimal from the float with proper precision
  decimalStr := fmt.Sprintf("%.2f", *req.Price)  // Format to 2 decimal places
  err = price.Scan(decimalStr)  // Scan the string representation
  if err != nil {
      s.respondError(w, http.StatusInternalServerError, "Failed to convert price")
      return
  }

  // Convert string to pgtype.Date for PublishedDate
  var publishedDate pgtype.Date
  date, err := time.Parse("2006-01-02", req.PublishedDate)
  if err != nil {
      s.respondError(w, http.StatusBadRequest, "Invalid date format")
      return
  }
  publishedDate.Time = date
  publishedDate.Valid = true

  updatedBook, err := s.queries.UpdateBook(r.Context(), db.UpdateBookParams{
    Title: req.Title,
    Description:description,
    Price: price,
    PublishedDate: publishedDate,
    ID: id,
  })

  if err != nil {
    s.respondError(w, http.StatusInternalServerError, err.Error())
    return
  }

  response := map[string]interface{}{
    "message": "Book updated successfully",
    "book": updatedBook,
  }
	s.respond(w, http.StatusOK, response)
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

  err = s.queries.DeleteBook(r.Context(), id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.respondError(w, http.StatusNotFound, "Book not found")
			return
		}
		s.respondError(w, http.StatusInternalServerError, err.Error())
	}

	s.respond(w, http.StatusNoContent, "Book deleted successfully")
}
