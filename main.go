package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/trustmedis/s3-file-mover/lib"
)

func main() {
	config := lib.LoadConfig()
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Upload existing files based on config
	if config.AUTOMOVE_EXISTING_FILES {
		for _, path := range config.WATCH_DIR {
			// scan files in each directory
			files, err := filepath.Glob(filepath.Join(path, "*"))
			if err != nil {
				log.Println("Error scanning files: ", err)
			}
			for _, file := range files {
				targetFilePath := filepath.Base(file)
				lib.UploadFile(config, file, targetFilePath)
			}
		}
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !event.Has(fsnotify.Remove) {
					targetFilePath := filepath.Base(event.Name)
					lib.UploadFile(config, event.Name, targetFilePath)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add paths to watch.
	for _, path := range config.WATCH_DIR {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal("ERROR : ", err)
		}
	}

	// Block main goroutine forever.
	if config.AUTO_CLEANUP {
		log.Println("**WARNING** : Auto cleanup is enabled. This will remove all files in", config.WATCH_DIR, "immediately.")
	}
	log.Println("Watching", strings.Join(config.WATCH_DIR, ", "))
	<-make(chan struct{})
}
