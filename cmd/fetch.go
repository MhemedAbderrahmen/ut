package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ut/config"

	"github.com/spf13/cobra"
)

var (
	outputPath     string
	forceOverwrite bool
	showProgress   bool
	isPrivate      bool
)

var (
	ErrAPIKeyInvalid = errors.New("invalid api key")
)

var downloadCmd = &cobra.Command{
	Use:   "fetch <fileKey>",
	Short: "Download a file from UploadThing",
	Long: `Download a file from UploadThing using a file key.
	
Examples:
  ut fetch abc123-example.jpg                    # Download to current directory
  ut fetch abc123-example.jpg -o myfile.jpg     # Download with custom name
  ut fetch abc123-example.jpg -o ./downloads/   # Download to specific directory
  ut fetch abc123-example.jpg --private         # Download private file (requires API key)
  ut fetch abc123-example.jpg --progress        # Show download progress`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileKey := args[0]
		err := runDownload(fileKey)
		if err != nil {
			if errors.Is(err, config.ErrConfigNotFound) {
				fmt.Fprintln(os.Stderr, `API key is not configured.
Run 'ut config set-secret' before downloading private files.`)
			} else if errors.Is(err, ErrAPIKeyInvalid) {
				fmt.Fprintln(os.Stderr, "Invalid API key. Run 'ut config set-secret' to update it.")
			} else {
				fmt.Fprintf(os.Stderr, "Error downloading file: %v\n", err)
			}
			os.Exit(1)
		}
		fmt.Println("File downloaded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path or directory")
	downloadCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite existing file without prompt")
	downloadCmd.Flags().BoolVarP(&showProgress, "progress", "p", false, "Show download progress")
	downloadCmd.Flags().BoolVar(&isPrivate, "private", false, "Download private file (requires API key)")
}

type FileAccessResponse struct {
	URL string `json:"url"`
}

func runDownload(fileKey string) error {
	if strings.TrimSpace(fileKey) == "" {
		return fmt.Errorf("file key cannot be empty")
	}

	var fileURL string
	var filename string

	if isPrivate {
		signedURL, err := getSignedURL(fileKey)
		if err != nil {
			return fmt.Errorf("failed to get signed URL for private file: %w", err)
		}
		fileURL = signedURL
		filename = extractFilenameFromKey(fileKey)
	} else {
		fileURL = "https://utfs.io/f/" + fileKey
		filename = extractFilenameFromKey(fileKey)
	}

	_, err := url.ParseRequestURI(fileURL)
	if err != nil {
		return fmt.Errorf("invalid URL generated: %w", err)
	}

	outputFilePath, err := determineOutputPath(filename)
	if err != nil {
		return fmt.Errorf("failed to determine output path: %w", err)
	}

	if !forceOverwrite {
		if _, err := os.Stat(outputFilePath); err == nil {
			fmt.Printf("File '%s' already exists. Overwrite? (y/N): ", outputFilePath)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				return fmt.Errorf("download cancelled by user")
			}
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("unable to create output file %s: %w", outputFilePath, err)
	}
	defer outputFile.Close()

	fmt.Printf("Downloading %s...\n", filename)

	if showProgress {
		err = downloadWithProgress(fileURL, outputFile)
	} else {
		err = downloadFile(fileURL, outputFile)
	}

	if err != nil {
		os.Remove(outputFilePath)
		return fmt.Errorf("download failed: %w", err)
	}

	fileInfo, _ := outputFile.Stat()
	fmt.Printf("Download complete: %s (%s)\n", outputFilePath, formatFileSize(fileInfo.Size()))

	return nil
}

func getSignedURL(fileKey string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config (API key required for private files): %w", err)
	}

	apiURL := "https://api.uploadthing.com/v6/requestFileAccess"
	reqBody := fmt.Sprintf(`{"fileKey": "%s"}`, fileKey)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Uploadthing-Api-Key", cfg.SecretKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", ErrAPIKeyInvalid
		}
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed: status %d, response: %s", resp.StatusCode, string(body))
	}

	var accessResp FileAccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&accessResp); err != nil {
		return "", fmt.Errorf("failed to parse API response: %w", err)
	}

	return accessResp.URL, nil
}

func extractFilenameFromKey(fileKey string) string {
	parts := strings.Split(fileKey, "-")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return fileKey
}

func determineOutputPath(filename string) (string, error) {
	if outputPath == "" {
		return filename, nil
	}

	if strings.HasSuffix(outputPath, "/") || strings.HasSuffix(outputPath, "\\") {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
		return filepath.Join(outputPath, filename), nil
	}

	if stat, err := os.Stat(outputPath); err == nil && stat.IsDir() {
		return filepath.Join(outputPath, filename), nil
	}

	parentDir := filepath.Dir(outputPath)
	if parentDir != "." {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory: %w", err)
		}
	}

	return outputPath, nil
}

func downloadFile(fileURL string, outputFile *os.File) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func downloadWithProgress(fileURL string, outputFile *os.File) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var fileSize int64
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
			fileSize = size
		}
	}

	progressWriter := &ProgressWriter{
		Total:      fileSize,
		Downloaded: 0,
		StartTime:  time.Now(),
	}

	_, err = io.Copy(outputFile, io.TeeReader(resp.Body, progressWriter))
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println()
	return nil
}

type ProgressWriter struct {
	Total      int64
	Downloaded int64
	StartTime  time.Time
}

func (pw *ProgressWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	pw.Downloaded += int64(n)

	if pw.Downloaded%102400 == 0 || pw.Downloaded == pw.Total {
		pw.printProgress()
	}

	return n, nil
}

func (pw *ProgressWriter) printProgress() {
	elapsed := time.Since(pw.StartTime)
	speed := float64(pw.Downloaded) / elapsed.Seconds()

	if pw.Total > 0 {
		percentage := float64(pw.Downloaded) / float64(pw.Total) * 100
		fmt.Printf("\rProgress: %.1f%% (%s/%s) - %.2f KB/s",
			percentage,
			formatFileSize(pw.Downloaded),
			formatFileSize(pw.Total),
			speed/1024)
	} else {
		fmt.Printf("\rDownloaded: %s - %.2f KB/s",
			formatFileSize(pw.Downloaded),
			speed/1024)
	}
}
