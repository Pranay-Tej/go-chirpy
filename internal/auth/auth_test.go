package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password1 := "1234"
	password2 := "abcd"
	hashedPassword1, _ := HashPassword(password1)
	hashedPassword2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "valid",
			password: password1,
			hash:     hashedPassword1,
			wantErr:  false,
		},
		{
			name:     "valid",
			password: password2,
			hash:     hashedPassword2,
			wantErr:  false,
		},
		{
			name:     "invalid",
			password: password1,
			hash:     hashedPassword2,
			wantErr:  true,
		},
		{
			name:     "invalid",
			password: password2,
			hash:     hashedPassword1,
			wantErr:  true,
		},
		{
			name:     "invalid",
			password: password2,
			hash:     "",
			wantErr:  true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckPasswordHash(tc.hash, tc.password)
			if tc.wantErr != (err != nil) {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
