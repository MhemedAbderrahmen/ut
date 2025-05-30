package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ut",
	Short: "UploadThing CLI - Upload and manage files from your terminal",
	Long: `UploadThing CLI (ut) is a powerful command-line interface for UploadThing.
	
Upload files, download them, and manage your UploadThing storage directly
from your terminal with features like progress tracking, custom output paths,
and support for both public and private files.

Visit https://uploadthing.com to get your API key and start using the CLI.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}
