package controllers

import (
	"api-golang-compro/constanta"
	"api-golang-compro/models"
	"api-golang-compro/utils"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type logoInput struct {
	Name string `json:"name"`
	// Logo		 string    `gorm:"type:text" json:"link_button"`
	Flag      uint      `json:"flag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllLogo(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var logo []models.Logo

    sort := c.DefaultQuery("sort", "asc")
    // Default to ascending if not provided
    sortOrder := "ASC"
    if sort == "desc" {
        sortOrder = "DESC"
    }

    pagination := ExtractPagination(c)
    query := db.Where("flag = 1")

    // Get all query parameters and loop through them
    queryParams := c.Request.URL.Query()
     // Remove 'page' and 'perPage' keys from queryParams
    delete(queryParams, "page")
    delete(queryParams, "perPage")
    delete(queryParams, "sort")
    for column, values := range queryParams {
        value := values[0] // In case there are multiple values, we take the first one

        // Apply filtering condition if the value is not empty
        if value != "" {
            query = query.Where(column + " LIKE ?", "%"+value+"%")
        }
    }

    var totalCount int64
    query.Model(&logo).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&logo).Error
    if err != nil {
        SendError(c, "internal server error", err.Error())
        return
    }

    // Calculate "last_page" based on total pages
    lastPage := totalPages

    // Calculate "nextPage" and "prevPage"
    nextPage := pagination.Page + 1
    if nextPage > totalPages {
        nextPage = 1
    }

    prevPage := pagination.Page - 1
    if prevPage < 1 {
        prevPage = 1
    }

	// Mendapatkan alamat server dari permintaan
    serverAddress := "http://" + c.Request.Host

    // Mengubah setiap entri dalam data logo untuk menambahkan URL lengkap
    for i := range logo {
        logo[i].Logo = serverAddress + logo[i].Logo
    }

    response := map[string]interface{}{
        "data":        logo,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all logo",response)
}

func CreateLogo(c *gin.Context) {
	// Validate input
	var input logoInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Logo")
	if errfilename != nil {
    SendError(c, "Upload error", errfilename.Error())
    return
	}

    Logo := constanta.DIR_FILE + filename

	// Create
	logo := models.Logo{Name: input.Name, Logo: Logo, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&logo)

	SendResponse(c, logo, "success")
	activityMessage := "Create logo: " +input.Name
    activitylog(c,activityMessage)
}

func GetLogoByid(c *gin.Context) { // Get model if exist
	var logo models.Logo

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&logo).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	serverAddress := "http://" + c.Request.Host

    // Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
    logo.Logo = serverAddress + logo.Logo

	checkAndLogActivity(c,"Get logo by id "+ c.Param("id"),logo)
}

func UpdateLogo(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var logo models.Logo
	if err := db.Where("id = ?", c.Param("id")).First(&logo).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input logoInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Logo")
	if errfilename != nil {
    SendError(c, "File upload error", errfilename.Error())
    return
	}

    // Cek apakah ada file yang diunggah
	if filename != "" {
    // Hapus gambar lama jika ada
    if logo.Logo != "" {
        oldImage := "." + logo.Logo
        utils.RemoveFile(oldImage)
    	}
	}

    // Jika ada file yang diunggah, set nama file yang baru
    if filename != "" {
        logo.Logo = constanta.DIR_FILE + filename
    }

	oldName := logo.Name

	var updatedInput models.Logo
	updatedInput.Name = input.Name
	updatedInput.Logo = logo.Logo
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&logo).Updates(updatedInput)

	SendResponse(c, logo, "success")
    activityMessage := "Update logo:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeleteLogo(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var logo models.Logo
	if err := db.Where("id = ?", c.Param("id")).First(&logo).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&logo).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, logo, "success")
	activityMessage := "Delete logo: "+ logo.Name
    activitylog(c,activityMessage)
}