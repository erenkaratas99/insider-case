package entities

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"regexp"
)

const phoneRegexPattern = "^\\+90[5-9]\\d{9}$"

var phoneRegex = regexp.MustCompile(phoneRegexPattern)

type Validator interface {
	Validate(req interface{}) error
}

type validator struct{}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(req interface{}) error {
	switch r := req.(type) {
	case *MessageInfo:
		return r.Validate()
	case *GetAllRequest:
		return r.Validate()
	default:
		return errors.New("unsupported type for validation")
	}
}

func (m *MessageInfo) Validate() error {
	if m.To == "" || !phoneRegex.MatchString(m.To) {
		log.Infof("invalid phone number for message Id %s : %s", m.Id, m.To)
		return errors.New("invalid phone number")
	}
	if len(m.Content) < 1 || len(m.Content) > 1000 {
		log.Infof("invalid content length for message Id : %s", m.Id)
		return errors.New("invalid content length")
	}
	return nil
}

func (r *GetAllRequest) Validate() error {
	if r.Limit > 50 {
		r.Limit = 50
	}
	if r.Offset < 0 {
		r.Offset = 0
	}
	return nil
}
