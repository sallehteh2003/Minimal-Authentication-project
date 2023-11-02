package Validation

import (
	"errors"
	"strings"
)

const (
	EnglishChars        = "qwertyuiopasdfghjklzxcvbnm"
	CapitalEnglishChars = "QWERTYUIOPASDFGHJKLZXCVBNM"
	Numbers             = "0123456789"
	ValidSpecialChars   = "!@#$%^&*_.-+=?"
	InvalidSpecialChars = "`~':;][{}(),<> \\ | / \" "

	LimitForName     = 100
	LimitForEmail    = 100
	LimitForPassword = 200
)

type Validation struct {
	ValidDomains []string
}
type MessageValidation struct {
	error    string
	Email    bool
	Password bool
	Name     bool
}

// Create a new instance of Validation

func CreateValidation(validDomain []string) *Validation {

	return &Validation{

		ValidDomains: validDomain,
	}
}
func (v *Validation) ValidateData(email string, name string, password string) (*MessageValidation, error) {
	na := v.validateName(name)
	em := v.validateEmail(email)
	pa := v.validatePassword(password)
	if !na || !em || !pa {
		return &MessageValidation{
			error:    "Validation error",
			Email:    em,
			Password: pa,
			Name:     na,
		}, errors.New("validation error")
	}

	return nil, nil
}
func (v *Validation) validateEmail(email string) bool {
	if email == "" || !(strings.Contains(email, "@")) || len(email) > LimitForEmail {
		return false
	}
	emailPart := strings.Split(email, "@")
	if len(emailPart) != 2 {
		return false
	}
	if !v.checkValidDomain(emailPart[1]) {
		return false
	}
	temp := false
	for i, i2 := range emailPart[0] {
		if i == 0 && string(i2) == "." {
			return false
		}

		// check for doable dot
		// example '..'
		if string(i2) == "." {
			if temp {
				return false
			}
			temp = true
		} else {
			temp = false
		}

		// check for English and Numbers words
		if !strings.Contains(EnglishChars, strings.ToLower(string(i2))) && !strings.Contains(Numbers, string(i2)) {
			if string(i2) == "." || string(i2) == "_" {
				continue
			}
			return false
		}

	}
	return true
}

func (v *Validation) validateName(name string) bool {
	if name == "" || len(name) > LimitForName {
		return false
	}
	for _, i2 := range name {
		if !strings.Contains(EnglishChars, strings.ToLower(string(i2))) {
			return false
		}
	}

	return true
}

func (v *Validation) validatePassword(password string) bool {
	if password == "" || len(password) < 8 || len(password) > LimitForPassword {
		return false
	}
	EnChar := false
	EnCharCapital := false
	numbers := false
	SpCharValid := false

	for _, i2 := range password {
		// check for invalid char
		if strings.Contains(InvalidSpecialChars, string(i2)) {
			return false
		}

		if strings.Contains(EnglishChars, string(i2)) { // check for EnglishChars
			EnChar = true
			continue
		} else if strings.Contains(CapitalEnglishChars, string(i2)) { // check for CapitalEnglishChars
			EnCharCapital = true
			continue
		} else if strings.Contains(Numbers, string(i2)) { // check for Numbers
			numbers = true
			continue
		} else if strings.Contains(ValidSpecialChars, string(i2)) { // check for SpecialChar Valid
			SpCharValid = true
			continue
		} else {
			return false
		}
	}

	return EnChar && EnCharCapital && numbers && SpCharValid
}

func (v *Validation) checkValidDomain(domain string) bool {
	for _, d := range v.ValidDomains {
		if d == domain {
			return true
		}
	}
	return false
}
