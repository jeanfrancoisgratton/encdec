// encdec
// Tests the command tree as a whole: help consistency and the version output

package cmd

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// commands returns every command of the tree, rootCmd included
func commands() []*cobra.Command {
	var all []*cobra.Command
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		all = append(all, c)
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(rootCmd)
	return all
}

// Every command carries a Short description and an Example section, indented like cobra's own
func TestEveryCommandHasHelp(t *testing.T) {
	for _, c := range commands() {
		if c.Hidden || c.Name() == "help" || c.Name() == "completion" || (c.Parent() != nil && c.Parent().Name() == "completion") {
			continue
		}
		if c.Short == "" {
			t.Errorf("%s has no Short description", c.CommandPath())
		}
		if c.Example == "" {
			t.Errorf("%s has no Example section", c.CommandPath())
		}
		for _, line := range strings.Split(c.Example, "\n") {
			if !strings.HasPrefix(line, "  encdec ") {
				t.Errorf("%s: example line %q does not start with \"  encdec \"", c.CommandPath(), line)
			}
		}
	}
}

// hasFlag tells if c, or one of its subcommands, has the flag of that name or (one-letter) shorthand
func hasFlag(c *cobra.Command, name string) bool {
	for _, fs := range []*pflag.FlagSet{c.Flags(), c.PersistentFlags(), c.InheritedFlags()} {
		if fs.Lookup(name) != nil || (len(name) == 1 && fs.ShorthandLookup(name) != nil) {
			return true
		}
	}
	for _, sub := range c.Commands() {
		if hasFlag(sub, name) {
			return true
		}
	}
	return false
}

// Every flag the help text or the examples mention must exist on that command (or, for a
// parent command, on one of its subcommands)
func TestDocumentedFlagsExist(t *testing.T) {
	flagRef := regexp.MustCompile(`(?m)(?:^|[\s\[|])--?([a-zA-Z][a-z-]*)\b`)
	for _, c := range commands() {
		for _, m := range flagRef.FindAllStringSubmatch(c.Long+"\n"+c.Example, -1) {
			if !hasFlag(c, m[1]) {
				t.Errorf("%s mentions -%s, which it does not have", c.CommandPath(), m[1])
			}
		}
	}
}

func TestVersion(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("could not create pipe: %v", err)
	}
	stdout := os.Stdout
	os.Stdout = w
	rootCmd.SetArgs([]string{"version"})
	execErr := rootCmd.Execute()
	os.Stdout = stdout
	w.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		t.Fatalf("could not read version output: %v", err)
	}
	if execErr != nil {
		t.Fatalf("version returned error: %v", execErr)
	}
	if !strings.Contains(out.String(), "encdec v"+buildVersion+" ("+buildDate+"), Go version = v") {
		t.Fatalf("output = %q", out.String())
	}
}
