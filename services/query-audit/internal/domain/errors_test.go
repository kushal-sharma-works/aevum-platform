package domain

import (
	"errors"
	"testing"
)

func TestDomainErrorWithoutCause(t *testing.T) {
	err := NewDomainError(ErrInvalidQuery, "bad query")

	if err.Code != ErrInvalidQuery {
		t.Fatalf("expected code %s, got %s", ErrInvalidQuery, err.Code)
	}

	if err.Error() != "[INVALID_QUERY] bad query" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
}

func TestDomainErrorWithCause(t *testing.T) {
	cause := errors.New("backend unavailable")
	err := NewDomainErrorWithCause(ErrSearchFailed, "query failed", cause)

	if err.Cause != cause {
		t.Fatal("expected original cause to be preserved")
	}

	if err.Error() != "[SEARCH_FAILED] query failed: backend unavailable" {
		t.Fatalf("unexpected error string: %s", err.Error())
	}
}

func TestDomainErrorCodes(t *testing.T) {
	testCases := []struct {
		code    ErrorCode
		message string
		expect  string
	}{
		{ErrInvalidQuery, "invalid syntax", "[INVALID_QUERY] invalid syntax"},
		{ErrSearchFailed, "timeout", "[SEARCH_FAILED] timeout"},
	}

	for _, tc := range testCases {
		err := NewDomainError(tc.code, tc.message)
		if err.Error() != tc.expect {
			t.Errorf("code %s: got %s, want %s", tc.code, err.Error(), tc.expect)
		}
	}
}
