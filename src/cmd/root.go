// encdec
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"encdec/executor"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "encdec",
	Short: "Encode and decode a string, file, or directory",
	Long: `encdec encrypts and decrypts a string (the default), a file (-f, --file), or every regular file
below a directory (-d, --directory), with AES-256-CFB, or with AES-256-GCM when -g, --gcm is passed.`,
	Example: "  encdec encode -s PASSPHRASE \"some secret text\"\n" +
		"  encdec decode -f -s PASSPHRASE secrets.txt\n" +
		"  encdec version",
}

// versionCmd prints the version and build date injected at link time (see types.go)
var versionCmd = &cobra.Command{
	Use:     "version",
	Example: "  encdec version",
	Short:   "Shows the software version",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("encdec v" + buildVersion + " (" + buildDate + "), Go version = v" + strings.TrimPrefix(runtime.Version(), "go") + " (" + runtime.GOARCH + ")"))
	},
}

// operationArgs validates encode/decode arguments: directory mode takes an optional rootdir,
// string and file modes need at least the string or source file
func operationArgs(cmd *cobra.Command, args []string) error {
	if executor.DirectoryOps {
		return cobra.MaximumNArgs(1)(cmd, args)
	}
	return cobra.MinimumNArgs(1)(cmd, args)
}

// gcmHelp closes the Long help of both encode and decode
const gcmHelp = `With -g, --gcm, data is encrypted with AES-256-GCM (Argon2id-derived key)
instead of the default AES-256-CFB. The two formats are NOT interchangeable:
data encoded with -g must be decoded with -g, and data encoded without it must
be decoded without it.`

// encodeCmd encrypts a string, a file (-f) or a directory tree (-d)
var encodeCmd = &cobra.Command{
	Use:     "encode {<STRING> | -f <SOURCEFILE> [<DESTFILE>] | -d [<ROOTDIR>]}",
	Aliases: []string{"enc"},
	Example: "  encdec encode [-s PASSPHRASE] [-g] [-q] <STRING>\n" +
		"  encdec encode [-s PASSPHRASE] [-g] [-k] [-F] -f <SOURCEFILE> [<DESTFILE>]\n" +
		"  encdec encode [-s PASSPHRASE] [-g] [-k] [-F] -d [<ROOTDIR>]\n" +
		"  encdec enc -s PASSPHRASE \"some secret text\"   # prints the Base64 ciphertext\n" +
		"  encdec enc -f secrets.txt                     # encrypts secrets.txt in place\n" +
		"  encdec enc -f -k secrets.txt                  # keeps secrets.txt, writes secrets.txt.enc\n" +
		"  encdec enc -f secrets.txt out.enc             # writes out.enc, leaves secrets.txt alone\n" +
		"  encdec enc -g -d /path/to/dir                 # encrypts every file below /path/to/dir with AES-256-GCM",
	Short: "Encodes a string, file, or directory",
	Long: `A string is encrypted to stdout. A file is encrypted in place, unless a destination file is given;
-k keeps the original next to the encrypted copy (<SOURCEFILE>.enc) on in-place runs.
A directory (the current one by default) has every regular file below it encrypted in place, as with -f.
An existing destination is never overwritten without -F.

` + gcmHelp,
	Args: operationArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var s string
		var e *ce.CustomError

		if executor.DirectoryOps {
			rootdir := ""
			if len(args) > 0 {
				rootdir = args[0]
			}
			if err := executor.EncodeDirectory(rootdir); err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}
			if !executor.Quiet {
				if rootdir == "" {
					rootdir = "."
				}
				fmt.Println("Directory " + rootdir + hftx.Green(" encoded successfully"))
			}
			return
		}

		if !executor.FileOps {
			// encode a string
			s, e = executor.Encode(args[0])
			if e != nil {
				fmt.Println(e.Error())
				os.Exit(1)
			}
			if executor.Quiet {
				fmt.Println(s)
			} else {
				fmt.Println("Encoded string is: " + hftx.Green(s))
			}
			os.Exit(0)
		}

		// encode a file

		dst := ""
		if len(args) > 1 {
			dst = args[1]
		}
		if err := executor.EncodeFile(args[0], dst); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		if !executor.Quiet {
			fmt.Println("File " + args[0] + hftx.Green(" encoded successfully"))
		}
	},
}

