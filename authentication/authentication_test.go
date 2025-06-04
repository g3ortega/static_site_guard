package authentication

import (
	"os"
	"testing"
)

func TestEligibleUser(t *testing.T) {
	os.Setenv("ELIGIBLE_USERNAMES", "alice,bob")
	if !eligibleUser("alice") {
		t.Fatalf("alice should be eligible")
	}
	if eligibleUser("charlie") {
		t.Fatalf("charlie should not be eligible")
	}
}
