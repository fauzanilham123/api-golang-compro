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

type portfolioInput struct {
	// Image       string    `json:"image" `
	Title       string    `json:"title"`
	Id_category uint      `json:"id_category"`
	Flag        uint      `json:"flag"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetAllPortfolio(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var portfolio []models.Portfolio

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
    query.Model(&portfolio).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Preload("Category").Find(&portfolio).Error
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
    for i := range portfolio {
        portfolio[i].Image = serverAddress + portfolio[i].Image
    }

    response := map[string]interface{}{
        "data":        portfolio,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all portfolio",response)
}


func CreatePortfolio(c *gin.Context) {
	// Validate input
	var input portfolioInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}
	
	filename, errfilename := HandleUploadFile(c, "Image")
	if errfilename != nil {
    SendError(c, "Upload error", errfilename.Error())
    return
	}

    image := constanta.DIR_FILE + filename

	// Create
	portfolio := models.Portfolio{Image: image, Title: input.Title, Category_homeID: input.Id_category, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&portfolio)

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	SendResponse(c, portfolio, "success")
    activityMessage := "Create portfolio: " +portfolio.Title
    activitylog(c,activityMessage)
}

func GetPortfolioById(c *gin.Context) { // Get model if exist
	var portfolio models.Portfolio

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&portfolio).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

    serverAddress := "http://" + c.Request.Host

    // Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
    portfolio.Image = serverAddress + portfolio.Image
    

	
    checkAndLogActivity(c,"Get portfolio by id "+ c.Param("id"),portfolio)
}

func UpdatePortfolio(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var portfolio models.Portfolio
	if err := db.Where("id = ?", c.Param("id")).First(&portfolio).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input portfolioInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename, errfilename := HandleUploadFile(c, "Image")
	if errfilename != nil {
    SendError(c, "File upload error", errfilename.Error())
    return
	}

    // Cek apakah ada file yang diunggah
	if filename != "" {
    // Hapus gambar lama jika ada
    if portfolio.Image != "" {
        oldImage := "." + portfolio.Image
        utils.RemoveFile(oldImage)
    	}
	}

    // Jika ada file yang diunggah, set nama file yang baru
    if filename != "" {
        portfolio.Image = constanta.DIR_FILE + filename
    }
    
    oldName := portfolio.Title

	var updatedInput models.Portfolio
	updatedInput.Image = portfolio.Image
	updatedInput.Title = input.Title
	updatedInput.Category_homeID = input.Id_category
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&portfolio).Updates(updatedInput)

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	SendResponse(c, portfolio, "success")
    activityMessage := "Update portfolio:'" + oldName + "' to '" + input.Title + "'"
    activitylog(c,activityMessage)

}

func DeletePortfolio(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var portfolio models.Portfolio
	if err := db.Where("id = ?", c.Param("id")).First(&portfolio).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&portfolio).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

	// Return success response
	SendResponse(c, portfolio, "success")
    activityMessage := "Delete portfolio: "+ portfolio.Title
    activitylog(c,activityMessage)
}