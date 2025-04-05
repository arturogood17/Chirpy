package auth

import (
	"testing"
)

func Test_HashingPasswords(t *testing.T) {
	password_1 := "Pepito1234"
	hashedp1, _ := HashPassword(password_1)
	password_2 := "Pepito2569"
	hashedp2, _ := HashPassword(password_2)
	cases := []struct {
		password string
		hash     string
		WantErr  bool
	}{
		{
			password: password_1,
			hash:     hashedp1,
			WantErr:  false,
		},
		{
			password: password_2,
			hash:     hashedp2,
			WantErr:  false,
		},
		{
			password: password_1,
			hash:     hashedp2,
			WantErr:  true,
		},
	}
	for _, test := range cases {
		t.Run(test.password, func(t *testing.T) {
			err := CheckPasswordHash(test.hash, test.password)
			if (err != nil) != test.WantErr {
				t.Errorf("Test failed. Check password hashing failed: %v", err)
			}
		})
	}
}
