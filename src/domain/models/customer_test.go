package models

import (
	"goddd/src/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomer_NewCustomer(t *testing.T) {
	// Build our needed testcase data struct
	testCases := []struct {
		test   string
		name   string
		email  string
		expErr error
	}{
		{
			test:   "Empty Name validation",
			name:   "",
			expErr: errors.InvalidFormat.New("customer name should not be empty"),
		},
		{
			test:   "Valid Name",
			name:   "Johnny Jo",
			email:  "a@b.c",
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		// Run Tests
		t.Run(tc.test, func(t *testing.T) {
			// Create a new customer
			customer, err := NewCustomer(tc.name, tc.email)
			assert.Equal(t, err, tc.expErr)
			if tc.expErr == nil {
				assert.Equal(t, customer.Name, tc.name)
				assert.Equal(t, customer.Email, tc.email)
			}
		})
	}
}
