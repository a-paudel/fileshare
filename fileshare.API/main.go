package main

import (
	"context"
	"encoding/json"
	"fileshare/models"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

var ctx = context.Background()

func generateCode() string {
	// three letter code
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	code := make([]rune, 3)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	// check if code exists in db
	exists, _ := models.Q.CheckIfCodeExists(ctx, string(code))
	if exists > 0 {
		// code exists, generate new code
		return generateCode()
	}
	return string(code)

}

// middleware to delete expired files
func deleteExpiredFilesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get list of expired files
		fmt.Println("Deleting expired files")
		expiredFiles, _ := models.Q.GetExpiredFiles(ctx)
		// for each file delete from disk and db
		for _, file := range expiredFiles {
			// remove from disk
			os.RemoveAll(fmt.Sprintf("./data/files/%s", file.Code))
			// remove from db
			models.Q.DeleteFileByCode(ctx, file.Code)
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	startServer()
}

func startServer() {
	app := chi.NewRouter()
	app.Use(cors.Default().Handler)
	app.Use(middleware.Logger)
	app.Use(deleteExpiredFilesMiddleware)

	app.Post("/api/files", addFile)
	app.Get("/api/files/{code}", getFileDetails)
	app.Get("/api/files/{code}/download", getFile)

	fmt.Println("Server started at http://localhost:8000")
	http.ListenAndServe(":8000", app)
}

func addFile(w http.ResponseWriter, r *http.Request) {
	// get formFile from request
	formFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		// invalid input return error
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	filename := fileHeader.Filename
	// generate code
	code := generateCode()

	// save file to disk: ./data/files/code/filename
	os.MkdirAll("./data/files/"+code, os.ModePerm)
	diskFile, err := os.Create("./data/files/" + code + "/" + filename)
	if err != nil {
		http.Error(w, "Unable to create file on disk", http.StatusInternalServerError)
		return
	}
	defer diskFile.Close()

	io.Copy(diskFile, formFile)

	// save file to db
	err = models.Q.CreateFile(ctx, models.CreateFileParams{
		Code:     code,
		Filename: filename,
		Filesize: fileHeader.Size,
	})
	if err != nil {
		http.Error(w, "Unable to save file to db", http.StatusInternalServerError)
		return
	}

	// get file from db
	dbFile, err := models.Q.GetFileByCode(ctx, code)
	if err != nil {
		http.Error(w, "Unable to get file from db", http.StatusInternalServerError)
		return
	}
	// return code to user
	response, err := json.Marshal(dbFile)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func getFileDetails(w http.ResponseWriter, r *http.Request) {
	// get code from path param
	code := chi.URLParam(r, "code")
	// get file from db
	file, err := models.Q.GetFileByCode(ctx, code)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	// write file details to response
	response, err := json.Marshal(file)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	// get code from path param
	code := chi.URLParam(r, "code")

	// get file from db
	dbFile, err := models.Q.GetFileByCode(ctx, code)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := fmt.Sprintf("./data/files/%s/%s", code, dbFile.Filename)
	// serve file from disk: ./data/files/code/filename
	diskFile, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// defer deletion of file from disk and db
	defer func() {
		os.RemoveAll(fmt.Sprintf("./data/files/%s", code))
		models.Q.DeleteFileByCode(ctx, code)
	}()

	// write file headers to response
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", dbFile.Filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", dbFile.Filesize))

	// write file to response
	io.Copy(w, diskFile)
}
