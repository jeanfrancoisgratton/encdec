package executor

import (
	"testing"
)

// An arbitrary passphrase used across the tests; helperFunctions derives the
// AES-256 key from its SHA256 sum, so its length is irrelevant.
const testPassphrase = "0123456789abcdef0123456789abcdef"

// withPassphrase temporarily sets the global Passphrase for the duration of a test,
// restoring whatever was there before once the test finishes.
func withPassphrase(t *testing.T, passphrase string) {
	t.Helper()
	prev := Passphrase
	Passphrase = passphrase
	t.Cleanup(func() { Passphrase = prev })
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)

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
	withPassphrase(t, testPassphrase)

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

// Passphrases of any length, the empty one included, are valid.
func TestEncodeDecodeAcceptsAnyPassphraseLength(t *testing.T) {
	cases := []struct{ name, passphrase string }{
		{"empty", ""},
		{"single byte", "x"},
		{"long", "a rather long passphrase that goes well beyond 32 bytes"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withPassphrase(t, tc.passphrase)

			enc, cerr := Encode("payload")
			if cerr != nil {
				t.Fatalf("Encode with passphrase %q returned error: %s", tc.passphrase, cerr.Error())
			}
			dec, cerr := Decode(enc)
			if cerr != nil {
				t.Fatalf("Decode with passphrase %q returned error: %s", tc.passphrase, cerr.Error())
			}
			if dec != "payload" {
				t.Errorf("roundtrip with passphrase %q: got %q, want %q", tc.passphrase, dec, "payload")
			}
		})
	}
}

// Not passing -s must leave the passphrase empty rather than fall back to some
// built-in key.
func TestDefaultPassphraseIsEmpty(t *testing.T) {
	if Passphrase != "" {
		t.Errorf("default passphrase is %q, want the empty string", Passphrase)
	}
}

// The wrong passphrase must not give the plaintext back.
func TestDecodeWithWrongPassphrase(t *testing.T) {
	withPassphrase(t, testPassphrase)
	enc, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("Encode returned error: %s", cerr.Error())
	}

	withPassphrase(t, "some other passphrase")
	dec, cerr := Decode(enc)
	if cerr == nil && dec == "payload" {
		t.Error("decoding with the wrong passphrase returned the original plaintext")
	}
}

// A blob shorter than the AES block size cannot contain an IV.
func TestDecodeRejectsShortCiphertext(t *testing.T) {
	withPassphrase(t, testPassphrase)

	// "AAAA" decodes to 3 bytes, which is < aes.BlockSize (16).
	if _, cerr := Decode("AAAA"); cerr == nil {
		t.Error("expected Decode to reject ciphertext shorter than the block size")
	}
}
