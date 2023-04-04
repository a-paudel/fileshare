package routes

import (
	"fileshare/models"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func RegisterFileRoutes(app *echo.Echo) {
	router := app.Group("")

	router.GET("", uploadGet)
	router.POST("", uploadPost)

	router.GET("/:code", view)
	router.GET("/:code/download", download)
}

// upload get
func uploadGet(c echo.Context) error {
	return c.Render(200, "upload", nil)
}

// upload post
func uploadPost(c echo.Context) error {
	htmlFile, err := c.FormFile("file")
	if err != nil {
		return c.String(400, "file not found"+"\n"+err.Error())
	}
	// file size limit 100MB
	if htmlFile.Size > 100*1000*1000 {
		return c.Render(400, "errors/file_size", nil)
	}

	src, err := htmlFile.Open()
	if err != nil {
		return c.String(400, "file could not be opened"+"\n"+err.Error())
	}
	defer src.Close()
	var content = make([]byte, htmlFile.Size)
	_, err = src.Read(content)
	if err != nil {
		return c.String(400, "file could not be read"+"\n"+err.Error())
	}

	var file models.File
	file.Create(htmlFile.Filename, content)

	var url = c.Echo().URL(view, file.Code)
	return c.Redirect(302, url)
}

// view
func view(c echo.Context) error {
	var file models.File
	file.DeleteExpired()
	err := models.Db.Where("code = ?", c.Param("code")).First(&file).Error

	if err != nil {
		return c.Render(404, "errors/not_found", nil)

	}
	var expiryDate = file.CreatedAt.Add(time.Hour * 24)
	var downloadUrl = c.Echo().URL(download, file.Code)
	var data = map[string]interface{}{
		"file":        file,
		"expiryDate":  expiryDate.Format("2006-01-02 15:04:05"),
		"downloadUrl": downloadUrl,
	}

	return c.Render(200, "view", data)
}

// download
func download(c echo.Context) error {
	var file models.File
	file.DeleteExpired()

	err := models.Db.Where("code = ?", c.Param("code")).First(&file).Error
	if err != nil {
		return c.Render(404, "errors/not_found", nil)
	}

	defer func() {
		os.Remove(file.Path)
		models.Db.Delete(&file)
	}()

	return c.Attachment(file.Path, file.Name)
}
