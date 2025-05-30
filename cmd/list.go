package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"ut/config"

	"github.com/spf13/cobra"
)

type FilesResponse struct {
	HasMore bool       `json:"hasMore"`
	Files   []FileInfo `json:"files"`
}

type FileInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	FileKey    string `json:"key"`
	UploadedAt int64  `json:"uploadedAt"`
}

var (
	verbose bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all uploaded files",
	Long:  `List all files uploaded to your UploadThing storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := listFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing files: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed file information")
}

func listFiles() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	emptyBody := make(map[string]interface{})
	jsonBody, err := json.Marshal(emptyBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	apiURL := "https://api.uploadthing.com/v6/listFiles"
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Uploadthing-Api-Key", cfg.SecretKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed: status %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var filesResp FilesResponse
	err = json.Unmarshal(body, &filesResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(filesResp.Files) == 0 {
		fmt.Println("No files found.")
		return nil
	}

	fmt.Printf("Found %d files:\n\n", len(filesResp.Files))

	if verbose {
		for _, file := range filesResp.Files {
			uploadedTime := time.Unix(file.UploadedAt, 0)
			fmt.Printf("📄 %s\n", file.Name)
			fmt.Printf("   File Key: %s\n", file.FileKey)
			fmt.Printf("   Size: %s\n", formatFileSize(file.Size))
			fmt.Printf("   Uploaded: %s\n", uploadedTime.Format("2006-01-02 15:04:05"))
			fmt.Printf("   ID: %s\n\n", file.ID)
		}
	} else {
		for _, file := range filesResp.Files {
			fmt.Printf("📄 %-30s %s\n", file.Name, file.FileKey)
		}
	}

	if filesResp.HasMore {
		fmt.Println("\n... more files available (pagination not implemented)")
	}

	return nil
}
