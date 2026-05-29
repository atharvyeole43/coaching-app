package utils

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	english "github.com/go-playground/validator/v10/translations/en"
	"github.com/sirupsen/logrus"
)

// InitializeValidator initializes the validator with custom translations and validation rules.
func InitializeValidator() (*validator.Validate, ut.Translator) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		logrus.Error("Translator not found")
	}

	tokenValidator := validator.New()

	// Register default translations for English
	if err := english.RegisterDefaultTranslations(tokenValidator, trans); err != nil {
		logrus.Error(err)
	}

	// Register custom translations and validation rules here
	_ = tokenValidator.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = tokenValidator.RegisterTranslation("required_if", trans, func(ut ut.Translator) error {
		return ut.Add("required_if", "{0} is a required.", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required_if", fe.Field())
		return t
	})

	tokenValidator.RegisterTranslation("min", trans, func(ut ut.Translator) error {
		return ut.Add("min", "The {0} field must be at least {1} characters long.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), fe.Param())
		return t
	})

	tokenValidator.RegisterTranslation("max", trans, func(ut ut.Translator) error {
		return ut.Add("max", "The {0} field cannot be longer than {1} characters.", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), fe.Param())
		return t
	})

	return tokenValidator, trans
}

// ValidateRequest performs validation on a request using the provided validator and translator.
func ValidateRequest(c *gin.Context, requestBody interface{}) interface{} {
	tokenValidator, trans := InitializeValidator() // Initialize the validator and translator

	tokenValidatorErr := tokenValidator.Struct(requestBody)

	if tokenValidatorErr != nil {
		var errorMessage string
		var counter int

		tokenValidatorErrors, ok := tokenValidatorErr.(validator.ValidationErrors)
		if !ok {
			logrus.Error("Failed to assert ValidationErrors")
			return "Validation error"
		}

		size := len(tokenValidatorErrors)

		for _, e := range tokenValidatorErrors {
			counter++
			errorMessage = errorMessage + e.Translate(trans) // Use the provided translator
			if counter != size {
				errorMessage = errorMessage + "|"
			}
		}

		logrus.Error("errorMessage", errorMessage)
		return errorMessage
	}

	return nil
}

func IsValidName(name string) bool {
	for _, char := range name {
		if (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') && char != ' ' {
			return false
		}
	}
	return true
}

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
func IsValidMobile(mobile string) bool {
	mobileRegex := regexp.MustCompile(`^[0-9]{10}$`)
	return mobileRegex.MatchString(mobile)
}

func ValidateField(value interface{}, field map[string]interface{}) error {
	valMap, _ := field["validation"].(map[string]interface{})
	label, _ := field["label"].(string)
	fieldType, _ := field["type"].(string)

	strVal := fmt.Sprintf("%v", value)

	if r, ok := valMap["required"].(bool); ok && r {
		if strVal == "" || strVal == "<nil>" {
			return fmt.Errorf("%s is required", label)
		}
	}

	if strVal == "" || strVal == "<nil>" {
		return nil
	}

	if minL, ok := valMap["minLength"].(float64); ok && len(strVal) < int(minL) {
		return fmt.Errorf("%s must be at least %d characters", label, int(minL))
	}
	if maxL, ok := valMap["maxLength"].(float64); ok && len(strVal) > int(maxL) {
		return fmt.Errorf("%s cannot exceed %d characters", label, int(maxL))
	}

	if pattern, ok := valMap["pattern"].(string); ok && pattern != "" {
		matched, _ := regexp.MatchString(pattern, strVal)
		if !matched {
			return fmt.Errorf("%s has invalid format", label)
		}
	}

	switch fieldType {
	case "email":
		if _, err := mail.ParseAddress(strVal); err != nil {
			return fmt.Errorf("%s must be a valid email", label)
		}

	case "number":
		switch v := value.(type) {
		case float64:
			if v != float64(int(v)) {
				return fmt.Errorf("%s must be a valid integer", label)
			}
		case string:
			if _, err := strconv.Atoi(v); err != nil {
				return fmt.Errorf("%s must be a valid integer", label)
			}
		default:
			return fmt.Errorf("%s must be a valid integer", label)
		}

	case "decimal":
		switch v := value.(type) {
		case float64:
		case string:
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				return fmt.Errorf("%s must be a valid decimal", label)
			}
		default:
			return fmt.Errorf("%s must be a valid decimal", label)
		}

	case "date":
		if _, err := time.Parse("2006-01-02", strVal); err != nil {
			return fmt.Errorf("%s must be in YYYY-MM-DD format", label)
		}

	case "datetime":
		if _, err := time.Parse("2006-01-02 15:04:05", strVal); err != nil {
			return fmt.Errorf("%s must be in YYYY-MM-DD HH:MM:SS format", label)
		}

	case "time":
		if _, err := time.Parse("15:04:05", strVal); err != nil {
			return fmt.Errorf("%s must be in HH:MM:SS format", label)
		}

	case "checkbox":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%s must be a boolean", label)
		}

	case "url", "file":
		if _, err := url.ParseRequestURI(strVal); err != nil {
			return fmt.Errorf("%s must be a valid URL", label)
		}

	case "select", "radio":
		if options, ok := field["options"].([]interface{}); ok && len(options) > 0 {
			found := false
			for _, opt := range options {
				if fmt.Sprintf("%v", opt) == strVal {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("%s has invalid option value", label)
			}
		}
	}

	return nil
}
