package collection

import (
	"strings"
	"testing"
)

func Test_ValidateName(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "Valid name",
			input:      "My Collection",
			wantResult: "My Collection",
			wantErr:    false,
		},
		{
			name:       "Two character name is valid",
			input:      "AI",
			wantResult: "AI",
			wantErr:    false,
		},
		{
			name:       "Name is trimmed",
			input:      "  My Collection  ",
			wantResult: "My Collection",
			wantErr:    false,
		},
		{
			name:       "Name is too short",
			input:      "a",
			wantErr:    true,
			wantErrMsg: "name must be between 2 and 255 characters",
		},
		{
			name:       "Name is too long",
			input:      strings.Repeat("a", 256),
			wantErr:    true,
			wantErrMsg: "name must be between 2 and 255 characters",
		},
		{
			name:       "Name with null byte is rejected",
			input:      "My \x00 Collection",
			wantErr:    true,
			wantErrMsg: "name must not contain control characters",
		},
		{
			name:       "Name with tab is rejected",
			input:      "My \t Collection",
			wantErr:    true,
			wantErrMsg: "name must not contain control characters",
		},
		{
			name:       "Multi-byte characters are counted as characters not bytes",
			input:      strings.Repeat("ą", 128),
			wantResult: strings.Repeat("ą", 128),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotErr := ValidateName(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateName() failed: %v", gotErr)
				}
				if gotErr.Error() != tt.wantErrMsg {
					t.Errorf("ValidateName() error message = %v, want %v", gotErr.Error(), tt.wantErrMsg)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateName() succeeded unexpectedly")
			}
			if tt.wantResult != "" && gotResult != tt.wantResult {
				t.Errorf("ValidateName() result = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_ValidateDescription(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:       "Valid description",
			input:      "A helpful collection description.",
			wantResult: "A helpful collection description.",
			wantErr:    false,
		},
		{
			name:    "Empty description",
			input:   "",
			wantErr: false,
		},
		{
			name:       "Description is trimmed",
			input:      "  A helpful collection description.  ",
			wantResult: "A helpful collection description.",
			wantErr:    false,
		},
		{
			name:       "Description is too long",
			input:      strings.Repeat("a", 513),
			wantErr:    true,
			wantErrMsg: "description must be less than 512 characters",
		},
		{
			name:       "Description with null byte is rejected",
			input:      "A description with \x00 null byte",
			wantErr:    true,
			wantErrMsg: "description must not contain control characters",
		},
		{
			name:       "Description with newline is rejected",
			input:      "A description with \n newline",
			wantErr:    true,
			wantErrMsg: "description must not contain control characters",
		},
		{
			name:       "Multi-byte characters are counted as characters not bytes",
			input:      strings.Repeat("ą", 257),
			wantResult: strings.Repeat("ą", 257),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotErr := ValidateDescription(tt.input)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateDescription() failed: %v", gotErr)
				}
				if gotErr.Error() != tt.wantErrMsg {
					t.Errorf("ValidateDescription() error message = %v, want %v", gotErr.Error(), tt.wantErrMsg)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateDescription() succeeded unexpectedly")
			}
			if tt.wantResult != "" && gotResult != tt.wantResult {
				t.Errorf("ValidateDescription() result = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
