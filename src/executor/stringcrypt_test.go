package executor

import (
	"testing"
)

// A valid 32-byte AES-256 key used across the tests.
const testKey = "0123456789abcdef0123456789abcdef"

// withKey temporarily sets the global SecretKey for the duration of a test,
// restoring whatever was there before once the test finishes.
func withKey(t *testing.T, key string) {
	t.Helper()
	prev := SecretKey
	SecretKey = key
	t.Cleanup(func() { SecretKey = prev })
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	withKey(t, testKey)

	cases := []struct {
		name string
		in   string
	}{
		{"simple", "hello world"},
		{"empty", ""},
		{"unicode", "héllo wörld — �✓"},
		{"long", "The quick brown fox jumps over the lazy dog, repeatedly and at length."},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, cerr := Encode(tc.in)
			if cerr != nil {
				t.Fatalf("Encode(%q) returned error: %s", tc.in, cerr.Error())
			}

			dec, cerr := Decode(enc)
			if cerr != nil {
				t.Fatalf("Decode(%q) returned error: %s", enc, cerr.Error())
			}

			if dec != tc.in {
				t.Errorf("roundtrip mismatch: got %q, want %q", dec, tc.in)
			}
		})
	}
}

// The random IV must make two encryptions of the same plaintext differ.
func TestEncodeUsesRandomIV(t *testing.T) {
	withKey(t, testKey)

	a, cerr := Encode("same input")
	if cerr != nil {
		t.Fatalf("first Encode returned error: %s", cerr.Error())
	}
	b, cerr := Encode("same input")
	if cerr != nil {
		t.Fatalf("second Encode returned error: %s", cerr.Error())
	}

	if a == b {
		t.Errorf("expected two encryptions of the same string to differ, both were %q", a)
	}
}

func TestEncodeRejectsBadKeyLength(t *testing.T) {
	withKey(t, "too-short-key")

	if _, cerr := Encode("data"); cerr == nil {
		t.Error("expected Encode to fail with a non-32-byte key, got nil error")
	}
}

func TestDecodeRejectsBadKeyLength(t *testing.T) {
	withKey(t, "too-short-key")

	if _, cerr := Decode("anything"); cerr == nil {
		t.Error("expected Decode to fail with a non-32-byte key, got nil error")
	}
}

// A blob shorter than the AES block size cannot contain an IV.
func TestDecodeRejectsShortCiphertext(t *testing.T) {
	withKey(t, testKey)

	// "AAAA" decodes to 3 bytes, which is < aes.BlockSize (16).
	if _, cerr := Decode("AAAA"); cerr == nil {
		t.Error("expected Decode to reject ciphertext shorter than the block size")
	}
}
