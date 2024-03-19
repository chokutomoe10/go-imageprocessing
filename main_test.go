package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
)

func setupRequest(t *testing.T, url string, body *bytes.Buffer, contentType string) *http.Request {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)
	return req
}

func createMultipartFormFile(t *testing.T, writer *multipart.Writer, i int, filePath string) {
	fileByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file %s: %s", filePath, err)
		return
	}

	ext := filepath.Ext(filePath)
	part, err := writer.CreateFormFile("Images", fmt.Sprintf("sample%d%s", i+1, ext))
	if err != nil {
		t.Fatalf("failed to create form file: %s", err)
		return
	}

	if _, err := io.Copy(part, bytes.NewReader(fileByte)); err != nil {
		t.Fatalf("failed to copy file contents: %s", err)
		return
	}
}

func makeRequestAndAssert(t *testing.T, url string, handler http.HandlerFunc, expectedStatusCode int, expectedBody string, body *bytes.Buffer, contentType string) {
	rr := httptest.NewRecorder()
	req := setupRequest(t, url, body, contentType)
	handler(rr, req)

	if rr.Code != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, expectedStatusCode)
	}

	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestConvertHandler(t *testing.T) {
	filePaths := []string{
		"uploads/PNG_transparency_demonstration_1.png",
		"uploads/Ducati_side_shadow.png",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, filePath := range filePaths {
		createMultipartFormFile(t, writer, i, filePath)
	}
	writer.Close()

	makeRequestAndAssert(t, "/convert", Convert, http.StatusOK, "Successfully Converted Files", body, writer.FormDataContentType())
}

func TestResizeHandler(t *testing.T) {
	filePaths := []string{
		"uploads/PNG_transparency_demonstration_1.png",
		"uploads/Ducati_side_shadow.png",
		"uploads/view.jpeg",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, filePath := range filePaths {
		createMultipartFormFile(t, writer, i, filePath)
	}

	sizeX := 400
	sizeY := 300

	sizeXStr := strconv.Itoa(sizeX)
	sizeYStr := strconv.Itoa(sizeY)

	writer.WriteField("SizeX", sizeXStr)
	writer.WriteField("SizeY", sizeYStr)

	writer.Close()

	makeRequestAndAssert(t, "/resize", Resize, http.StatusOK, "Successfully Resized Files", body, writer.FormDataContentType())
}

func TestCompressHandler(t *testing.T) {
	filePaths := []string{
		"uploads/PNG_transparency_demonstration_1.png",
		"uploads/Ducati_side_shadow.png",
		"uploads/view.jpeg",
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for i, filePath := range filePaths {
		createMultipartFormFile(t, writer, i, filePath)
	}

	sizeX := 400
	sizeY := 300
	sizeJPEG := 50

	sizeXStr := strconv.Itoa(sizeX)
	sizeYStr := strconv.Itoa(sizeY)
	sizeJPEGStr := strconv.Itoa(sizeJPEG)

	writer.WriteField("SizeX", sizeXStr)
	writer.WriteField("SizeY", sizeYStr)
	writer.WriteField("SizeJPEG", sizeJPEGStr)

	writer.Close()

	makeRequestAndAssert(t, "/compress", Compress, http.StatusOK, "Successfully Compressed Files", body, writer.FormDataContentType())
}
