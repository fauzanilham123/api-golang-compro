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

type portfolioHepytechInput struct {
	Id_category   uint   `json:"id_category"`
	Name          string `json:"name"`
	Description          string `gorm:"text" json:"description"`
	// Image         string `gorm:"text" json:"image"`
	Desc_problem  string `gorm:"text" json:"desc_problem"`
	Desc_solution string `gorm:"text" json:"desc_solution"`
	Impact_title  string `json:"impact_title"`
	Impact_desc   string `gorm:"text" json:"impact_desc"`
	// Impact_icon   string `gorm:"text" json:"impact_icon"`
	Slug          string `json:"slug"`
	Flag          uint   `json:"flag"`
}

func GetAllPortfolioHepytech(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var portfolio []models.PortfolioHepytech

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

	// Count the total number of records
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
		"data":         portfolio,
		"current_page": pagination.Page,
		"last_page":    lastPage,
		"per_page":     pagination.PerPage,
		"nextPage":     nextPage,
		"prevPage":     prevPage,
		"totalPages":   totalPages,
		"totalCount":   totalCount,
	}

	checkAndLogActivity(c,"Get all portfolio_hepytech",response)
}

func CreatePortfolioHepytech(c *gin.Context) {
	// Validate input
	var input portfolioHepytechInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename1, errfilename1 := HandleUploadFile(c, "Image")
	if errfilename1 != nil {
    SendError(c, "Upload error", errfilename1.Error())
    return
	}

	image := constanta.DIR_FILE + filename1
	

	// Create
	portfolio := models.PortfolioHepytech{CategoryID: input.Id_category, Name: input.Name, Description: input.Description, Image: image, Desc_problem: input.Desc_problem,Desc_solution: input.Desc_solution, Slug: input.Slug, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&portfolio)

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
	if err != nil {
		SendError(c, "error", err.Error())
		return
	}

	SendResponse(c, portfolio, "success")
	activityMessage := "Create portfolio_hepytech: " +input.Name
    activitylog(c,activityMessage)
}

func GetPortfolioHepytechById(c *gin.Context) { // Get model if exist
	var portfolio models.PortfolioHepytech

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

	checkAndLogActivity(c,"Get portfolio_hepytech by id "+ c.Param("id"),portfolio)
}

func UpdatePortfolioHepytech(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var portfolio models.PortfolioHepytech
	if err := db.Where("id = ?", c.Param("id")).First(&portfolio).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input portfolioHepytechInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename1, errfilename1 := HandleUploadFile(c, "Image")
    if errfilename1 != nil {
        SendError(c, "File upload error", errfilename1.Error())
        return
    }
    // Cek apakah ada file yang diunggah
    if filename1 != "" {
        // Hapus gambar lama jika ada
        if portfolio.Image != "" {
            oldImage := "." + portfolio.Image
            utils.RemoveFile(oldImage)
        }
        // Jika ada file yang diunggah, set nama file yang baru
        portfolio.Image = constanta.DIR_FILE + filename1
    }
	oldName := portfolio.Name
	
	var updatedInput models.PortfolioHepytech
	updatedInput.CategoryID = input.Id_category
	updatedInput.Name = input.Name
	updatedInput.Description = input.Description
	updatedInput.Image = portfolio.Image
	updatedInput.Desc_problem = input.Desc_problem
	updatedInput.Desc_solution = input.Desc_solution
	updatedInput.Slug = input.Slug
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&portfolio).Updates(updatedInput)

	err := db.Preload("Category").First(&portfolio, portfolio.Id).Error
	if err != nil {
		SendError(c, "error", err.Error())
		return
	}

	SendResponse(c, portfolio, "success")
    activityMessage := "Update portfolio_hepytech:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeletePortfolioHepytech(c *gin.Context) {
	// Get model if exist
	db := c.MustGet("db").(*gorm.DB)
	var portfolio models.PortfolioHepytech
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
	activityMessage := "Delete portfolio_hepytech: "+ portfolio.Name
    activitylog(c,activityMessage)
}

func GetPortfolioHepytechBySlug(c *gin.Context) { // Get model if exist
	var portfolio models.PortfolioHepytech

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("slug = ?", c.Param("slug")).First(&portfolio).Error; err != nil {
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

	checkAndLogActivity(c,"Get portfolio_hepytech by slug "+ c.Param("slug"),portfolio)
}