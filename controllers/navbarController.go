package controllers

import (
	"api-golang-compro/models"
	"math"

	// "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NavbarInput struct {
	Name        string    `json:"name"`
	Link_button string    `gorm:"type:text" json:"link_button"`
	Flag        uint      `json:"flag"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetAllNavbar(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var navbar []models.Navbar

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
    query.Model(&navbar).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&navbar).Error
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

    response := map[string]interface{}{
        "data":        navbar,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all navbar",response)
}


func CreateNavbar(c *gin.Context) {
	// Validate input
	var input NavbarInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	// Create
	navbar := models.Navbar{Name: input.Name, Link_button: input.Link_button, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&navbar)

	SendResponse(c, navbar, "success")
	activityMessage := "Create navbar: " +input.Name
    activitylog(c,activityMessage)
}

func GetNavbarById(c *gin.Context) { // Get model if exist
	var navbar models.Navbar

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&navbar).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c,"Get navbar by id "+ c.Param("id"),navbar)
}

func UpdateNavbar(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var navbar models.Navbar
	if err := db.Where("id = ?", c.Param("id")).First(&navbar).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input NavbarInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}
	oldName := navbar.Name
	
	var updatedInput models.Navbar
	updatedInput.Name = input.Name
	updatedInput.Link_button = input.Link_button
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&navbar).Updates(updatedInput)

	SendResponse(c, navbar, "success")
    activityMessage := "Update navbar:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeleteNavbar(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var navbar models.Navbar
	if err := db.Where("id = ?", c.Param("id")).First(&navbar).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&navbar).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, navbar, "success")
	activityMessage := "Delete navbar: "+ navbar.Name
    activitylog(c,activityMessage)
}