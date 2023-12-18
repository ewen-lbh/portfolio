package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/ewen-lbh/portfolio/pages"
	"github.com/ewen-lbh/portfolio/shared"
	"github.com/fatih/color"
	ortfodb "github.com/ortfo/db"
	"gopkg.in/yaml.v3"
)

func startServer(db ortfodb.Database, collections shared.Collections, sites []shared.Site, technologies []shared.Technology, tags []shared.Tag, locale string, port int) {
	server := http.NewServeMux()

	layouted := func(page templ.Component) templ.Component {
		return pages.Layout(page, collections.URLsToNames(true, locale), sites, locale)
	}

	server.Handle("/", templ.Handler(layouted(pages.Index(db, locale))))

	registeredCount := 0
	for _, work := range db.Works() {
		server.Handle("/"+work.ID, templ.Handler(layouted(pages.Work(work, locale))))
		registeredCount++
	}
	fmt.Printf("[%s] Registered %d works\n", locale, registeredCount)

	for id, collection := range collections {
		fmt.Printf("[%s] Registering collection %s along with %d aliases\n", locale, id, len(collection.Aliases))
		server.Handle("/"+id, templ.Handler(layouted(pages.Collection(collection, db, tags, technologies, locale))))
		for _, pathname := range collection.Aliases {
			server.Handle("/"+pathname, templ.Handler(layouted(pages.Collection(collection, db, tags, technologies, locale))))
		}
	}

	for _, site := range sites {
		server.Handle("/to/"+site.Name, http.RedirectHandler(site.URL, http.StatusSeeOther))
	}

	fmt.Printf("[%s] Server started on http://localhost:%d\n", locale, port)
	http.ListenAndServe(":"+fmt.Sprint(port), server)
}

func loadDataFile[T any](path string, into *[]T) {
	raw, err := os.ReadFile(filepath.Join(".", path))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(raw, &into)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[  ] Loaded %d items from %s\n", len(*into), path)
}

func loadDataFileMap[K comparable, V any](path string, into *map[K]V) {
	raw, err := os.ReadFile(filepath.Join(".", path))
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(raw, &into)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[  ] Loaded %d items from %s\n", len(*into), path)
}

func main() {
	var db ortfodb.Database
	databaseRaw, err := os.ReadFile(filepath.Join(".", "database.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(databaseRaw, &db)
	if err != nil {
		panic(err)
	}

	if db.Partial() {
		color.Yellow("[!!] Database is partial, some works are missing data.")
	}

	if shared.IsDev() {
		fmt.Println("[  ] Running in dev mode")
	} else {
		fmt.Println("[  ] Running in production mode")
	}

	locales := db.Languages()
	fmt.Printf("[  ] Works database has %d works in %d locales (%s)\n", len(db.Works()), len(locales), strings.Join(locales, ", "))

	var collections map[string]shared.Collection
	loadDataFileMap("collections.yaml", &collections)

	var sites []shared.Site
	loadDataFile("sites.yaml", &sites)

	var technologies []shared.Technology
	loadDataFile("technologies.yaml", &technologies)

	var tags []shared.Tag
	loadDataFile("tags.yaml", &tags)

	var wg sync.WaitGroup
	wg.Add(len(locales))
	for i, locale := range db.Languages() {
		go startServer(db, collections, sites, technologies, tags, locale, 8081+i)
	}
	staticServer := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepath.Join("/home/ewen/projects/ortfo/db/dist/media")))
	staticServer.Handle("/", fs)
	http.ListenAndServe(":8080", staticServer)
	wg.Wait()
}
