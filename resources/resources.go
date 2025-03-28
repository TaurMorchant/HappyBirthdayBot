package res

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"math/rand"
)

//go:embed *
var templateFS embed.FS

type ImageKey string

var (
	No_picture        ImageKey = ""
	Vjuh              ImageKey = "img/cats/vjuh"
	Many_of_cats      ImageKey = "img/cats/many_cats"
	Cool_cat          ImageKey = "img/cats/cool_cat"
	Sad_cat           ImageKey = "img/cats/sad_cat"
	Do_not_understand ImageKey = "img/cats/do_not_understand"
	Do_not_scream     ImageKey = "img/cats/do_not_scream"
	Suspicious_cat    ImageKey = "img/cats/suspicious_cat"
	Wishlist          ImageKey = "img/cats/wishlist"
	Angry_cats        ImageKey = "img/cats/angry_cats"
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
