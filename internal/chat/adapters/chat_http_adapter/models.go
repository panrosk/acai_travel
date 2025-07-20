package chathttpadapter

import "github.com/go-playground/validator/v10"

type ChatRequestDTO struct {
	ConversationID string `json:"conversationId" validate:"required,uuid4"`
	Message        struct {
		Role    string `json:"role" validate:"required,oneof=user system"`
		Content string `json:"content" validate:"required"`
	} `json:"message" validate:"required"`
	UserID string `json:"userId" validate:"required,uuid4"`
}

type ValidationErrorResponse struct {
	FailedField string `json:"field"`
	Tag         string `json:"rule"`
	Value       string `json:"value,omitempty"`
}

var validate = validator.New()

func (dto *ChatRequestDTO) Validate() []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse

	err := validate.Struct(dto)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ValidationErrorResponse{
				FailedField: err.Namespace(), // ejemplo: "ChatRequestDTO.Message.Role"
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}
	return errors
}
