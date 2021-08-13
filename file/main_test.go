package file

import "testing"

func Test_InputParse(t *testing.T) {
	user_input := LoadUserInput("test_data.txt")

	if user_input.Username != "ufukty" {
		t.Errorf("LoadUserInput function couldn't parse Username correctly.\nOutput: %s", user_input.Username)
	}

	if user_input.Password != "Hello World!" {
		t.Errorf("LoadUserInput function couldn't parse Password correctly.\nOutput: %s", user_input.Password)
	}

	if user_input.Otp_nonce != "111111" {
		t.Errorf("LoadUserInput function couldn't parse Otp_nonce correctly.\nOutput: %s", user_input.Otp_nonce)
	}
}
