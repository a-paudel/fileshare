package models

import (
	"math/rand"
	"os"
	"time"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	Name string
	Path string
	Code string `gorm:"uniqueIndex"`
}

func generateCode() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	var code = make([]rune, 3)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	// check if code exists
	var file File
	err := Db.Where("code = ?", code).First(&file).Error
	if err == nil {
		// file exists, generate new code
		return generateCode()
	}
	return string(code)
}

func (f *File) DeleteExpired() {
	var before = time.Now().Add(-time.Hour * 24)
	var files []File
	Db.Where("created_at < ?", before).Find(&files)
	for _, file := range files {
		os.Remove(file.Path)
		Db.Delete(&file)
	}
}

func (f *File) Create(filename string, content []byte) {
	var code = generateCode()
	// var extension = strings.Split(filename, ".")[1]
	// save file
	var path = "files/" + code

	// create folder if not exists
	os.Mkdir("files", os.ModePerm)
	os.WriteFile(path, content, os.ModePerm)

	f.Name = filename
	f.Code = code
	f.Path = path

	Db.Create(f)
}
