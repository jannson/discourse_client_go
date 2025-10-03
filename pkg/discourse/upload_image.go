package discourse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadImageResponse represents the response from Discourse image upload API.
type UploadImageResponse struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// UploadImage uploads an image to Discourse and returns its URL, width, and height.
func UploadImage(client *Client, filename string) (resp *UploadImageResponse, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	_ = writer.WriteField("type", "composer")
	writer.Close()

	req, err := http.NewRequest("POST", client.host+"/uploads.json", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Api-Key", client.apiKey)
	req.Header.Set("Api-Username", client.username)

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("upload failed: %s", string(body))
	}

	resp = &UploadImageResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	return resp, err
}
