package task

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Input struct {
	Title  *string `json:"title"`
	Status *Status `json:"status"`
}

func (i *Input) Validate(required bool) error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Title,
			validation.When(required || i.Title != nil, validation.Required.Error("обязательное поле")),
			validation.When(i.Title != nil, validation.Length(3, 200).Error("должно быть 3-200 символов")),
		),
		validation.Field(&i.Status,
			validation.When(required || i.Status != nil, validation.Required.Error("обязательное поле")),
			validation.When(i.Status != nil, validation.In(StatusTodo, StatusDoing, StatusDone).Error("допустимые статусы: todo, doing, done")),
		),
	)
}
