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
	"time"

	"github.com/a-h/templ"
	"github.com/ewen-lbh/portfolio/pages"
	pages_shs "github.com/ewen-lbh/portfolio/pages/shs"
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

func startSHSServer(wg *sync.WaitGroup, port int, regularSiteURL string, sites []shared.Site, db ortfodb.Database, collections shared.Collections, technologies []shared.Technology, tags []shared.Tag) {
	navigation := pages.Navigation([]pages.NavigationLink{
		{Text: "home", Link: "/"},
		{Text: "education", Link: "/education"},
		// {Text: "sustainability", Link: "/sustainability"},
		{Text: "engagement", Link: "/engagement"},
		{Text: "international", Link: "/international"},
		{Text: "career", Link: "/career"},
		{Text: "projects", Link: regularSiteURL},
		{Text: "ppp", Link: "https://media.ewen.works/ppp.pdf"},
	})

	server := http.NewServeMux()
	server.Handle("/", templ.Handler(pages.Layout(pages_shs.Home(db, collections, tags, technologies), navigation, sites, "en")))
	server.Handle("/education", templ.Handler(pages.Layout(pages_shs.EducationPage(), navigation, sites, "en")))
	server.Handle("/international", templ.Handler(pages.Layout(pages_shs.InternationalPage(), navigation, sites, "en")))
	server.Handle("/engagement", templ.Handler(pages.Layout(pages_shs.EngagementPage(), navigation, sites, "en")))
	server.Handle("/career", templ.Handler(pages.Layout(pages_shs.CareerPage(), navigation, sites, "en")))
	http.ListenAndServe(fmt.Sprintf(":%d", port), server)
}

