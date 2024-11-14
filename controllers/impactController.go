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

type impactInput struct {
	PortfolioHepytechId uint   `json:"id_portfolio"`
	Impact_title        string `json:"impact_title"`
	Impact_desc         string `json:"impact_desc"`
	// Impact_icon  			string    `gorm:"text" json:"impact_icon"`
	Flag      			uint    `json:"flag"`
}

func GetAllImpact(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var impact []models.Impact

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
			query = query.Where(column+" LIKE ?", "%"+value+"%")
		}
	}

	var totalCount int64
	query.Model(&impact).Where("flag = 1").Count(&totalCount)

	// Calculate the total pages
	totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

	// Calculate the offset for pagination
	offset := (pagination.Page - 1) * pagination.PerPage

	// Apply pagination and sorting
	err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Preload("Portfolio.Category").Find(&impact).Error
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

	// Mengubah setiap entri dalam data impact untuk menambahkan URL lengkap
	for i := range impact {
		impact[i].Impact_icon = serverAddress + impact[i].Impact_icon
	}

	response := map[string]interface{}{
		"data":         impact,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	checkAndLogActivity(c, "Get all impact", response)
}

func CreateImpact(c *gin.Context) {
	// Validate input
	var input impactInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Impact_icon")
	if errfilename != nil {
		SendError(c, "Upload error", errfilename.Error())
		return
	}

	impact_icon := constanta.DIR_FILE + filename

	// Create
	impact := models.Impact{PortfolioHepytechId: input.PortfolioHepytechId, Impact_title: input.Impact_title, Impact_desc: input.Impact_desc, Impact_icon: impact_icon, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&impact)

	err := db.Preload("Portfolio.Category").First(&impact, impact.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	SendResponse(c, impact, "success")
	activityMessage := "Create impact: " + input.Impact_title
	activitylog(c, activityMessage)
}

func GetImpactByid(c *gin.Context) { // Get model if exist
	var impact models.Impact

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).Preload("Portfolio.Category").First(&impact).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	serverAddress := "http://" + c.Request.Host

	// Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
	impact.Impact_icon = serverAddress + impact.Impact_icon

	checkAndLogActivity(c, "Get impact by id "+c.Param("id"), impact)
}

func UpdateImpact(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var impact models.Impact
	if err := db.Where("id = ?", c.Param("id")).First(&impact).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input impactInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Impact_icon")
	if errfilename != nil {
		SendError(c, "File upload error", errfilename.Error())
		return
	}

	// Cek apakah ada file yang diunggah
	if filename != "" {
		// Hapus gambar lama jika ada
		if impact.Impact_icon != "" {
			oldImage := "." + impact.Impact_icon
			utils.RemoveFile(oldImage)
		}
	}

	// Jika ada file yang diunggah, set nama file yang baru
	if filename != "" {
		impact.Impact_icon = constanta.DIR_FILE + filename
	}

	oldName := impact.Impact_title

	var updatedInput models.Impact
	updatedInput.PortfolioHepytechId = input.PortfolioHepytechId
	updatedInput.Impact_title = input.Impact_title
	updatedInput.Impact_desc = input.Impact_desc
	updatedInput.Impact_icon = impact.Impact_icon
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&impact).Updates(updatedInput)

	err := db.Preload("Portfolio.Category").First(&impact, impact.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	SendResponse(c, impact, "success")
	activityMessage := "Update impact:'" + oldName + "' to '" + input.Impact_title + "'"
	activitylog(c, activityMessage)
}

func DeleteImpact(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var impact models.Impact
	if err := db.Where("id = ?", c.Param("id")).First(&impact).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&impact).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	err := db.Preload("Portfolio.Category").First(&impact, impact.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	// Return success response
	SendResponse(c, impact, "success")
	activityMessage := "Delete impact: " + impact.Impact_title
	activitylog(c, activityMessage)
}