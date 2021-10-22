package config

import (
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	expect := "fakeEnvVar"
	os.Setenv("PAGARME", expect)
	actual, err := GetEnvVar("PAGARME")
	if err != nil {
		t.Errorf("TestGetEnvVar failed: %s", err.Error())
	}
	if actual != expect {
		t.Errorf("Error getting variable, got: %s, want: %s", actual, expect)
	}
}

func TestGetEnvVarWithError(t *testing.T) {
	expect := errMissingEnvVar
	_, err := GetEnvVar("")
	if err != expect {
		t.Errorf("Error getting variable, got: %s, want: %s", err, expect)
	}
}