func startPagesServer(wg *sync.WaitGroup, db ortfodb.Database, collections shared.Collections, sites []shared.Site, technologies []shared.Technology, tags []shared.Tag, blogEntries []shared.BlogEntry, translations *Translations, port int) {
	server := http.NewServeMux()

	registeredPaths := []string{}

	handlePage := func(path string, page templ.Component) {
		for _, registeredPath := range registeredPaths {
			if registeredPath == path {
				color.Yellow("[%s] Page /%s is already registered, skipping", translations.language, path)
				return
			}
		}

		if os.Getenv("ENV") == "production" {
			// Try to get rendered static page first
			content, err := GetPage(filepath.Join("dist", translations.language), path)
			if err != nil {
				color.Cyan("[%s] Could not get static page for /%s: %s, falling back to server rendering until next restart", translations.language, path, err)
			} else {
				server.Handle(fmt.Sprintf("/%s", path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write(content)
				}))
				registeredPaths = append(registeredPaths, path)
				return
			}
		}

		navigation := pages.Navigation(pages.RegularLayoutLinks(collections.URLsToNames(true, translations.language), translations.language, path))

		translator := HttpTranslator{
			translations: translations,
			ch:           templ.Handler(pages.Layout(page, navigation, sites, translations.language)),
		}

		server.Handle(fmt.Sprintf("/%s", path), translator)
		registeredPaths = append(registeredPaths, path)
	}

	redirect := func(from, to string) {
		if !strings.HasPrefix(to, "https://") && !strings.HasPrefix(to, "mailto:") && !strings.HasPrefix(to, "http://") {
			to = fmt.Sprintf("/%s", to)
		}
		// fmt.Printf("[%s] Registering redirect /%s -> %s\n", locale, from, to)
		server.Handle(fmt.Sprintf("/%s", from), http.RedirectHandler(to, http.StatusSeeOther))
	}

	handlePage("", pages.Index(db, translations.language))
	for _, work := range db.Works() {
		handlePage(work.ID, pages.Work(
			work,
			shared.TagsOf(tags, work.Metadata),
			shared.TechsOf(technologies, work.Metadata),
			collections.ThatIncludeWork(work, shared.Keys(db.Works()), tags, technologies),
			translations.language,
		))

		for _, alias := range work.Metadata.Aliases {
			redirect(alias, work.ID)
			redirect(fmt.Sprintf("%s.json", alias), fmt.Sprintf("%s.json", work.ID))
		}

		encoded, _ := json.Marshal(work)
		server.Handle(fmt.Sprintf("/%s.json", work.ID), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(encoded)
		}))
	}

	handlePage("blog", pages.BlogIndex(blogEntries))
	for _, entry := range blogEntries {
		err := entry.GetPageviews("/blog")
		if err != nil {
			fmt.Printf("[!!] Could not get pageviews for blog entry %s: %s\n", "blog/"+entry.Slug, err)
		}

		handlePage("blog/"+entry.Slug, pages.BlogEntry(entry, db))
	}

	for id, collection := range collections {
		handlePage(id, pages.Collection(collection, db, tags, technologies, translations.language))
		for _, pathname := range collection.Aliases {
			redirect(pathname, id)
		}
	}

	for _, technology := range technologies {
		handlePage("using/"+technology.Slug, pages.TechnologyPage(technology, db, translations.language))
		for _, pathname := range technology.Aliases {
			redirect("using/"+pathname, "using/"+technology.Slug)
		}
	}

	for _, tag := range tags {
		handlePage(tag.URLName(), pages.TagPage(tag, db, translations.language))
		for _, pathname := range tag.Aliases {
			redirect(pathname, tag.URLName())
		}
	}

	for _, site := range sites {
		redirect(filepath.Join("to", site.Name), site.URL)
	}

	handlePage("about", pages.AboutPage(translations.language))
	// handlePage("resume", pages.DynamicResume(db, technologies, translations.language))
	redirect("resume", shared.Asset("resume.pdf"))
	redirect("resume.pdf", shared.Asset("resume.pdf"))

	handlePage("contact", pages.ContactPage(false))
	handlePage("contact/sent", pages.ContactPage(true))

	server.HandleFunc("/mail", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		subject := strings.TrimSpace(r.FormValue("subject"))
		body := strings.TrimSpace(r.FormValue("body"))
		from := strings.TrimSpace(r.FormValue("from"))

		err := SendMailToSelf(from, subject, body)
		if err != nil {
			color.Red("Could not send mail: %s", err)
			http.Error(w, "Could not send mail", 500)
			// Write mail to file named with from and today's datetime
			// and content set to subject and body
			nowStr := time.Now().Format("2006-01-02-15-04-05")
			writeTo := filepath.Join("mails", fmt.Sprintf("%s-%s.txt", from, nowStr))
			os.MkdirAll(filepath.Dir(writeTo), 0755)
			os.WriteFile(writeTo, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)), 0644)
		} else {
			http.Redirect(w, r, "/contact/sent", http.StatusSeeOther)
		}
	})

	go http.ListenAndServe(":"+fmt.Sprint(port), server)
	fmt.Printf("[%s] Server started on http://localhost:%d\n", translations.language, port)

	if os.Getenv("ENV") == "static" {
		defer wg.Done()
		fmt.Printf("[%s] Statically rendering %d paths\n", translations.language, len(registeredPaths))
		for _, path := range registeredPaths {
			err := StaticallyRender(filepath.Join("dist", translations.language), fmt.Sprintf("http://localhost:%d", port), path)
			if err != nil {
				color.Red("[%s] Couln't render %s: %s", translations.language, path, err)
				os.Exit(1)
			}
		}
		unusedTranslationsCount := len(translations.UnusedMessages())
		if unusedTranslationsCount > 0 {
			if shared.WantToRemoveUnusedMessages() {
				translations.DeleteUnusedMessages()
				translations.SavePO()
				color.Cyan("[%s] Removed %d unused messages", translations.language, unusedTranslationsCount)
			} else {
				color.Yellow("[%s] %s contains %d unused messages, see %s", translations.language, translations.PoFilePath(), unusedTranslationsCount, translations.UnusedMessagesFilePath())
			}
		}
		translations.WriteUnusedMessages()
		if len(translations.missingMessages) > 0 {
			color.Red("[%s] Some content is not translated. See errors above.", translations.language)
			os.Exit(1)
		}
	}
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
	// Ugly fix while ortfodb isn't fixed
	for i, work := range db {
		if work.ID == "#meta" {
			continue
		}
		if work.Metadata.AdditionalMetadata["made_with"] == nil {
			continue
		}
		newWork := work
		for _, techName := range work.Metadata.AdditionalMetadata["made_with"].([]interface{}) {
			newWork.Metadata.MadeWith = append(newWork.Metadata.MadeWith, techName.(string))
		}
		db[i] = newWork
	}
	locales := db.Languages()
	sort.Strings(locales)
	fmt.Printf("[  ] Works database has %d works in %d locales (%s)\n", len(db.Works()), len(locales), strings.Join(locales, ", "))
	return db, locales
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

	translations, err := LoadTranslations(locales)
	if err != nil {
		color.Yellow("[!!] Couldn't load translations: %s", err)
		os.Exit(1)
		return
	}

	var collections map[string]shared.Collection
	loadDataFileMap("collections.yaml", &collections)
	for id, c := range collections {
		newCollection := c
		newCollection.ID = id
		collections[id] = newCollection
	}

	var sites []shared.Site
	loadDataFile("sites.yaml", &sites)

	var technologies []shared.Technology
	loadDataFile("technologies.yaml", &technologies)

	var tags []shared.Tag
	loadDataFile("tags.yaml", &tags)

	blogEntries, err := loadBlogEntries("blog")
	if err != nil {
		fmt.Printf("while loading blog entries: %s\n", err.Error())
		os.Exit(1)
	}

	if os.Getenv("SKIP_WAKATIME") != "1" {
		for i := range technologies {
			fmt.Printf("[  ] Calculating time spent on %s\n", technologies[i].Name)
			_, err := technologies[i].CalculateTimeSpent(technologies)
			if err != nil {
				color.Yellow("[!!] While calculating time spent on %s: %s", technologies[i].Name, err)
			} else if technologies[i].TimeSpent.Seconds() > 0 {
				fmt.Printf("[  ] Time spent with %s is %s, via wakatime\n", technologies[i].Name, technologies[i].TimeSpent)
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(locales))

	for i, locale := range locales {
		go startPagesServer(&wg, db, collections, sites, technologies, tags, blogEntries, translations[locale], 8081+i)
	}
	go startFileServer(8079, "public")
	go startFileServer(8080, "media")
	var origin string
	if shared.IsDev() {
		origin = "http://localhost:8081"
	} else {
		origin = "https://ewen.works"
	}
	go startSHSServer(&wg, 8666, origin, sites, db, collections, technologies, tags)

	wg.Wait()
}