// decodeCmd decrypts a string, a file (-f) or a directory tree (-d) produced by encodeCmd
var decodeCmd = &cobra.Command{
	Use:     "decode {<STRING> | -f <SOURCEFILE> [<DESTFILE>] | -d [<ROOTDIR>]}",
	Aliases: []string{"dec"},
	Example: "  encdec decode [-s PASSPHRASE] [-g] [-q] <STRING>\n" +
		"  encdec decode [-s PASSPHRASE] [-g] [-k] [-F] -f <SOURCEFILE> [<DESTFILE>]\n" +
		"  encdec decode [-s PASSPHRASE] [-g] [-k] [-F] -d [<ROOTDIR>]\n" +
		"  encdec dec -s PASSPHRASE AAECAwQF...          # the Base64 string printed by encode\n" +
		"  encdec dec -f secrets.txt                     # decrypts secrets.txt in place\n" +
		"  encdec dec -f -k secrets.txt                  # keeps secrets.txt, writes secrets.txt.dec\n" +
		"  encdec dec -f out.enc secrets.txt             # writes secrets.txt, leaves out.enc alone\n" +
		"  encdec dec -g -d /path/to/dir                 # decrypts every file below /path/to/dir with AES-256-GCM",
	Short: "Decodes a string, file, or directory",
	Long: `A string is decrypted to stdout. A file is decrypted in place, unless a destination file is given;
-k keeps the original next to the decrypted copy (<SOURCEFILE>.dec) on in-place runs.
A directory (the current one by default) has every regular file below it decrypted in place, as with -f.
An existing destination is never overwritten without -F.
Decoding needs the same passphrase, and the same -g setting, that encoding used.

` + gcmHelp,
	Args: operationArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var s string
		var e *ce.CustomError

		if executor.DirectoryOps {
			rootdir := ""
			if len(args) > 0 {
				rootdir = args[0]
			}
			if err := executor.DecodeDirectory(rootdir); err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}
			if !executor.Quiet {
				if rootdir == "" {
					rootdir = "."
				}
				fmt.Println("Directory " + rootdir + hftx.Green(" decoded successfully"))
			}
			return
		}

		if !executor.FileOps {
			// decode a string
			s, e = executor.Decode(args[0])
			if e != nil {
				fmt.Println(e.Error())
				os.Exit(1)
			}
			if executor.Quiet {
				fmt.Println(s)
			} else {
				fmt.Println("Decoded string is: " + hftx.Green(s))
			}
			os.Exit(0)
		}

		// decode a file

		dst := ""
		if len(args) > 1 {
			dst = args[1]
		} else {
			dst = ""
		}
		if err := executor.DecodeFile(args[0], dst); err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		if !executor.Quiet {
			fmt.Println("File " + args[0] + hftx.Green(" decoded successfully"))
		}
	},
}

// Execute runs the command line; it is called once, by main()
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(encodeCmd, decodeCmd, versionCmd)

	rootCmd.PersistentFlags().StringVarP(&executor.Passphrase, "secret", "s", "", "Passphrase to encrypt/decrypt the data; optional, of any length (defaults to the empty passphrase)")
	rootCmd.PersistentFlags().BoolVarP(&executor.Quiet, "quiet", "q", false, "Only show the encrypted/decrypted string")
	rootCmd.PersistentFlags().BoolVarP(&executor.DEBUG, "debug", "", false, "Debug mode: show extra output")
	rootCmd.PersistentFlags().BoolVarP(&executor.UseGCM, "gcm", "g", false, "Use AES-256-GCM (stronger, slower) instead of the default AES-256-CFB; the two formats are not interchangeable")
	decodeCmd.PersistentFlags().BoolVarP(&executor.FileOps, "file", "f", false, "Are we dealing with a file or not")
	encodeCmd.PersistentFlags().BoolVarP(&executor.FileOps, "file", "f", false, "Are we dealing with a file or not")
	decodeCmd.PersistentFlags().BoolVarP(&executor.DirectoryOps, "directory", "d", false, "Recursively decode all regular files below rootdir (defaults to the current directory)")
	encodeCmd.PersistentFlags().BoolVarP(&executor.DirectoryOps, "directory", "d", false, "Recursively encode all regular files below rootdir (defaults to the current directory)")
	decodeCmd.MarkFlagsMutuallyExclusive("file", "directory")
	encodeCmd.MarkFlagsMutuallyExclusive("file", "directory")
	decodeCmd.PersistentFlags().BoolVarP(&executor.Keep, "keep", "k", false, "Keep the original file (in-place runs only, ie. when no destfile is given)")
	encodeCmd.PersistentFlags().BoolVarP(&executor.Keep, "keep", "k", false, "Keep the original file (in-place runs only, ie. when no destfile is given)")
	decodeCmd.PersistentFlags().BoolVarP(&executor.Force, "force", "F", false, "Overwrite the destination file if it already exists")
	encodeCmd.PersistentFlags().BoolVarP(&executor.Force, "force", "F", false, "Overwrite the destination file if it already exists")
}
