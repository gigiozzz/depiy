package utility

import (
	"testing"
)

func TestTruncateString(t *testing.T) {
	data := "0123456789"

	expectedName := "01234567"
	actualName := TruncateString(data, 8)
	if actualName != expectedName {
		t.Fatalf("Invalid truncation for long %q. Expected %q, got %q", data, expectedName, actualName)
	}

	shortData := "0123"
	expectedName = shortData
	actualName = TruncateString(shortData, 8)
	if actualName != expectedName {
		t.Fatalf("Invalid truncation for short %q. Expected %q, got %q", shortData, expectedName, actualName)
	}

}

func TestApplayerGenerateSha256(t *testing.T) {
	data := "docker.io/gigiozzz/bundle-test-op"
	expectedName := "a4e2c0a3b13df95de94508ab6f4ef32647c7e259f37dbdbcc52064d1b5db39e7"
	actualName := GenerateSha256(data)
	if actualName != expectedName {
		t.Fatalf("Invalid generation for %q. Expected %q, got %q", data, expectedName, actualName)
	}
}
