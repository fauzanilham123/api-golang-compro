package controllers

import (
	"api-golang-compro/constanta"
	"api-golang-compro/utils"
	"errors"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HandleUploadFile(c *gin.Context, form string) (string, error) {
    fileHeader, _ := c.FormFile(form)
    var filename string

    // Cek apakah ada file yang diunggah
    if fileHeader != nil {
        fileExtention := []string{".png", ".jpg", ".jpeg", ".svg"}
        isFileValidated := utils.FileValidationByExtension(fileHeader, fileExtention)
        if !isFileValidated {
            return "", errors.New("file not allowed")
        }

        // Batas ukuran file dalam byte (2MB)
        maxFileSize := int64(2 * 1024 * 1024) // 2MB

        // Periksa ukuran file
        if fileHeader.Size > maxFileSize {
            return "", errors.New("File size too large (max 2MB)")
        }

        extensionFile := filepath.Ext(fileHeader.Filename)
        filename = utils.RandomFileName(extensionFile)

        isSaved := utils.SaveFile(c, fileHeader, filename)
        if !isSaved {
            return "", errors.New("Internal server error, can't save file")
        }
    }

    // Jika tidak ada file yang diunggah, filename akan tetap kosong

    return filename, nil
}




func HandleRemoveFile(c *gin.Context){
	filename := c.Param("name")

	err := utils.RemoveFile(constanta.DIR_FILE + filename)
	if err != nil{
		SendError(c,"file name is required.",err.Error())
		return
	}
	SendResponse(c,filename,"file deleted")
}