package utils

import (
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var charset = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz123456789")

func RandomString(n int) string {
	rand.Seed(time.Now().UnixMilli())
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func FileValidation(fileHeader *multipart.FileHeader, fileType []string) bool {
	contentType := fileHeader.Header.Get("Content-Type")
	log.Println("content-type", contentType)
	result := false

	for _, typeFile := range fileType {
		if contentType == typeFile {
			result = true
			break
		}
	}

	return result
}

func FileValidationByExtension(fileHeader *multipart.FileHeader, fileExtension []string) bool {
	extension := filepath.Ext(fileHeader.Filename)
	log.Println("extension", extension)
	result := false

	for _, typeFile := range fileExtension {
		if extension == typeFile {
			result = true
			break
		}
	}

	return result
}

func RandomFileName(extensionFile string,) string {
	uniqueID := uuid.New()
    currentTime := time.Now().UTC().Format("20061206")
    filename := fmt.Sprintf("file-%s-%s%s", currentTime, uniqueID, extensionFile)
    return filename
}


func SaveFile(c *gin.Context, fileHeader *multipart.FileHeader, filename string) bool {
	errUpload := c.SaveUploadedFile(fileHeader, fmt.Sprintf("./public/file/%s",filename))
	if errUpload != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errUpload.Error(),
		})
		
		return false
	} else {
		return true
	}
}

func RemoveFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		log.Println("Failed to remove file")
		return err
	}
	return nil
}

