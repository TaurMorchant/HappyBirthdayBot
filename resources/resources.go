package res

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
)

//go:embed *
var templateFS embed.FS

type ImageKey string

const (
	Vjuh         ImageKey = "img/cats/vjuh.jpg"
	Many_of_cats ImageKey = "img/cats/many_of_cats.jpg"
	Cool_cat     ImageKey = "img/cats/cool_cat"
	Sad_cat      ImageKey = "img/cats/sad_cat"
)

func ReadFile() {
	// Чтение файла из встроенной FS
	data, err := templateFS.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

func GetImage(imageKey ImageKey) ([]byte, error) {
	if isDir(imageKey) {
		return getRandomImage(imageKey)
	} else {
		return templateFS.ReadFile(string(imageKey))
	}
}

func getRandomImage(imageKey ImageKey) ([]byte, error) {
	files, err := templateFS.ReadDir(string(imageKey))
	if err != nil {
		return nil, err
	}

	// Фильтруем только файлы (исключаем поддиректории)
	var images []string
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}

	if len(images) == 0 {
		return nil, err
	}

	randomImage := images[rand.Intn(len(images))]

	return templateFS.ReadFile(string(imageKey) + "/" + randomImage)
}

func isDir(key ImageKey) bool {
	_, err := templateFS.ReadFile(string(key))
	if err != nil {
		log.Println("err", err.Error())
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			return true
		}
	}
	return false
}
