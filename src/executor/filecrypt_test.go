package executor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// quiet silences the tool's decorative output for the duration of a test.
func quiet(t *testing.T) {
	t.Helper()
	prev := Quiet
	Quiet = true
	t.Cleanup(func() { Quiet = prev })
}

// keep sets the global Keep flag for the duration of a test.
func keep(t *testing.T, v bool) {
	t.Helper()
	prev := Keep
	Keep = v
	t.Cleanup(func() { Keep = prev })
}

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("could not write temp file %s: %v", path, err)
	}
	return path
}

// With -k (keep), both source and destination survive and a round-trip
// reproduces the original bytes.
func TestEncodeDecodeFileRoundtrip(t *testing.T) {
	withKey(t, testKey)
	quiet(t)
	keep(t, true)

	plaintext := bytes.Repeat([]byte("encdec file mode test payload\n"), 5000) // ~145 KB, spans several chunks
	src := writeTempFile(t, "plain.txt", plaintext)
	encPath := src + ".enc"
	decPath := src + ".dec"

	if cerr := EncodeFile(src, encPath); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}

	// The ciphertext must exist and differ from the plaintext.
	encData, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatalf("could not read encrypted file: %v", err)
	}
	if bytes.Equal(encData, plaintext) {
		t.Fatal("encrypted file is identical to the plaintext")
	}

	if cerr := DecodeFile(encPath, decPath); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}

	decData, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatalf("could not read decrypted file: %v", err)
	}
	if !bytes.Equal(decData, plaintext) {
		t.Errorf("decrypted file does not match original (got %d bytes, want %d)", len(decData), len(plaintext))
	}
}

// Without -k, the source is replaced in place: after encode+decode the
// original path holds the original content and no temp files remain.
func TestEncodeDecodeFileInPlace(t *testing.T) {
	withKey(t, testKey)
	quiet(t)
	keep(t, false)

	plaintext := []byte("in-place round trip\n")
	src := writeTempFile(t, "secret.txt", plaintext)

	if cerr := EncodeFile(src, ""); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}
	if _, err := os.Stat(src + ".enc"); !os.IsNotExist(err) {
		t.Errorf("expected temporary .enc file to be gone, stat err = %v", err)
	}

	if cerr := DecodeFile(src, ""); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}
	if _, err := os.Stat(src + ".dec"); !os.IsNotExist(err) {
		t.Errorf("expected temporary .dec file to be gone, stat err = %v", err)
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("could not read file after in-place round trip: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("in-place round trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestEncodeFileMissingSource(t *testing.T) {
	withKey(t, testKey)
	quiet(t)
	keep(t, true)

	missing := filepath.Join(t.TempDir(), "does-not-exist")
	if cerr := EncodeFile(missing, missing+".enc"); cerr == nil {
		t.Error("expected EncodeFile to fail on a missing source file, got nil error")
	}
}
