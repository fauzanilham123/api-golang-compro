package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
 
type BaseController struct {
}
 
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type Pagination struct {
    Page    int
    PerPage int
}
 
func  SendResponse(ctx *gin.Context, result interface{}, message string) {
	response := Response{
		Success: true,
		Code:    200,
		Message: message,
		Result:  result,
	}

	ctx.JSON(http.StatusOK, response)
}

func SendError(ctx *gin.Context, errorMessage string, errorMessages string) {
	response := Response{
		Success: false,
		Code:   400,
		Message: errorMessage,
		Result:  nil,
	}

	if len(errorMessages) > 0 {
		response.Result = errorMessages
	}

	ctx.JSON(http.StatusBadRequest, response)
}
 
 
func ExtractPagination(c *gin.Context) Pagination {
    page := c.DefaultQuery("page", "1")
    perPage := c.DefaultQuery("perPage", "10")

    pageInt, err := strconv.Atoi(page)
    if err != nil || pageInt < 1 {
        pageInt = 1
    }

    perPageInt, err := strconv.Atoi(perPage)
    if err != nil || perPageInt < 1 {
        perPageInt = 10
    }

    return Pagination{Page: pageInt, PerPage: perPageInt}
}

func PaginateQuery(db *gorm.DB, pagination Pagination) *gorm.DB {
    return db.Offset((pagination.Page - 1) * pagination.PerPage).Limit(pagination.PerPage)
}