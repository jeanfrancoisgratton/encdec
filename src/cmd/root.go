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
	Short: "Encode and decode a string or file to-from BASE64",
	//Version: hftx.White("1.30.00 (2026.04.13), Go version : " + strings.TrimPrefix(runtime.Version(), "go")),
	Version: "1.32.00 (2026.04.14), Go version : " + strings.TrimPrefix(runtime.Version(), "go"),
}

var clCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Shows changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changelog()
	},
}

var encodeCmd = &cobra.Command{
	Use:     "encode",
	Aliases: []string{"enc"},
	Example: "encdec enc {[-f sourcefile [destfile]] | sourcestring}",
	Short:   "Encodes a string or a file",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var s string
		var e *ce.CustomError

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
			fmt.Println("File " + args[0] + hftx.Green("encoded successfully"))
		}
	},
}

var decodeCmd = &cobra.Command{
	Use:     "decode",
	Aliases: []string{"dec"},
	Example: "encdec dec {[-f sourcefile [destfile]] | sourcestring}",
	Short:   "Decodes a string or a file",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var s string
		var e *ce.CustomError

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
			fmt.Println("File " + args[0] + hftx.Green("decoded successfully"))
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
	rootCmd.AddCommand(clCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)

	rootCmd.PersistentFlags().BoolVarP(&executor.Quiet, "quiet", "q", false, "Only show the encrypted/decrypted string")
	rootCmd.PersistentFlags().BoolVarP(&executor.PromptForKeys, "prompt", "p", false, "Should we prompt for a secret key")
	rootCmd.PersistentFlags().BoolVarP(&executor.DEBUG, "debug", "", false, "Debug mode: show extra output")
	decodeCmd.PersistentFlags().BoolVarP(&executor.FileOps, "file", "f", false, "Are we dealing with a file or not")
	encodeCmd.PersistentFlags().BoolVarP(&executor.FileOps, "file", "f", false, "Are we dealing with a file or not")
	decodeCmd.PersistentFlags().BoolVarP(&executor.Keep, "keep", "k", false, "Should we keep the original file")
	encodeCmd.PersistentFlags().BoolVarP(&executor.Keep, "keep", "k", false, "Should we keep the original file")
}

func changelog() {
	//fmt.Printf("\x1b[2J")
	fmt.Printf("\x1bc")

	fmt.Print(`
VERSION		DATE			COMMENT
-------		----			-------
1.32.00		2026.05.11		Go version bump (1.26.3), binary packaging overhaul for RHEL and ArchLinux
1.31.00		2026.04.14		GO version bump (1.26.2), fixed inconsistent error handling; decoding a string actually returned a re-encoded one
1.30.00		2025.11.17		GO version bump (1.25.4), major package and builddeps update. Added a forgotten error path
1.21.03		2024.12.19		GO version bump (1.23.4)
1.21.02		2024.08.13		Variables reshuffling
1.21.01		2024.08.12		Inverted quiet-verbose switch
1.20.01		2024.08.12		Better file handling for destination file, added github actions, go version bump
1.10.00		2024.06.25		Added -q switch, moved to github's helperFunctions package
1.02.00		2023.11.06		Fixed argument count error, version numbering scheme change
1.000		2023.08.02		Updated changelogs and some forgotten release numbers in packaging scripts
0.200		2023.07.31		added file encryption/decryption capabilities
0.100		2023.07.09		stub
`)
}
