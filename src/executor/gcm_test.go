package executor

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

// gcm sets the global UseGCM flag for the duration of a test.
func gcm(t *testing.T, v bool) {
	t.Helper()
	prev := UseGCM
	UseGCM = v
	t.Cleanup(func() { UseGCM = prev })
}

// Not passing -g must keep the legacy Encode*/Decode* behaviour.
func TestDefaultIsNotGCM(t *testing.T) {
	if UseGCM {
		t.Error("UseGCM defaults to true, want false")
	}
}

func TestGCMEncodeDecodeRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)
	gcm(t, true)

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

// The random salt and nonce must make two encryptions of the same plaintext differ.
func TestGCMEncodeIsRandomized(t *testing.T) {
	withPassphrase(t, testPassphrase)
	gcm(t, true)

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

func TestGCMAcceptsEmptyPassphrase(t *testing.T) {
	withPassphrase(t, "")
	gcm(t, true)

	enc, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("Encode with the empty passphrase returned error: %s", cerr.Error())
	}
	dec, cerr := Decode(enc)
	if cerr != nil {
		t.Fatalf("Decode with the empty passphrase returned error: %s", cerr.Error())
	}
	if dec != "payload" {
		t.Errorf("roundtrip with the empty passphrase: got %q, want %q", dec, "payload")
	}
}

// Unlike CFB, GCM is authenticated: the wrong passphrase must be an error,
// not garbage.
func TestGCMDecodeWithWrongPassphraseFails(t *testing.T) {
	withPassphrase(t, testPassphrase)
	gcm(t, true)
	enc, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("Encode returned error: %s", cerr.Error())
	}

	withPassphrase(t, "some other passphrase")
	if dec, cerr := Decode(enc); cerr == nil {
		t.Errorf("expected Decode with the wrong passphrase to fail, got %q", dec)
	}
}

func TestGCMDecodeRejectsTamperedCiphertext(t *testing.T) {
	withPassphrase(t, testPassphrase)
	gcm(t, true)
	enc, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("Encode returned error: %s", cerr.Error())
	}

	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		t.Fatalf("Encode output is not valid Base64: %v", err)
	}
	raw[len(raw)-1] ^= 0x01
	tampered := base64.StdEncoding.EncodeToString(raw)

	if dec, cerr := Decode(tampered); cerr == nil {
		t.Errorf("expected Decode to reject tampered ciphertext, got %q", dec)
	}
}

func TestGCMDecodeRejectsMalformedInput(t *testing.T) {
	withPassphrase(t, testPassphrase)
	gcm(t, true)

	cases := []struct{ name, in string }{
		{"not base64", "this is not base64!"},
		{"shorter than the salt", "AAAA"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, cerr := Decode(tc.in); cerr == nil {
				t.Errorf("expected Decode(%q) to fail", tc.in)
			}
		})
	}
}

// The legacy and GCM formats are not interchangeable, in either direction.
func TestGCMAndLegacyFormatsAreIncompatible(t *testing.T) {
	withPassphrase(t, testPassphrase)

	gcm(t, false)
	legacy, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("legacy Encode returned error: %s", cerr.Error())
	}

	gcm(t, true)
	if dec, cerr := Decode(legacy); cerr == nil {
		t.Errorf("GCM Decode accepted legacy ciphertext, got %q", dec)
	}
	modern, cerr := Encode("payload")
	if cerr != nil {
		t.Fatalf("GCM Encode returned error: %s", cerr.Error())
	}

	gcm(t, false)
	if dec, cerr := Decode(modern); cerr == nil && dec == "payload" {
		t.Error("legacy Decode recovered the plaintext from GCM ciphertext")
	}
}

func TestGCMEncodeDecodeFileRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, true)
	gcm(t, true)

	plaintext := bytes.Repeat([]byte("encdec GCM file mode test payload\n"), 5000)
	src := writeTempFile(t, "plain.txt", plaintext)
	encPath := src + ".enc"
	decPath := src + ".dec"

	if cerr := EncodeFile(src, encPath); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}
	encrypted, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatalf("could not read encrypted file: %v", err)
	}
	if bytes.Equal(encrypted, plaintext) {
		t.Fatal("encrypted file is identical to its plaintext")
	}

	if cerr := DecodeFile(encPath, decPath); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}
	decrypted, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatalf("could not read decrypted file: %v", err)
	}
	if !bytes.Equal(decrypted, plaintext) {
		t.Error("decrypted file does not match the original")
	}
}

func TestGCMEncodeDecodeFileInPlace(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)
	gcm(t, true)

	plaintext := []byte("in-place GCM payload\n")
	src := writeTempFile(t, "plain.txt", plaintext)

	if cerr := EncodeFile(src, ""); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(src); err != nil || bytes.Equal(got, plaintext) {
		t.Fatalf("file was not encrypted in place: %q (err %v)", got, err)
	}

	if cerr := DecodeFile(src, ""); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(src); err != nil || !bytes.Equal(got, plaintext) {
		t.Errorf("in-place GCM round trip failed: %q (err %v)", got, err)
	}
}

func TestGCMDecodeFileWithWrongPassphraseFails(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, true)
	gcm(t, true)

	src := writeTempFile(t, "plain.txt", []byte("payload\n"))
	encPath := src + ".enc"
	if cerr := EncodeFile(src, encPath); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}

	withPassphrase(t, "some other passphrase")
	decPath := src + ".dec"
	if cerr := DecodeFile(encPath, decPath); cerr == nil {
		t.Fatal("expected DecodeFile with the wrong passphrase to fail")
	}
	if _, err := os.Stat(decPath); !os.IsNotExist(err) {
		t.Errorf("a failed decryption left %s behind (stat err %v)", decPath, err)
	}
}

func TestGCMEncodeDecodeDirectoryRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)
	gcm(t, true)

	root := t.TempDir()
	files := map[string][]byte{
		"root.txt":         []byte("file in the root\n"),
		"nested/child.txt": bytes.Repeat([]byte("nested payload\n"), 100),
	}
	for name, contents := range files {
		writeDirectoryTestFile(t, filepath.Join(root, name), contents)
	}

	if cerr := EncodeDirectory(root); cerr != nil {
		t.Fatalf("EncodeDirectory returned error: %s", cerr.Error())
	}
	for name, plaintext := range files {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("could not read encrypted file %s: %v", name, err)
		}
		if bytes.Equal(got, plaintext) {
			t.Errorf("encrypted file %s is identical to its plaintext", name)
		}
	}

	if cerr := DecodeDirectory(root); cerr != nil {
		t.Fatalf("DecodeDirectory returned error: %s", cerr.Error())
	}
	for name, plaintext := range files {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("could not read decrypted file %s: %v", name, err)
		}
		if !bytes.Equal(got, plaintext) {
			t.Errorf("decrypted file %s does not match its original contents", name)
		}
	}
}
