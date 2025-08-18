package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// sanity checking the email address, parsing this all at startup and storing it for performance
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

//Generaly this is a map of potential errors and checks to see if any errors return true
//

// define new valid type which contains a map of validation errors for forms
type Validator struct {
	FieldErrors map[string]string
}

// Valid returns true if fielderrors map has NO entrys
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Add func adds error message to map
func (v *Validator) AddFieldError(key, message string) {
	//do need to init map first if not already
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// checfield adds err message
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// notblank returns true if value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// maxchars() returns true if value has no more than n chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// permmited returns true if value is in a list of allowed integers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// MinChars() returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches() returns true if a value matches a provided compiled regular
// expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
