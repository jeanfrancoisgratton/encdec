// encdec : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/cmd/root.go

package cmd

import (
	"encdec/executor"
	"fmt"
	"os"
	"runtime"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v4/terminalfx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "encdec",
	Short:   "Encode and decode a string or file to-from AES-256",
	Version: hftx.White(fmt.Sprintf("1.21.03-0-%s (2024.12.19)", runtime.GOARCH)),
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
	Aliases: []string{"enc", "encrypt"},
	Example: "encdec enc {[-f sourcefile [destfile]] | sourcestring}",
	Short:   "Encrypts a string or a file",
	Run: func(cmd *cobra.Command, args []string) {
		var cerr *ce.CustomError
		result := ""
		if len(args) < 1 {
			fmt.Println("You need to specify the source (string or filename)")
			os.Exit(1)
		}
		if !executor.FileOps {
			// decode a string
			result, cerr = executor.Encode(args[0])
			if cerr != nil {
				fmt.Println(cerr.Error())
				os.Exit(3)
			}
			if !executor.Quiet {
				result = fmt.Sprintf("Encoded string : %s\n", hftx.Green(result))
			}
			fmt.Println(result)
			os.Exit(0)
		}
		// decode a file

		dst := ""
		if len(args) > 1 {
			dst = args[1]
		} else {
			dst = ""
		}
		if err := executor.EncodeFile(args[0], dst); err != nil {
			fmt.Printf("Error encoding %s : %v", args[0], err)
			os.Exit(3)
		}
	},
}

var decodeCmd = &cobra.Command{
	Use:     "decode",
	Aliases: []string{"dec", "decrypt"},
	Example: "encdec dec {[-f sourcefile [destfile]] | sourcestring}",
	Short:   "Decrypts a string or a file",
	Run: func(cmd *cobra.Command, args []string) {
		var cerr *ce.CustomError
		result := ""
		if len(args) < 1 {
			fmt.Println("You need to specify the source (string or filename)")
			os.Exit(1)
		}
		if !executor.FileOps {
			// decode a string
			result, cerr = executor.Decode(args[0])
			if cerr != nil {
				fmt.Println(cerr.Error())
				os.Exit(2)
			}
			if !executor.Quiet {
				result = fmt.Sprintf("Decoded string : %s\n", hftx.Green(result))
			}
			fmt.Println(result)
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
			fmt.Printf("Error decoding %s : %v", args[0], err)
			os.Exit(2)
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

	rootCmd.PersistentFlags().BoolVarP(&executor.Quiet, "quiet", "q", true, "Only show the encrypted/decrypted string")
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
1.30.00		2025.11.15		GO version bump (1.25.4), major package and builddeps update. Added a forgotten error path
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
