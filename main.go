package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/ewen-lbh/portfolio/pages"
	"github.com/ewen-lbh/portfolio/shared"
	"github.com/fatih/color"
	ortfodb "github.com/ortfo/db"
	"gopkg.in/yaml.v3"
)

func startFileServer(port int, root string) {
	staticServer := http.NewServeMux()
	fs := http.FileServer(http.Dir(filepath.Join(".", root)))
	staticServer.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf(":%d", port), staticServer)
}

func startServer(db ortfodb.Database, collections shared.Collections, sites []shared.Site, technologies []shared.Technology, tags []shared.Tag, locale string, port int) {
	server := http.NewServeMux()

	registeredPaths := []string{}

	handlePage := func(path string, page templ.Component) {
		for _, registeredPath := range registeredPaths {
			if registeredPath == path {
				color.Yellow("[%s] Page /%s is already registered, skipping", locale, path)
				return
			}
		}
		fmt.Printf("[%s] Registering page /%s\n", locale, path)
		server.Handle(fmt.Sprintf("/%s", path), templ.Handler(pages.Layout(page, collections.URLsToNames(true, locale), sites, locale)))
		registeredPaths = append(registeredPaths, path)
	}

	redirect := func(from, to string) {
		if !strings.HasPrefix(to, "https://") && !strings.HasPrefix(to, "mailto:") {
			to = fmt.Sprintf("/%s", to)
		}
		fmt.Printf("[%s] Registering redirect /%s -> %s\n", locale, from, to)
		server.Handle(fmt.Sprintf("/%s", from), http.RedirectHandler(to, http.StatusSeeOther))
	}

	handlePage("", pages.Index(db, locale))
	for _, work := range db.Works() {
		handlePage(work.ID, pages.Work(work, locale))
	}

	for id, collection := range collections {
		handlePage(id, pages.Collection(collection, db, tags, technologies, locale))
		for _, pathname := range collection.Aliases {
			redirect(pathname, id)
		}
	}

	for _, technology := range technologies {
		handlePage("using/"+technology.Slug, pages.TechnologyPage(technology, db, locale))
		for _, pathname := range technology.Aliases {
			redirect("using/"+pathname, "using/"+technology.Slug)
		}
	}

	for _, tag := range tags {
		handlePage(tag.URLName(), pages.TagPage(tag, db, locale))
		for _, pathname := range tag.Aliases {
			redirect(pathname, tag.URLName())
		}
	}

	for _, site := range sites {
		redirect(filepath.Join("to", site.Name), site.URL)
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
	db, locales := loadDatabase()

	if db.Partial() {
		color.Yellow("[!!] Database is partial, some works are missing data.")
	}

	if shared.IsDev() {
		fmt.Println("[  ] Running in dev mode")
	} else {
		fmt.Println("[  ] Running in production mode")
	}

	var collections map[string]shared.Collection
	loadDataFileMap("collections.yaml", &collections)

	var sites []shared.Site
	loadDataFile("sites.yaml", &sites)

	var technologies []shared.Technology
	loadDataFile("technologies.yaml", &technologies)

	var tags []shared.Tag
	loadDataFile("tags.yaml", &tags)

	var wg sync.WaitGroup
	wg.Add(1)

	for i, locale := range locales {
		go startServer(db, collections, sites, technologies, tags, locale, 8081+i)
	}
	go startFileServer(8079, "public")
	go startFileServer(8080, "media")

	wg.Wait()
}

func loadDatabase() (ortfodb.Database, []string) {
	var db ortfodb.Database
	databaseRaw, err := os.ReadFile(filepath.Join(".", "database.json"))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(databaseRaw, &db)
	if err != nil {
		panic(err)
	}
	locales := db.Languages()
	sort.Strings(locales)
	fmt.Printf("[  ] Works database has %d works in %d locales (%s)\n", len(db.Works()), len(locales), strings.Join(locales, ", "))
	return db, locales
}
