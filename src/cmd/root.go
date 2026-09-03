// encdec : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/root.go

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
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("encdec v1.6.0 (2026.09.03), Go version = v" + strings.TrimPrefix(runtime.Version(), "go")))
	},
}

func operationArgs(cmd *cobra.Command, args []string) error {
	if executor.DirectoryOps {
		return cobra.MaximumNArgs(1)(cmd, args)
	}
	return cobra.MinimumNArgs(1)(cmd, args)
}

var encodeCmd = &cobra.Command{
	Use:     "encode",
	Aliases: []string{"enc"},
	Example: "encdec enc {[-f sourcefile [destfile]] | [-d [rootdir]] | sourcestring}",
	Short:   "Encodes a string, file, or directory",
	Args:    operationArgs,
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

var decodeCmd = &cobra.Command{
	Use:     "decode",
	Aliases: []string{"dec"},
	Example: "encdec dec {[-f sourcefile [destfile]] | [-d [rootdir]] | sourcestring}",
	Short:   "Decodes a string, file, or directory",
	Args:    operationArgs,
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
