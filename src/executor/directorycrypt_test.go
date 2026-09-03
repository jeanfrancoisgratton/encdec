// encdec
// Written by Jean-François Gratton <jean-francois.gratton@aylo.com>
// Original timestamp : 2026.09.03 14:54:51
// Original filename : src/executor/directorycrypt_test.go

package executor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeDirectoryTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("could not create directory for %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("could not write %s: %v", path, err)
	}
}

func TestEncodeDecodeDirectoryRoundtrip(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)

	root := t.TempDir()
	files := map[string][]byte{
		"root.txt":              []byte("file in the root\n"),
		"nested/child.txt":      bytes.Repeat([]byte("nested payload\n"), 100),
		"nested/deeper/.hidden": []byte("hidden file\n"),
		"empty":                 {},
	}
	for name, contents := range files {
		writeDirectoryTestFile(t, filepath.Join(root, name), contents)
	}
	if err := os.Mkdir(filepath.Join(root, "empty-directory"), 0o755); err != nil {
		t.Fatalf("could not create empty directory: %v", err)
	}

	outside := writeTempFile(t, "outside.txt", []byte("must not be encoded\n"))
	link := filepath.Join(root, "outside-link")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatalf("could not create symbolic link: %v", err)
	}

	if cerr := EncodeDirectory(root); cerr != nil {
		t.Fatalf("EncodeDirectory returned error: %s", cerr.Error())
	}
	for name, plaintext := range files {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("could not read encoded file %s: %v", name, err)
		}
		if bytes.Equal(got, plaintext) {
			t.Errorf("encoded file %s is identical to its plaintext", name)
		}
	}
	if got, err := os.ReadFile(outside); err != nil || string(got) != "must not be encoded\n" {
		t.Errorf("symbolic-link target was modified: %q (err %v)", got, err)
	}

	if cerr := DecodeDirectory(root); cerr != nil {
		t.Fatalf("DecodeDirectory returned error: %s", cerr.Error())
	}
	for name, plaintext := range files {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("could not read decoded file %s: %v", name, err)
		}
		if !bytes.Equal(got, plaintext) {
			t.Errorf("decoded file %s does not match its original contents", name)
		}
	}
}

func TestEncodeDirectoryDefaultsToCurrentDirectory(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)

	previousDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get current directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDirectory); err != nil {
			t.Errorf("could not restore current directory: %v", err)
		}
	})

	root := t.TempDir()
	path := filepath.Join(root, "payload.txt")
	plaintext := []byte("current-directory payload\n")
	writeDirectoryTestFile(t, path, plaintext)
	if err := os.Chdir(root); err != nil {
		t.Fatalf("could not change current directory: %v", err)
	}

	if cerr := EncodeDirectory(""); cerr != nil {
		t.Fatalf("EncodeDirectory returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(path); err != nil || bytes.Equal(got, plaintext) {
		t.Errorf("file in current directory was not encoded: %q (err %v)", got, err)
	}
	if cerr := DecodeDirectory(""); cerr != nil {
		t.Fatalf("DecodeDirectory returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, plaintext) {
		t.Errorf("current-directory round trip failed: %q (err %v)", got, err)
	}
}

func TestEncodeDirectoryWithKeepPreservesSources(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, true)
	force(t, false)

	root := t.TempDir()
	path := filepath.Join(root, "nested", "payload.txt")
	plaintext := []byte("preserve me\n")
	writeDirectoryTestFile(t, path, plaintext)

	if cerr := EncodeDirectory(root); cerr != nil {
		t.Fatalf("EncodeDirectory returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, plaintext) {
		t.Errorf("source was not preserved: %q (err %v)", got, err)
	}
	if got, err := os.ReadFile(path + ".enc"); err != nil || bytes.Equal(got, plaintext) {
		t.Errorf("encoded copy was not created: %q (err %v)", got, err)
	}
}

func TestDirectoryOperationRejectsFileRoot(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)

	path := writeTempFile(t, "not-a-directory", []byte("payload\n"))
	if cerr := EncodeDirectory(path); cerr == nil {
		t.Fatal("expected EncodeDirectory to reject a file root, got nil error")
	}
}

func TestEncodeDirectoryAcceptsSymlinkedRoot(t *testing.T) {
	withPassphrase(t, testPassphrase)
	quiet(t)
	keep(t, false)
	force(t, false)

	parent := t.TempDir()
	root := filepath.Join(parent, "real-root")
	path := filepath.Join(root, "payload.txt")
	plaintext := []byte("symlinked-root payload\n")
	writeDirectoryTestFile(t, path, plaintext)

	linkedRoot := filepath.Join(parent, "linked-root")
	if err := os.Symlink(root, linkedRoot); err != nil {
		t.Fatalf("could not create root symbolic link: %v", err)
	}
	if cerr := EncodeDirectory(linkedRoot); cerr != nil {
		t.Fatalf("EncodeDirectory returned error: %s", cerr.Error())
	}
	if got, err := os.ReadFile(path); err != nil || bytes.Equal(got, plaintext) {
		t.Errorf("file below symbolic-link root was not encoded: %q (err %v)", got, err)
	}
}
