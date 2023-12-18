package shared

type LocalizableString map[string]string

func (ls LocalizableString) Localized(locale string) string {
	return ls[locale]
}

type Site struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Purpose  string `yaml:"purpose"`
	Username string `yaml:"username"`
}

type Technology struct {
	Slug         string   `yaml:"slug"`
	Name         string   `yaml:"name"`
	By           string   `yaml:"by"`
	Files        []string `yaml:"files"`
	Aliases      []string `yaml:"aliases"`
	LearnMoreURL string   `yaml:"learn more at"`
	Description  string   `yaml:"description"`
}

type Tag struct {
	Singular     string   `yaml:"singular"`
	Plural       string   `yaml:"plural"`
	Aliases      []string `yaml:"aliases"`
	Description  string   `yaml:"description"`
	LearnMoreURL string   `yaml:"learn more at"`
}

type Collection struct {
	Title        map[string]string `yaml:"title"`
	Aliases      []string          `yaml:"aliases"`
	Includes     string            `yaml:"includes"`
	Description  map[string]string `yaml:"description"`
	LearnMoreURL string            `yaml:"learn more at"`
}

type Collections map[string]Collection

func (cs Collections) URLsToNames(canonicalOnly bool, locale string) map[string]string {
	urlsToNames := make(map[string]string)
	for id, collection := range cs {
		if !canonicalOnly {
			for _, alias := range collection.Aliases {
				urlsToNames[alias] = collection.Title[locale]
			}
		}
		urlsToNames[id] = collection.Title[locale]
	}
	return urlsToNames
}
