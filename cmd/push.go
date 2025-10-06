package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ut/config"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "push <filepath> [filepath2] [filepath3]...",
	Short: "Push one or more files to UploadThing",
	Long:  `Push one or more files to UploadThing using your secret API key configured.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, filePath := range args {
			fmt.Printf("[%d/%d] Uploading %s...\n", i+1, len(args), filepath.Base(filePath))
			err := uploadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error uploading file %s: %v\n", filePath, err)
				os.Exit(1)
			}
			fmt.Printf("[%d/%d] ✓ %s uploaded successfully!\n", i+1, len(args), filepath.Base(filePath))
		}
		if len(args) > 1 {
			fmt.Printf("All %d files uploaded successfully!\n", len(args))
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

type UploadFilesRequest struct {
	Files              []FileMetadata `json:"files"`
	ACL                string         `json:"acl,omitempty"`
	ContentDisposition string         `json:"contentDisposition,omitempty"`
}

type FileMetadata struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	CustomID string `json:"customId,omitempty"`
}

type UploadFilesResponse struct {
	Data []PresignedUpload `json:"data"`
}

type PresignedUpload struct {
	URL                string            `json:"url"`
	Fields             map[string]string `json:"fields"`
	Key                string            `json:"key"`
	FileName           string            `json:"fileName"`
	FileType           string            `json:"fileType"`
	FileUrl            string            `json:"fileUrl"`
	ContentDisposition string            `json:"contentDisposition"`
}

func uploadFile(filePath string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileName := filepath.Base(file.Name())
	fileSize := fileInfo.Size()

	contentType := "application/octet-stream"
	if ext := filepath.Ext(fileName); ext != "" {
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".pdf":
			contentType = "application/pdf"
		case ".txt":
			contentType = "text/plain"
		case ".json":
			contentType = "application/json"
		case ".xml":
			contentType = "application/xml"
		case ".csv":
			contentType = "text/csv"
		}
	}

	uploadReq := UploadFilesRequest{
		Files: []FileMetadata{
			{
				Name: fileName,
				Size: fileSize,
				Type: contentType,
			},
		},
		ACL:                "public-read",
		ContentDisposition: "inline",
	}

	reqBody, err := json.Marshal(uploadReq)
	if err != nil {
		return fmt.Errorf("failed to marshal upload request: %w", err)
	}

	fmt.Printf("Requesting presigned URL...\n")

	apiURL := "https://api.uploadthing.com/v6/uploadFiles"
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Uploadthing-Api-Key", cfg.SecretKey)

	client := &http.Client{Timeout: 60 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read upload response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get presigned URL: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	var uploadResp UploadFilesResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return fmt.Errorf("failed to unmarshal upload response: %w", err)
	}

	if len(uploadResp.Data) == 0 {
		return fmt.Errorf("no presigned upload data received from UploadThing")
	}

	presignedUpload := uploadResp.Data[0]
	fmt.Printf("Got presigned URL: %s\n", presignedUpload.URL)

	file.Seek(0, 0)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range presignedUpload.Fields {
		err := writer.WriteField(key, value)
		if err != nil {
			return fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	uploadFileReq, err := http.NewRequest(http.MethodPost, presignedUpload.URL, body)
	if err != nil {
		return fmt.Errorf("failed to create file upload request: %w", err)
	}

	uploadFileReq.Header.Set("Content-Type", writer.FormDataContentType())

	fmt.Printf("Uploading file to storage...\n")

	uploadFileResp, err := client.Do(uploadFileReq)
	if err != nil {
		return fmt.Errorf("file upload request failed: %w", err)
	}
	defer uploadFileResp.Body.Close()

	uploadFileRespBody, err := io.ReadAll(uploadFileResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read file upload response: %w", err)
	}

	if uploadFileResp.StatusCode < 200 || uploadFileResp.StatusCode >= 300 {
		return fmt.Errorf("file upload failed: status %d, response: %s", uploadFileResp.StatusCode, string(uploadFileRespBody))
	}

	fmt.Printf("Upload successful!\n")
	fmt.Printf("File key: %s\n", presignedUpload.Key)
	fmt.Printf("File URL: %s\n", presignedUpload.FileUrl)

	return nil
}
