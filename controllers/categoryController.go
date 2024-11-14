package controllers

import (
	"api-golang-compro/models"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type categoryInput struct {
	Name string `json:"name"`
	Flag uint   `json:"flag"`
}

func GetAllCategory(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var category []models.Category

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
    query.Model(category).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&category).Error
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
        "data":        category,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all category",response)
}

func CreateCategory(c *gin.Context) {
    // Validate input
    var input careerInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }

    // Create 
    category := models.Category{Name: input.Name, Flag: 1, CreatedAt: time.Now()}
    db := c.MustGet("db").(*gorm.DB)
    db.Create(&category)

    SendResponse(c, category, "success")
    activityMessage := "Create category: " +input.Name
    activitylog(c,activityMessage)
}

func GetCategoryById(c *gin.Context) { // Get model if exist
    var category models.Category

    db := c.MustGet("db").(*gorm.DB)
    if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
        SendError(c,"Record not found",err.Error())
        return
    }

    checkAndLogActivity(c,"Get category by id "+ c.Param("id"),category)
}

func UpdateCategory(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var category models.Category
    if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
        SendError(c,"Record not found",err.Error())
        return
    }

    // Validate input
    var input categoryInput
    if err := c.ShouldBind(&input); err != nil {
        SendError(c,"error",err.Error())
        return
    }
    oldName := category.Name

    var updatedInput models.Category
    updatedInput.Name = input.Name
    updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()
    

    db.Model(&category).Updates(updatedInput)

    SendResponse(c, category, "success")
    activityMessage := "Update category:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeleteCategory(c *gin.Context) {
    // Get model if exist
    db := c.MustGet("db").(*gorm.DB)
    var category models.Category
    if err := db.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
        SendError(c, "Record not found", err.Error())
        return
    }

    // Set the flag to 0
    if err := db.Model(&category).Update("flag", 0).Error; err != nil {
        SendError(c, "Failed to delete", err.Error())
        return
    }

    // Return success response
    SendResponse(c, category, "success")
    activityMessage := "Delete category: "+ category.Name
    activitylog(c,activityMessage)
}
