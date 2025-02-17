package grpc

import (
	"context"
	"encoding/json"
	"errors"

	validator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidationField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationResult struct {
	Message string            `json:"message"`
	Fields  []ValidationField `json:"fields"`
}

//nolint:gochecknoglobals // we need this map to store the tag messages
var tagMessages = map[string]string{
	"required": "Required field",
	"email":    "Invalid email",
	"uuid":     "Invalid ID",
	"min":      "Value is smaller than minimum allowed",
	"max":      "Value is greater than maximum allowed",
}

func tagMessage(tag string) string {
	if msg, ok := tagMessages[tag]; ok {
		return msg
	}
	return "invalid value"
}

func (v ValidationResult) FromErrors(msg string, errs []validator.FieldError) ValidationResult {
	v.Message = msg
	for _, err := range errs {
		v.Fields = append(v.Fields, ValidationField{Field: err.Field(), Error: tagMessage(err.Tag())})
	}
	return v
}
func (v ValidationResult) JSON() string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal validation result")
		return ""
	}
	return string(b)
}

func ValidateInterceptor(validate *validator.Validate) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate.StructCtx(ctx, req); err != nil {
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(err, &invalidValidationError) {
				return nil, status.Error(
					codes.InvalidArgument,
					new(ValidationResult).FromErrors("Invalid request", nil).JSON(),
				)
			}
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return nil, status.Error(
					codes.InvalidArgument,
					new(ValidationResult).FromErrors("Invalid request", validationErrors).JSON(),
				)
			}
		}
		return handler(ctx, req)
	}
}
