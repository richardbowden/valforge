package main

import (
	"testing"
	valgen "tests/internal/valgen"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name      string
		user      User
		wantErr   bool
		errFields []string
	}{
		{
			name: "valid user",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "john@example.com",
				Color: "blue",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			user: User{
				Name:  "",
				Age:   25,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "john@example.com",
				Color: "blue",
			},
			wantErr:   true,
			errFields: []string{"name"},
		},
		{
			name: "age below 18",
			user: User{
				Name:  "Jane Doe",
				Age:   17,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "jane@example.com",
				Color: "red",
			},
			wantErr:   true,
			errFields: []string{"age"},
		},
		{
			name: "age exactly 18 (boundary test)",
			user: User{
				Name:  "Jane Doe",
				Age:   18,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "jane@example.com",
				Color: "red",
			},
			wantErr: false,
		},
		{
			name: "password too short",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "pass",
				Pwd2:  "pass",
				Email: "john@example.com",
				Color: "blue",
			},
			wantErr:   true,
			errFields: []string{"pwd1"},
		},
		{
			name: "password exactly 6 characters (boundary test)",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "pass12",
				Pwd2:  "pass12",
				Email: "john@example.com",
				Color: "blue",
			},
			wantErr: false,
		},
		{
			name: "passwords do not match",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "password123",
				Pwd2:  "password456",
				Email: "john@example.com",
				Color: "blue",
			},
			wantErr:   true,
			errFields: []string{"pwd2"},
		},
		{
			name: "invalid email",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "invalid-email",
				Color: "blue",
			},
			wantErr:   true,
			errFields: []string{"email"},
		},
		{
			name: "empty email",
			user: User{
				Name:  "John Doe",
				Age:   25,
				Pwd1:  "password123",
				Pwd2:  "password123",
				Email: "",
				Color: "blue",
			},
			wantErr:   true,
			errFields: []string{"email"},
		},
		{
			name: "multiple validation errors",
			user: User{
				Name:  "",
				Age:   15,
				Pwd1:  "123",
				Pwd2:  "456",
				Email: "bad-email",
				Color: "",
			},
			wantErr:   true,
			errFields: []string{"name", "age", "pwd1", "pwd2", "email", "color"},
		},
		{
			name: "all fields empty/invalid",
			user: User{
				Name:  "",
				Age:   0,
				Pwd1:  "",
				Pwd2:  "",
				Email: "",
				Color: "",
			},
			wantErr:   true,
			errFields: []string{"name", "age", "pwd1", "email", "color"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				verr, ok := err.(*valgen.ValidationError)
				if !ok {
					t.Errorf("expected *valgen.ValidationError, got %T", err)
					return
				}

				// Check that expected fields have errors
				for _, field := range tt.errFields {
					found := false
					for _, fieldErr := range verr.Errors {
						if fieldErr.Field == field {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected error for field %q, but none found", field)
					}
				}

				// Verify we got the expected number of errors
				if len(verr.Errors) != len(tt.errFields) {
					t.Errorf("expected %d errors, got %d", len(tt.errFields), len(verr.Errors))
				}
			}
		})
	}
}

func TestUser_Validate_PasswordCompare(t *testing.T) {
	user1 := User{
		Name:  "Test User",
		Age:   25,
		Pwd1:  "password123",
		Pwd2:  "password124", // last char different
		Email: "test@example.com",
		Color: "blue",
	}

	err1 := user1.Validate()
	if err1 == nil {
		t.Error("expected error for mismatched passwords")
	}

	user2 := User{
		Name:  "Test User",
		Age:   25,
		Pwd1:  "password123",
		Pwd2:  "xassword123", // first char different
		Email: "test@example.com",
		Color: "blue",
	}

	err2 := user2.Validate()
	if err2 == nil {
		t.Error("expected error for mismatched passwords")
	}

	if (err1 == nil) != (err2 == nil) {
		t.Error("password comparison should behave consistently")
	}
}

func TestUser_Validate_EdgeCases(t *testing.T) {
	t.Run("very long password", func(t *testing.T) {
		longPwd := string(make([]byte, 10000))
		for i := range longPwd {
			longPwd = longPwd[:i] + "a" + longPwd[i+1:]
		}

		user := User{
			Name:  "Test",
			Age:   25,
			Pwd1:  longPwd,
			Pwd2:  longPwd,
			Email: "test@example.com",
			Color: "blue",
		}

		err := user.Validate()
		if err != nil {
			t.Errorf("should accept very long password: %v", err)
		}
	})

	t.Run("negative age", func(t *testing.T) {
		user := User{
			Name:  "Test",
			Age:   -5,
			Pwd1:  "password123",
			Pwd2:  "password123",
			Email: "test@example.com",
			Color: "blue",
		}

		err := user.Validate()
		if err == nil {
			t.Error("expected error for negative age")
		}
	})
}

func TestUser_Validate_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		user          User
		expectedField string
		expectedMsg   string
	}{
		{
			name:          "empty name message",
			user:          User{Age: 20, Pwd1: "password", Pwd2: "password", Email: "test@test.com"},
			expectedField: "name",
			expectedMsg:   "name is required",
		},
		{
			name:          "age too low message",
			user:          User{Name: "Test", Age: 10, Pwd1: "password", Pwd2: "password", Email: "test@test.com"},
			expectedField: "age",
			expectedMsg:   "age must be greater than or equal to 18",
		},
		{
			name:          "password too short message",
			user:          User{Name: "Test", Age: 20, Pwd1: "12345", Pwd2: "12345", Email: "test@test.com"},
			expectedField: "pwd1",
			expectedMsg:   "pwd1 must be at least 6 characters",
		},
		{
			name:          "password mismatch message",
			user:          User{Name: "Test", Age: 20, Pwd1: "password1", Pwd2: "password2", Email: "test@test.com"},
			expectedField: "pwd2",
			expectedMsg:   "pwd2 must match Pwd1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			verr := err.(*valgen.ValidationError)
			found := false
			for _, fieldErr := range verr.Errors {
				if fieldErr.Field == tt.expectedField && fieldErr.Message == tt.expectedMsg {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("expected error message %q for field %q", tt.expectedMsg, tt.expectedField)
			}
		})
	}
}
