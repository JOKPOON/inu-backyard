package validator

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	errs "github.com/team-inu/inu-backyard/entity/error"
	"github.com/team-inu/inu-backyard/internal/config"
)

type PayloadValidator interface {
	Validate(payload interface{}, ctx *fiber.Ctx) (bool, error)
	ValidateAuth(ctx *fiber.Ctx) (string, error)
}

type payloadValidator struct {
	validator *validator.Validate
	config    *config.AuthConfig
}

func NewPayloadValidator(config *config.AuthConfig) PayloadValidator {
	return &payloadValidator{
		validator: validator.New(),
		config:    config,
	}
}

func (v *payloadValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	sid := ctx.Cookies(v.config.Session.CookieName) // TODO: move cookie name to constant
	if sid == "" {
		return "", errs.New(errs.ErrAuthHeader, "missing auth header")
	}
	return sid, nil
}

func (v *payloadValidator) Validate(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if len(ctx.Body()) != 0 {
		if err := ctx.BodyParser(payload); err != nil {
			log.Println(err)
			return false, errs.New(errs.ErrBodyParser, err.Error())
		}
	}
	if err := ctx.ParamsParser(payload); err != nil {
		return false, errs.New(errs.ErrParamsParser, err.Error())
	}
	if err := ctx.QueryParser(payload); err != nil {
		return false, errs.New(errs.ErrQueryParser, err.Error())
	}
	if err := fileParser(payload, ctx); err != nil {
		return false, errs.New(errs.ErrBodyParser, err.Error())
	}

	if validationErrors := v.validateStruct(payload); validationErrors != nil {
		return false, errs.NewPayloadError(validationErrors)
	}

	return true, nil
}

func (v *payloadValidator) validateStruct(payload interface{}) []errs.ValidationErrorDetail {
	v.validator.RegisterTagNameFunc(getJSONFieldName)
	err := v.validator.Struct(payload)
	if err == nil {
		return nil
	}

	errors := err.(validator.ValidationErrors)
	var validationErrors []errs.ValidationErrorDetail
	for _, e := range errors {
		validationErrors = append(validationErrors, errs.ValidationErrorDetail{
			Field: e.Field(),
			Tag:   e.Tag(),
		})
	}
	return validationErrors
}

func getJSONFieldName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "" {
		return ""
	}
	return name
}

func fileParser(payload interface{}, ctx *fiber.Ctx) error {
	v := reflect.ValueOf(payload)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("interface must be a pointer to struct")
	}
	v = v.Elem() // Unwrap interface or pointer

	var form *multipart.Form

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fileKey := field.Tag.Get("file")

		if fileKey != "" {
			// Parse multipart form if not parsed before
			if form == nil {
				ctxForm, err := ctx.MultipartForm()
				if err != nil {
					return fmt.Errorf("cannot parse the file")
				}
				form = ctxForm
			}

			if field.Type == reflect.TypeOf((*multipart.File)(nil)).Elem() {
				// Parse a single file.
				// If the payload contains multiple files, the first file is being parsed
				for name, headers := range form.File {
					if name == fileKey {
						file, err := headers[0].Open()
						if err != nil {
							return fmt.Errorf("cannot parse the file")
						}
						v.Field(i).Set(reflect.ValueOf(file))
					}
				}
			} else if field.Type == reflect.TypeOf((*[]multipart.File)(nil)).Elem() {
				// Parse multiple files
				for name, headers := range form.File {
					if name != fileKey {
						continue
					}
					var files []multipart.File
					for _, header := range headers {
						file, err := header.Open()
						if err != nil {
							return fmt.Errorf("cannot parse the file")
						}
						files = append(files, file)
					}
					v.Field(i).Set(reflect.ValueOf(files))
					break
				}
			}
		}
	}
	return nil
}

func ToLower(payload interface{}) {
	v := reflect.ValueOf(payload)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("input must be a pointer to a struct")
	}

	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(strings.ToLower(field.String()))
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					ToLower(field.Index(j).Addr().Interface())
				}
			}
		case reflect.Struct:
			ToLower(field.Addr().Interface())
		}
	}
}
