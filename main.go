package main

import (
	"fmt"
  "sync"
  "net/http"
  "os"
  "io"
  "strings"
  "strconv"
)

const taskLimit = 20
const imageDir = "images"
const pokedexLast = 1010

func Task(id int) {
  idStr := padStart(strconv.Itoa(id), "0", 3)
  path := getImageFilePath(idStr)
  url := getImageUrl(idStr)

  if _, err := os.Stat(path); os.IsNotExist(err) {
    downloadImage(url, path)
  }
}

func mkdir(path string) error {
  if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

  return nil
}

func getImageUrl(id string) string {
  return fmt.Sprintf("https://www.serebii.net/pokedex-sv/icon/%s.png", id)
}

func getImageFileName(id string) string {
  return fmt.Sprintf("%s.png", id)
}

func getImageFilePath(id string) string {
  return fmt.Sprintf("%s/%s", imageDir, getImageFileName(id))
}

func padStart(str, pad string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(pad, length-len(str)) + str
}

func downloadImage(url string, filepath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request error: %s", response.Status)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
  err := mkdir(imageDir)
  if err != nil {
    fmt.Println(err)
    return
  }

	var wg sync.WaitGroup
	taskLimit := make(chan struct{}, 20)

	for i := 1; i <= pokedexLast; i++ {
		wg.Add(1)
		taskLimit <- struct{}{}

		go func(id int) {
			defer wg.Done()
			Task(id)
			<-taskLimit
		}(i)
	}

	wg.Wait()
	fmt.Println("All tasks completed.")
}
