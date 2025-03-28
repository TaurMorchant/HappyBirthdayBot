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

const (
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

func GetImage(imageKey ImageKey) ([]byte, error) {
	log.Printf("[DEBUG] get image %s", imageKey)
	if isDir(imageKey) {
		log.Printf("[DEBUG] it is directory")
		return getRandomImage(imageKey)
	} else {
		log.Printf("[DEBUG] it is file")
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
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			return true
		}
	}
	return false
}
