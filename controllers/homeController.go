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

type homeInput struct {
	// Logo                       string    `gorm:"text" json:"logo"`
	// Background_image_section_1 string    `gorm:"text" json:"background_image_section_1"`
	Title_section_1            string    `json:"title_section_1"`
	Description_section_1      string    `gorm:"text" json:"description_section_1"`
	Button_section_1           string    `json:"button_section_1"`
	Sub_title_section_2        string    `json:"sub_title_section_2"`
	Title_section_2            string    `json:"title_section_2"`
	Description_section_2      string    `gorm:"text" json:"description_section_2"`
	Button_section_2           string    `json:"button_section_2"`
	// Image_section_2            string    `gorm:"text" json:"image_section_2"`
	Sub_title_section_3        string    `json:"sub_title_section_3"`
	Title_section_3            string    `json:"title_section_3"`
	Sub_title_section_4        string    `json:"sub_title_section_4"`
	Title_section_4            string    `json:"title_section_4"`
	Description_contact_us     string    `gorm:"text" json:"description_contact_us"`
	Button_contact_us          string    `json:"button_contact_us"`
	Sub_title_section_5        string    `json:"sub_title_section_5"`
	Title_section_5            string    `json:"title_section_5"`
	Button_section_5           string    `json:"button_section_5"`
	Link_facebook              string    `gorm:"text" json:"link_facebook"`
	Link_linkedln              string    `gorm:"text" json:"link_linkedln"`
	Link_instagram             string    `gorm:"text" json:"link_instagram"`
	Flag                       uint      `json:"flag"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

func GetAllHome(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
    var home []models.Home

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
    query.Model(&home).Where("flag = 1").Count(&totalCount)

    // Calculate the total pages
    totalPages := int(math.Ceil(float64(totalCount) / float64(pagination.PerPage)))

    // Calculate the offset for pagination
    offset := (pagination.Page - 1) * pagination.PerPage

    // Apply pagination and sorting
    err := query.Order("id " + sortOrder).Offset(offset).Limit(pagination.PerPage).Find(&home).Error
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
    for i := range home {
        home[i].Logo = serverAddress + home[i].Logo
        home[i].Background_image_section_1 = serverAddress + home[i].Background_image_section_1
        home[i].Image_section_2 = serverAddress + home[i].Image_section_2
    }

    response := map[string]interface{}{
        "data":        home,
        "current_page": pagination.Page,
        "last_page":   lastPage,
        "per_page":    pagination.PerPage,
        "nextPage":    nextPage,
        "prevPage":    prevPage,
        "totalPages":  totalPages,
        "totalCount":  totalCount,
    }

    checkAndLogActivity(c,"Get all home",response)
}


func CreateHome(c *gin.Context) {
	// Validate input
	var input homeInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename1, errfilename1 := HandleUploadFile(c, "Logo")
	if errfilename1 != nil {
    SendError(c, "Upload error", errfilename1.Error())
    return
	}
	
	filename2, errfilename2 := HandleUploadFile(c, "Background_image_section_1")
	if errfilename2 != nil {
    SendError(c, "Upload error", errfilename2.Error())
    return
	}
	
	filename3, errfilename3 := HandleUploadFile(c, "Image_section_2")
	if errfilename3 != nil {
    SendError(c, "Upload error", errfilename3.Error())
    return
	}
	logo := constanta.DIR_FILE + filename1
	background_image_section_1 := constanta.DIR_FILE + filename2
	image_section_2 := constanta.DIR_FILE + filename3

	// Create
	home := models.Home{
		Logo:                       logo,
		Background_image_section_1: background_image_section_1,
		Title_section_1:            input.Title_section_1,
		Description_section_1:      input.Description_section_1,
		Button_section_1:           input.Button_section_1,
		Sub_title_section_2:        input.Sub_title_section_2,
		Title_section_2:            input.Title_section_2,
		Description_section_2:      input.Description_section_2,
		Button_section_2:           input.Button_section_2,
		Image_section_2:            image_section_2,
		Sub_title_section_3:        input.Sub_title_section_3,
		Title_section_3:            input.Title_section_3,
		Sub_title_section_4:        input.Sub_title_section_4,
		Title_section_4:            input.Title_section_4,
		Description_contact_us:     input.Description_contact_us,
		Button_contact_us:          input.Button_contact_us,
		Sub_title_section_5:        input.Sub_title_section_5,
		Title_section_5:            input.Title_section_5,
		Button_section_5:           input.Button_section_5,
		Link_facebook:              input.Link_facebook,
		Link_linkedln:              input.Link_linkedln,
		Link_instagram:             input.Link_instagram,
		Flag:                       1,
		CreatedAt:                  time.Now(),}
	db := c.MustGet("db").(*gorm.DB)
	db.Create(&home)

	SendResponse(c, home, "success")
	activityMessage := "Create home: " +input.Title_section_1
    activitylog(c,activityMessage)
}

func GetHomeByid(c *gin.Context) { // Get model if exist
	var home models.Home

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&home).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	serverAddress := "http://" + c.Request.Host

    // Mengubah setiap entri dalam data portofolio untuk menambahkan URL lengkap
    home.Logo = serverAddress + home.Logo
    home.Background_image_section_1 = serverAddress + home.Background_image_section_1
    home.Image_section_2 = serverAddress + home.Image_section_2

	checkAndLogActivity(c,"Get home by id "+ c.Param("id"),home)
}

func UpdateHome(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var home models.Home
	if err := db.Where("id = ?", c.Param("id")).First(&home).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Validate input
	var input homeInput
	if err := c.ShouldBind(&input); err != nil {
		SendError(c, "error", err.Error())
		return
	}

	filename1, errfilename1 := HandleUploadFile(c, "Logo")
    if errfilename1 != nil {
        SendError(c, "File upload error", errfilename1.Error())
        return
    }
    // Cek apakah ada file yang diunggah
    if filename1 != "" {
        // Hapus gambar lama jika ada
        if home.Logo != "" {
            oldImage := "." + home.Logo
            utils.RemoveFile(oldImage)
        }
        // Jika ada file yang diunggah, set nama file yang baru
        home.Logo = constanta.DIR_FILE + filename1
    }

    filename2, errfilename2 := HandleUploadFile(c, "Background_image_section_1")
    if errfilename2 != nil {
        SendError(c, "File upload error", errfilename2.Error())
        return
    }
    // Cek apakah ada file yang diunggah
    if filename2 != "" {
        // Hapus gambar lama jika ada
        if home.Background_image_section_1 != "" {
            oldIcon := "." + home.Background_image_section_1
            utils.RemoveFile(oldIcon)
        }
        // Jika ada file yang diunggah, set nama file yang baru
        home.Background_image_section_1 = constanta.DIR_FILE + filename2
    }
	
    filename3, errfilename3 := HandleUploadFile(c, "Image_section_2")
    if errfilename3 != nil {
        SendError(c, "File upload error", errfilename3.Error())
        return
    }
    // Cek apakah ada file yang diunggah
    if filename3 != "" {
        // Hapus gambar lama jika ada
        if home.Image_section_2 != "" {
            oldIcon := "." + home.Image_section_2
            utils.RemoveFile(oldIcon)
        }
        // Jika ada file yang diunggah, set nama file yang baru
        home.Image_section_2 = constanta.DIR_FILE + filename3
    }

	
    oldName := home.Title_section_1

	var updatedInput models.Home
		updatedInput.Logo=                       home.Logo
		updatedInput.Background_image_section_1= home.Background_image_section_1
		updatedInput.Title_section_1=            input.Title_section_1
		updatedInput.Description_section_1=      input.Description_section_1
		updatedInput.Button_section_1=           input.Button_section_1
		updatedInput.Sub_title_section_2=        input.Sub_title_section_2
		updatedInput.Title_section_2=            input.Title_section_2
		updatedInput.Description_section_2=      input.Description_section_2
		updatedInput.Button_section_2=           input.Button_section_2
		updatedInput.Image_section_2=            home.Image_section_2
		updatedInput.Sub_title_section_3=        input.Sub_title_section_3
		updatedInput.Title_section_3=            input.Title_section_3
		updatedInput.Sub_title_section_4=        input.Sub_title_section_4
		updatedInput.Title_section_4=            input.Title_section_4
		updatedInput.Description_contact_us=     input.Description_contact_us
		updatedInput.Button_contact_us=          input.Button_contact_us
		updatedInput.Sub_title_section_5=        input.Sub_title_section_5
		updatedInput.Title_section_5=            input.Title_section_5
		updatedInput.Button_section_5=           input.Button_section_5
		updatedInput.Link_facebook=              input.Link_facebook
		updatedInput.Link_linkedln=              input.Link_linkedln
		updatedInput.Link_instagram=             input.Link_instagram
		updatedInput.Flag=                       input.Flag
		updatedInput.UpdatedAt=                  time.Now()

	db.Model(&home).Updates(updatedInput)

	SendResponse(c, home, "success")
	activityMessage := "Update home:'" + oldName + "' to '" + input.Title_section_1 + "'"
    activitylog(c,activityMessage)
}

func DeleteHome(c *gin.Context) {
	// Get model if exiForm
	db := c.MustGet("db").(*gorm.DB)
	var home models.Home
	if err := db.Where("id = ?", c.Param("id")).First(&home).Error; err != nil {
		SendError(c, "Record not found", err.Error())
		return
	}

	// Set the flag to 0
	if err := db.Model(&home).Update("flag", 0).Error; err != nil {
		SendError(c, "Failed to delete", err.Error())
		return
	}

	// Return success response
	SendResponse(c, home, "success")
	activityMessage := "Delete home: "+ home.Title_section_1
    activitylog(c,activityMessage)
}