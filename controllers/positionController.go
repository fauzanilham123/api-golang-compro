package controllers

import (
	"api-golang-compro/models"
	"math"

	// "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type positionInput struct {
	Name string `json:"name"`
	Flag uint   `json:"flag"`
}

func GetAllPosition(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var position []models.Position

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
    query.Model(&position).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&position).Error
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
        "data":        position,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all position",response)
}


func CreatePosition(c *gin.Context) {
    // Validate input
    var input positionInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }

    // Create 
    position := models.Position{Name: input.Name, Flag: 1, CreatedAt: time.Now()}
    db := c.MustGet("db").(*gorm.DB)
    db.Create(&position)

    SendResponse(c, position, "success")
    activityMessage := "Create position: " +input.Name
    activitylog(c,activityMessage)
}


func GetPositionById(c *gin.Context) { // Get model if exist
    var position models.Position

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
        SendError(c,"error",err.Error())
        return
    }

    
    checkAndLogActivity(c,"Get position by id "+ c.Param("id"),position)
}

func UpdatePosition(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var position models.Position
    if err := db.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
        SendError(c,"error",err.Error())
        return
    }

    // Validate input
    var input positionInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }
    
    oldName := position.Name
    var updatedInput models.Position
    updatedInput.Name = input.Name
    updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()
    

    db.Model(&position).Updates(updatedInput)

    SendResponse(c, position, "success")
    activityMessage := "Update position:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeletePosition(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var position models.Position
    if err := db.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
        SendError(c, "Record not found", err.Error())
        return
    }

    // Set the flag to 0
    if err := db.Model(&position).Update("flag", 0).Error; err != nil {
        SendError(c, "Failed to delete", err.Error())
        return
    }

    // Return success response
    SendResponse(c, position, "success")
    activityMessage := "Delete position: "+ position.Name
    activitylog(c,activityMessage)
}