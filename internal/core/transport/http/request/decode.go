package core_http_request

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	core_errors "github.com/severholod/amazing-todo/internal/core/errors"
	"net/http"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("error decoding json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("error validating request: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
