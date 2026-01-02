package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Mahesh252k/students-api/internal/storage"
	"github.com/Mahesh252k/students-api/internal/types"
	"github.com/Mahesh252k/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		// 1. Decode JSON
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body not allowed")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 2. Request validation
		if err := validator.New().Struct(student); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidatorError(validatorErr))
			return
		}

		// 3. CALL THE DATABASE to create student
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			// If MySQL fails, return 500
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// 4. Send success response with the new ID
		response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"success": "OK",
			"id":      lastId,
		})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("GetById - getting a student", slog.String("id", id))

		parsedId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format")))
			return
		}
		student, err := storage.GetStudentById(parsedId)
		if err != nil {
			slog.Error("GetById - failed to get student", slog.String("id", id))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
