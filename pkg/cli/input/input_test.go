package input

import "testing"

func TestValidateRequired(t *testing.T) {
	input := "test"
	valid := validate(input, "required")
	if !valid {
		t.Errorf("Validate(%s) valid = %v; wanted = %v", input, valid, true)
	}
}

func TestValidateRequiredFailure(t *testing.T) {
	input := ""
	valid := validate(input, "required")
	if valid {
		t.Errorf("Validate(%s) valid = %v; wanted = %v", input, valid, false)
	}
}

func TestValidateLength(t *testing.T) {
	input := "12345679"
	valid := validate(input, "length:8")
	if !valid {
		t.Errorf("Validate(%s) valid = %v; wanted = %v", input, valid, true)
	}
}

func TestValidateRequiredLength(t *testing.T) {
	input := "12345678"
	valid := validate(input, "length:8")
	if !valid {
		t.Errorf("Validate(%s) valid = %v; wanted = %v", input, valid, true)
	}
}

func TestValidateRequiredLengthFailure(t *testing.T) {
	input := "1234567"
	valid := validate(input, "required,length:8")
	if valid {
		t.Errorf("Validate(%s) valid = %v; wanted = %v", input, valid, true)
	}
}
