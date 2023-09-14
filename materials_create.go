package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
)

const maxUploadSize = 20 * 1024 * 1024

func materialsCreate(w http.ResponseWriter, r *http.Request) {
	SetCorsPolicies(w, r)

	team, err := GetTeam(PathParam(r, 0))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := &Response{
			Message: "team not found",
			Status:  http.StatusNotFound,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
			return
		}
		w.Write(responseJSON)
		return
	}

	_, err = GetMaterials(team.ID)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		response := &Response{
			Message: "materials already exists",
			Status:  http.StatusConflict,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
			return
		}
		w.Write(responseJSON)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	err = r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		fmt.Println(err)

		w.WriteHeader(http.StatusBadRequest)
		response := &Response{
			Message: "failed to parse multipart form",
			Status:  http.StatusBadRequest,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	file, _, err := r.FormFile("pdf")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := &Response{
			Message: "failed to get pdf",
			Status:  http.StatusBadRequest,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}
	defer file.Close()

	webpFiles, err := convertPdfToWebp(file, team.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to convert pdf to webp",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	urls, err := uploadFileToAzure(webpFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := &Response{
			Message: "failed to upload file to azure",
			Status:  http.StatusInternalServerError,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			failedToMarshalResponse(w)
		}
		w.Write(responseJSON)
		return
	}

	for _, url := range urls {
		u, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := &Response{
				Message: "failed to generate team id",
				Status:  http.StatusInternalServerError,
			}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				failedToMarshalResponse(w)
			}
			w.Write(responseJSON)
			return
		}
		uu := u.String()

		material := &Material{
			ID:     uu,
			TeamID: team.ID,
			Url:    url,
		}

		result := DB.Create(material)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := &Response{
				Message: "failed to create material",
				Status:  http.StatusInternalServerError,
			}
			responseJSON, err := json.Marshal(response)
			if err != nil {
				failedToMarshalResponse(w)
			}
			w.Write(responseJSON)
			return
		}
	}

	os.RemoveAll(filepath.Dir(webpFiles[0]))

	w.Write([]byte("{\"message\": \"ok\", \"status\": 200}"))
}

func uploadFileToAzure(filePaths []string) ([]string, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	if accountName == "" || accountKey == "" || containerName == "" {
		return nil, fmt.Errorf("failed to get azure storage env")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	var urls []string
	for _, filePath := range filePaths {
		blobURL := containerURL.NewBlockBlobURL(filepath.Base(filePath))

		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, err = azblob.UploadFileToBlockBlob(context.TODO(), file, blobURL, azblob.UploadToBlockBlobOptions{})
		if err != nil {
			return nil, err
		}
		urls = append(urls, blobURL.String())
	}

	return urls, nil
}

func convertPdfToWebp(pdfFile multipart.File, teamId string) ([]string, error) {
	// ここで作成する一時ディレクトリの削除は関数の呼び出し元が行うこと
	tmpDir, err := os.MkdirTemp("", "pdf_conversion")
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp("", "input_*.pdf")
	if err != nil {
		return nil, err
	}
	io.Copy(tmpFile, pdfFile)
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()

	var outBuf bytes.Buffer
	cmd := exec.Command("pdftoppm", "-png", tmpFilePath, filepath.Join(tmpDir, teamId))
	cmd.Stdout = &outBuf
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	outputFiles, err := filepath.Glob(filepath.Join(tmpDir, teamId+"*.png"))
	if err != nil {
		return nil, err
	}

	webpFiles := make([]string, 0, len(outputFiles))
	for _, file := range outputFiles {
		imgData, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		img, err := png.Decode(bytes.NewReader(imgData))
		if err != nil {
			return nil, err
		}

		webpPath := file + ".webp"
		webpFile, err := os.Create(webpPath)
		if err != nil {
			return nil, err
		}

		err = webp.Encode(webpFile, img, &webp.Options{Lossless: true})
		if err != nil {
			return nil, err
		}
		webpFile.Close()

		webpFiles = append(webpFiles, webpPath)
	}

	return webpFiles, nil
}
