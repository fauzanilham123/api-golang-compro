package controllers

import (
	"api-golang-compro/constanta"
	"api-golang-compro/models"
	"api-golang-compro/utils"
	"math"

	// "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type serviceInput struct {
	// Icon      string    `gorm:"text" json:"icon"`
	Title     string    `json:"title"`
	Flag      uint      `json:"flag"`
	Description string  `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllService(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var service []models.Service

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
    query.Model(&service).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&service).Error
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

    // Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
    for i := range service {
        service[i].Icon = serverAddress + service[i].Icon
    }

    response := map[string]interface{}{
        "data":        service,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all service",response)
}

func CreateService(c *gin.Context) {
	// Validate input
	var input serviceInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Icon")
	if errfilename != nil {
    SendError(c, "Upload error", errfilename.Error())
    return
	}

    icon := constanta.DIR_FILE + filename

	// Create
	service := models.Service{Icon: icon, Title: input.Title, Description: input.Description, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&service)

	SendResponse(c, service, "success")
	activityMessage := "Create service: " +input.Title
    activitylog(c,activityMessage)
}

func GetServiceByid(c *gin.Context) { // Get model if exist
	var service models.Service

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&service).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	serverAddress := "http://" + c.Request.Host

    // Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
    service.Icon = serverAddress + service.Icon

	checkAndLogActivity(c,"Get service by id "+ c.Param("id"),service)
}

func UpdateService(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var service models.Service
	if err := db.Where("id = ?", c.Param("id")).First(&service).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input serviceInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Icon")
	if errfilename != nil {
    SendError(c, "File upload error", errfilename.Error())
    return
	}

    // Cek apakah ada file yang diunggah
	if filename != "" {
    // Hapus gambar lama jika ada
    if service.Icon != "" {
        oldImage := "." + service.Icon
        utils.RemoveFile(oldImage)
    	}
	}

    // Jika ada file yang diunggah, set nama file yang baru
    if filename != "" {
        service.Icon = constanta.DIR_FILE + filename
    }

	oldName := service.Title

	var updatedInput models.Service
	updatedInput.Icon = service.Icon
	updatedInput.Title = input.Title
	updatedInput.Description = input.Description
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&service).Updates(updatedInput)

	SendResponse(c, service, "success")
    activityMessage := "Update service:'" + oldName + "' to '" + input.Title + "'"
    activitylog(c,activityMessage)
}

func DeleteService(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var service models.Service
	if err := db.Where("id = ?", c.Param("id")).First(&service).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&service).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, service, "success")
	activityMessage := "Delete service: "+ service.Title
    activitylog(c,activityMessage)
}