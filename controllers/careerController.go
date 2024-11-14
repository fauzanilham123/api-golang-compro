package controllers

import (
	"api-golang-compro/models"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type careerInput struct {
	Id_category uint   `json:"id_category"`
	Id_position uint   `json:"id_position"`
	Name        string `json:"name"`
	Desc        string `json:"description"`
	Required    string `json:"required"`
	Flag        uint   `json:"flag"`
}



func  GetAllCareer(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var career []models.Career

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

    // Count the total number of records
    var totalCount int64
    query.Model(&career).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Preload("Category").Preload("Position").Find(&career).Error
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
        "data":        career,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all career", response)
}

func CreateCareer(c *gin.Context) {
    // Validate input
    var input careerInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }

    // Create 
    career := models.Career{CategoryID: input.Id_category,PositionID: input.Id_position, Name: input.Name, Description: input.Desc, Required: input.Required,Flag: 1,CreatedAt: time.Now()}
    db := c.MustGet("db").(*gorm.DB)
    db.Create(&career)

    err := db.Preload("Category").Preload("Position").First(&career, career.ID).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }
    
    SendResponse(c, career, "success")
    activityMessage := "Create career: " +input.Name
    activitylog(c,activityMessage)
    
}

func GetCareerById(c *gin.Context) { // Get model if exist
    var career models.Career

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&career).Error; err != nil {
        SendError(c,"Record not found",err.Error())
        return
    }

    err := db.Preload("Category").Preload("Position").First(&career, career.ID).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

    checkAndLogActivity(c,"Get career by id "+ c.Param("id"),career)
    
}

func UpdateCareer(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var career models.Career
    if err := db.Where("id = ?", c.Param("id")).First(&career).Error; err != nil {
        SendError(c,"Record not found",err.Error())
        return
    }

    // Validate input
    var input careerInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }

    oldName := career.Name

    var updatedInput models.Career
    updatedInput.CategoryID = input.Id_category
    updatedInput.PositionID = input.Id_position
    updatedInput.Name = input.Name
    updatedInput.Description = input.Desc
    updatedInput.Required = input.Required
    updatedInput.Flag = input.Flag
    updatedInput.UpdatedAt = time.Now()

    db.Model(&career).Updates(updatedInput)

    err := db.Preload("Category").Preload("Position").First(&career, career.ID).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

    SendResponse(c, career, "success")
    activityMessage := "Update career:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeleteCareer(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var career models.Career
    if err := db.Where("id = ?", c.Param("id")).First(&career).Error; err != nil {
        SendError(c, "Record not found", err.Error())
        return
    }

    // Set the flag to 0
    if err := db.Model(&career).Update("flag", 0).Error; err != nil {
        SendError(c, "Failed to delete", err.Error())
        return
    }

    err := db.Preload("Category").Preload("Position").First(&career, career.ID).Error
    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

    // Return success response
    SendResponse(c, career, "success")
    activityMessage := "Delete career: "+ career.Name
    activitylog(c,activityMessage)
}

