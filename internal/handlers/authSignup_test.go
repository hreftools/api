package handlers_test

import (
	"strings"
	"testing"

	"github.com/hreftools/api/internal/handlers"
)

func TestAuthSignupBody_Normalize(t *testing.T) {
	tests := []struct {
		name     string
		input    handlers.AuthSignupBody
		expected handlers.AuthSignupBody
	}{
		{
			name: "Produces no change to already normalized input",
			input: handlers.AuthSignupBody{
				Username: "user_name",
				Email:    "user@email.com",
				Password: "  whateva  ",
			},
			expected: handlers.AuthSignupBody{
				Username: "user_name",
				Email:    "user@email.com",
				Password: "  whateva  ",
			},
		},
		{
			name: "Trim and lowercase username and email",
			input: handlers.AuthSignupBody{
				Username: "  User_Name  ",
				Email:    "  User@Email.com  ",
				Password: "  whateva  ",
			},
			expected: handlers.AuthSignupBody{
				Username: "user_name",
				Email:    "user@email.com",
				Password: "  whateva  ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.input
			b.Normalize()

			if b.Username != tt.expected.Username {
				t.Errorf("Normalize() Username = %v, want %v", b.Username, tt.expected.Username)
			}
			if b.Email != tt.expected.Email {
				t.Errorf("Normalize() Email = %v, want %v", b.Email, tt.expected.Email)
			}
			if b.Password != tt.expected.Password {
				t.Errorf("Normalize() Password = %v, want %v", b.Password, tt.expected.Password)
			}
		})
	}
}

func TestAuthSignupBody_Validate(t *testing.T) {
	tests := []struct {
		name       string
		input      handlers.AuthSignupBody
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Valid input",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		// username
		{
			name: "Missing username",
			input: handlers.AuthSignupBody{
				Username: "",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username is required",
		},
		{
			name: "Username too short",
			input: handlers.AuthSignupBody{
				Username: "ab",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username must be min 3 characters",
		},
		{
			name: "Username too long",
			input: handlers.AuthSignupBody{
				Username: "thisusernameiswaytoolongtobevalid",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username must be max 32 characters",
		},
		{
			name: "Username contains invalid character",
			input: handlers.AuthSignupBody{
				Username: "invalid_%_user",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username can only contain lowercase characters, numbers, hyphens, and underscores",
		},
		{
			name: "Username starts with hyphen",
			input: handlers.AuthSignupBody{
				Username: "-invalid_user",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username cannot start with hyphen or underscore",
		},
		{
			name: "Username ends with hyphen",
			input: handlers.AuthSignupBody{
				Username: "invalid_user-",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username cannot end with hyphen or underscore",
		},
		{
			name: "Username starts with underscore",
			input: handlers.AuthSignupBody{
				Username: "_invalid_user",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username cannot start with hyphen or underscore",
		},
		{
			name: "Username ends with underscore",
			input: handlers.AuthSignupBody{
				Username: "invalid_user_",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username cannot end with hyphen or underscore",
		},
		{
			name: "Username is reserved",
			input: handlers.AuthSignupBody{
				Username: "admin",
				Email:    "valid@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "username is reserved",
		},
		// email
		{
			name: "Email is missing",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "email is required",
		},
		{
			name: "Email is invalid",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Email:    "invalidemail.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "email format is invalid",
		},
		{
			name: "Email's length exceeds 254 characters",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Email:    strings.Repeat("a", 245) + "@email.com",
				Password: "strongpassword",
			},
			wantErr:    true,
			wantErrMsg: "email must be at most 254 characters",
		},
		// password
		{
			name: "Missing password",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Email:    "valid@email.com",
			},
			wantErr:    true,
			wantErrMsg: "password is required",
		},
		{
			name: "Password too short",
			input: handlers.AuthSignupBody{
				Username: "valid_user",
				Email:    "valid@email.com",
				Password: "password",
			},
			wantErr:    true,
			wantErrMsg: "password must be at least 12 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.input
			gotErr := b.Validate()

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Validate() failed: %v", gotErr)
				}
				if gotErr.Error() != tt.wantErrMsg {
					t.Errorf("Validate() error = %q, want %q", gotErr.Error(), tt.wantErrMsg)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Validate() succeeded unexpectedly")
			}
		})
	}
}
