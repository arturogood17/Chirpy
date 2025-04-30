package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tokenSecret := "This is a test"
	userID := uuid.New()
	token, _ := MakeJWT(userID, tokenSecret)
	tests := []struct {
		input   string
		tokenS  string
		wantErr bool
	}{
		{
			input:   token,
			tokenS:  tokenSecret,
			wantErr: false,
		},
		{
			input:   token,
			tokenS:  "123456",
			wantErr: true,
		},
	}

	for _, test := range tests {
		_, err := ValidateJWT(test.input, test.tokenS)
		if (err != nil) != test.wantErr {
			t.Errorf("Error: Want %v, got: %v", test.wantErr, err)
			return
		}
	}
}
