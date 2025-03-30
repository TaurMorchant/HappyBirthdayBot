package res

import (
	"embed"
	"encoding/csv"
	"errors"
	"io/fs"
	"log"
	"math/rand"
)

//go:embed *
var templateFS embed.FS

type ImageKey string

var (
	NoPicture     ImageKey = ""
	Vjuh          ImageKey = "img/cats/vjuh"
	Many          ImageKey = "img/cats/many"
	Cool          ImageKey = "img/cats/cool"
	Sad           ImageKey = "img/cats/sad"
	Error         ImageKey = "img/cats/error"
	DoNotScream   ImageKey = "img/cats/do_not_scream"
	Suspicious    ImageKey = "img/cats/suspicious"
	Wishlist      ImageKey = "img/cats/wishlist"
	Angry         ImageKey = "img/cats/angry"
	Ok            ImageKey = "img/cats/ok"
	Waiting       ImageKey = "img/cats/waiting"
	HappyBirthday ImageKey = "img/cats/happy_birthday"
	Random        ImageKey = "img/cats/random"
	Main          ImageKey = "img/cats/main.jpg"
)

func GetImage(imageKey ImageKey) ([]byte, bool) {
	log.Printf("[DEBUG] get image %s", imageKey)
	if isDir(imageKey) {
		log.Printf("[DEBUG] it is directory")
		return getRandomImage(imageKey)
	} else {
		log.Printf("[DEBUG] it is file")
		result, err := templateFS.ReadFile(string(imageKey))
		if err != nil {
			log.Printf("[ERROR] get image %s failed, err = %s", imageKey, err.Error())
			return nil, false
		}
		return result, true
	}
}

func ReadCSV(filename string) ([][]string, error) {
	file, err := templateFS.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	return reader.ReadAll()
}

func getRandomImage(imageKey ImageKey) ([]byte, bool) {
	files, err := templateFS.ReadDir(string(imageKey))
	if err != nil {
		log.Printf("[ERROR] get random image %s failed, err = %s", imageKey, err.Error())
		return nil, false
	}

	// Фильтруем только файлы (исключаем поддиректории)
	var images []string
	for _, file := range files {
		if !file.IsDir() {
			images = append(images, file.Name())
		}
	}

	if len(images) == 0 {
		return nil, false
	}

	randomImage := images[rand.Intn(len(images))]

	imagePath := string(imageKey) + "/" + randomImage
	result, err := templateFS.ReadFile(imagePath)
	if err != nil {
		log.Printf("[ERROR] get image %s failed, err = %s", imagePath, err.Error())
		return nil, false
	}
	return result, true
}

func isDir(key ImageKey) bool {
	_, err := templateFS.ReadFile(string(key))
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			return true
		}
	}
	return false
}
