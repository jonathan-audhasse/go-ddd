package models

import (
	"testing"
)

func TestCustomer_NewCustomer(t *testing.T) {
	// Build our needed testcase data struct
	testCases := []struct {
		test        string
		name        string
		email       string
		expectedErr error
	}{
		{
			test:        "Empty Name validation",
			name:        "",
			expectedErr: ErrInvalidCustomer,
		}, {
			test:        "Valid Name",
			name:        "Johnny Jo",
			email:       "a@b.c",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new customer
			_, err := NewCustomer(tc.name, "")
			// Check if the error matches the expected error
			if err != tc.expectedErr {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
