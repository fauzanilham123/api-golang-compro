package controllers

import (
	"api-golang-compro/models"
	"math"

	// "net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type formInput struct {
		Name      			string    `json:"name"`
		Email      			string    `json:"email"`
		Message     		string    `json:"message"`
		Flag     			uint      `json:"flag"`
		CreatedAt 			time.Time `json:"created_at"`
		UpdatedAt 			time.Time `json:"updated_at"`
	}


func GetAllForm(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var form []models.Form

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
    query.Model(&form).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&form).Error
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
        "data":        form,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all form",response)
}


func CreateForm(c *gin.Context) {
	// Validate input
	var input formInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	// Create
	form := models.Form{Name: input.Name, Email: input.Email, Message: input.Message, Flag: 1, CreatedAt: time.Now()}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&form)

	SendResponse(c, form, "success")
	activityMessage := "Create form: " +input.Name
    activitylog(c,activityMessage)
}

func GetFormById(c *gin.Context) { // Get model if exist
	var form models.Form

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	checkAndLogActivity(c,"Get form by id "+ c.Param("id"),form)
}

func UpdateForm(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var form models.Form
	if err := db.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input formInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}
	oldName := form.Name

	var updatedInput models.Form
	updatedInput.Name = input.Name
	updatedInput.Email = input.Email
	updatedInput.Message = input.Message
	updatedInput.Flag = input.Flag
	updatedInput.UpdatedAt = time.Now()

	db.Model(&form).Updates(updatedInput)

	SendResponse(c, form, "success")
    activityMessage := "Update form:'" + oldName + "' to '" + input.Name + "'"
    activitylog(c,activityMessage)
}

func DeleteForm(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var form models.Form
	if err := db.Where("id = ?", c.Param("id")).First(&form).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&form).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, form, "success")
	activityMessage := "Delete form: "+ form.Name
    activitylog(c,activityMessage)
}