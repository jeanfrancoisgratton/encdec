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

// force sets the global Force flag for the duration of a test.
func force(t *testing.T, v bool) {
	t.Helper()
	prev := Force
	Force = v
	t.Cleanup(func() { Force = prev })
}

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("could not write temp file %s: %v", path, err)
	}
	return path
}

// With explicit destinations, both source and destination survive and a
// round-trip reproduces the original bytes.
func TestEncodeDecodeFileRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)
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
	withPassphrase(t, testPassphrase)
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

// Naming a destination must never move the result onto the source: the source
// stays put, untouched, even without -k.
func TestExplicitDestinationLeavesSourceAlone(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)

	plaintext := []byte("the source must survive\n")
	src := writeTempFile(t, "in.txt", plaintext)
	dst := filepath.Join(filepath.Dir(src), "out.enc")

	if cerr := EncodeFile(src, dst); cerr != nil {
		t.Fatalf("EncodeFile returned error: %s", cerr.Error())
	}

	got, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("source file is gone: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("source file was modified: got %q, want %q", got, plaintext)
	}

	ciphertext, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("destination file was not written: %v", err)
	}
	if bytes.Equal(ciphertext, plaintext) {
		t.Error("destination holds the plaintext")
	}

	// And the reverse trip, likewise, must not disturb its own source.
	back := filepath.Join(filepath.Dir(src), "out.dec")
	if cerr := DecodeFile(dst, back); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}
	if _, err := os.Stat(dst); err != nil {
		t.Errorf("decode removed its source file: %v", err)
	}
	if got, err := os.ReadFile(back); err != nil || !bytes.Equal(got, plaintext) {
		t.Errorf("round trip through explicit destinations failed: %q (err %v)", got, err)
	}
}

// An existing destination is not collateral damage: refuse it without -F.
func TestRefusesToOverwriteExistingDestination(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)

	plaintext := []byte("payload\n")
	src := writeTempFile(t, "in.txt", plaintext)
	precious := []byte("DO NOT CLOBBER ME\n")
	dst := writeTempFile(t, "out.enc", precious)

	if cerr := EncodeFile(src, dst); cerr == nil {
		t.Fatal("expected EncodeFile to refuse an existing destination, got nil error")
	}

	if got, _ := os.ReadFile(dst); !bytes.Equal(got, precious) {
		t.Errorf("destination was clobbered: got %q, want %q", got, precious)
	}
	if got, _ := os.ReadFile(src); !bytes.Equal(got, plaintext) {
		t.Errorf("source was disturbed by a refused run: got %q, want %q", got, plaintext)
	}
}

// The scratch file of an in-place run gets the same protection.
func TestRefusesToOverwriteExistingScratchFile(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)

	src := writeTempFile(t, "in.txt", []byte("payload\n"))
	precious := []byte("leftover from an aborted run\n")
	if err := os.WriteFile(src+".enc", precious, 0o600); err != nil {
		t.Fatalf("could not write scratch file: %v", err)
	}

	if cerr := EncodeFile(src, ""); cerr == nil {
		t.Fatal("expected EncodeFile to refuse an existing scratch file, got nil error")
	}
	if got, _ := os.ReadFile(src + ".enc"); !bytes.Equal(got, precious) {
		t.Errorf("scratch file was clobbered: got %q, want %q", got, precious)
	}
}

func TestForceOverwritesExistingDestination(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, true)

	plaintext := []byte("payload\n")
	src := writeTempFile(t, "in.txt", plaintext)
	dst := writeTempFile(t, "out.enc", []byte("stale content\n"))

	if cerr := EncodeFile(src, dst); cerr != nil {
		t.Fatalf("EncodeFile with -F returned error: %s", cerr.Error())
	}

	back := filepath.Join(filepath.Dir(src), "out.dec")
	if cerr := DecodeFile(dst, back); cerr != nil {
		t.Fatalf("DecodeFile returned error: %s", cerr.Error())
	}
	if got, _ := os.ReadFile(back); !bytes.Equal(got, plaintext) {
		t.Errorf("forced overwrite did not produce valid ciphertext: got %q, want %q", got, plaintext)
	}
}

func TestEncodeFileMissingSource(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, true)

	missing := filepath.Join(t.TempDir(), "does-not-exist")
	if cerr := EncodeFile(missing, missing+".enc"); cerr == nil {
		t.Error("expected EncodeFile to fail on a missing source file, got nil error")
	}
}
