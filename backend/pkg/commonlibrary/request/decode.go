package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Validator interface {
	Validate() error
}

func DecodeAndValidate(src io.ReadCloser, target Validator) error {
	if src == nil {
		return errors.New("empty body")
	}

	defer func(src io.ReadCloser) {
		_ = src.Close()
	}(src)

	if err := json.NewDecoder(src).Decode(target); err != nil {
		return err
	}

	return target.Validate()
}

func Decode(src io.ReadCloser, target any) error {
	if src == nil {
		return errors.New("empty body")
	}

	defer func(src io.ReadCloser) {
		_ = src.Close()
	}(src)

	return json.NewDecoder(src).Decode(target)
}

func PathUUID(w http.ResponseWriter, r *http.Request, key string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}
