package utils

import (
	"demo/internal/app/errs"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/dlclark/regexp2"
	"io"
)

func init() {
	govalidator.TagMap["username_valid"] = govalidator.Validator(func(str string) bool {
		return govalidator.Matches(str, "^[A-Za-z0-9]{1,20}$")
	})
	govalidator.TagMap["valid_pass"] = govalidator.Validator(func(str string) bool {
		reg, _ := regexp2.Compile(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#%^&*?\.\/\-\=]).{8,20}$`, 0)
		match, _ := reg.FindStringMatch(str)
		return match != nil
	})
	govalidator.TagMap["valid_email"] = govalidator.Validator(func(str string) bool {
		return govalidator.IsEmail(str)
	})
}

// ParseRequestData 请求数据校验
func ParseRequestData(r io.Reader, data any) error {
	if err := json.NewDecoder(r).Decode(data); err != nil {
		return errs.EcInvalidRequest
	}
	return Validator(data)
}

// Validator 数据校验
func Validator(data any) error {
	govalidator.SetFieldsRequiredByDefault(true)
	if result, err := govalidator.ValidateStruct(data); !result && err != nil {
		for _, e := range err.(govalidator.Errors).Errors() {
			return e
		}
	}
	return nil
}
