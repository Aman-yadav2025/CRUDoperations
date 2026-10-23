package student

import (
	"CRUDapp/internal/storage"
	"CRUDapp/internal/types"
	"CRUDapp/internal/utils/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	//"golang.org/x/mod/sumdb/storage"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Initiating Student Creation")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validate the request

		if err := validator.New().Struct(student); err != nil {
			var validateErrs validator.ValidationErrors
			if errors.As(err, &validateErrs) {
				response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
				return
			}

			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		Lastid, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			slog.Error("user Creatation failed")
			return
		}

		slog.Info("user Created successfully", slog.String("userID", fmt.Sprint(Lastid)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": Lastid})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting the student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			slog.Error("can't convert id to int64", slog.String("id",id))
			return
		}
		student, e := storage.GetStudentById(intId)

		if e != nil {
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(e))
			slog.Error("Can't get the student with id",slog.String("id",id))
			return
		}
		response.WriteJson(w,http.StatusOK,student)
	}
}
