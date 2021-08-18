package validatorhelper

import (
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

// MessageError DataStructure MessageError with field and message
type MessageError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorMessageTranslator function to customize the error message
func ErrorMessageTranslator(err error, translated ut.Translator) []MessageError {
	var messages []MessageError
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		var messageStruct MessageError
		messageStruct.Field = e.Field()
		messageStruct.Message = e.Translate(translated)
		messages = append(messages, messageStruct)
		// can translated each error one at a time.

	}
	return messages
}
