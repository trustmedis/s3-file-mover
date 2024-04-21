package main

import (
	"log"

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

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !event.Has(fsnotify.Remove) {
					lib.UploadFile(config, event.Name, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	for _, path := range config.WATCH_DIR {
		err = watcher.Add(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Block main goroutine forever.
	log.Println("Watching", config.WATCH_DIR)
	<-make(chan struct{})
}
