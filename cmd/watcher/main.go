package main

import (
	"MrAlbertX/server/internal/indexer"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"
)

type program struct {
	watcher *fsnotify.Watcher
	quit    chan struct{}
	indexer *indexer.ProgramIndexer
}

func (p *program) Start(s service.Service) error {
	p.quit = make(chan struct{})
	var err error
	p.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	log.Println("Service starting...")
	go p.run()
	return nil
}

func (p *program) run() {
	defer p.watcher.Close()

	var (
		debounceTimer *time.Timer
		debounceDelay = 10 * time.Second
	)

	reindexAction := func() {
		log.Println("Changes detected. Starting re-indexing...")
		if err := p.indexer.Reindex(); err != nil {
			log.Printf("Error during re-indexing: %v", err)
		} else {
			log.Println("Re-indexing completed successfully.")
		}
	}

	reindexAction()

	pathsToWatch := indexer.GetDefaultScanPaths()
	for _, path := range pathsToWatch {
		filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if err == nil && info.IsDir() {
				p.watcher.Add(subPath)
			}
			return nil
		})
	}
	log.Println("Watching initial paths.")

	for {
		select {
		case event, ok := <-p.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceDelay, reindexAction)
			}
		case err, ok := <-p.watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		case <-p.quit:
			log.Println("Service stopping.")
			return
		}
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.quit)
	log.Println("Service stopped.")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "MrAlbertX-Watcher",
		DisplayName: "Mr. Albert X Watcher Service",
		Description: "File system watcher for the Mr. Albert X assistant.",
	}

	prg := &program{}
	configDir, _ := os.UserConfigDir()
	cachePath := filepath.Join(configDir, "MrAlbertX", "program_index.json")
	prg.indexer = indexer.NewProgramIndexer(cachePath)

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
