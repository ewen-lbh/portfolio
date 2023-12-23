package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/chai2010/gettext-go/po"
	mapset "github.com/deckarep/golang-set"
	"github.com/ewen-lbh/portfolio/shared"
	"github.com/fatih/color"
	"golang.org/x/net/html"
)

const TRANSLATABLE_MARKER_ATTRIBUTE = "i18n"

var SourceLanguage = "en"

// Translations holds both the gettext catalog from the .mo file
// and a po file object used to update the .po file (e.g. when discovering new translatable strings)
type Translations struct {
	poFile          po.File
	seenMessages    mapset.Set
	missingMessages []po.Message
	language        string
}

func (t Translations) PoFilePath() string {
	return filepath.Join("i18n", fmt.Sprintf("%s.po", t.language))
}

type TranslationsCatalog map[string]*Translations

type HttpTranslator struct {
	translations *Translations
	ch           *templ.ComponentHandler
}

func (h HttpTranslator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if shared.IsDev() {
		color.Magenta("[%s] Reloading translations since we're in dev mode", h.translations.language)
		newCatalog, err := LoadTranslations([]string{h.translations.language})
		if err != nil {
			color.Yellow("[%s] Failed to reload translations: %s", h.translations.language, err)
		} else {
			*h.translations = *newCatalog[h.translations.language]
		}
	}
	w.Header().Set("Language", h.translations.language)
	if h.ch.Status != 0 {
		w.WriteHeader(h.ch.Status)
	}
	w.Header().Add("Content-Type", h.ch.ContentType)
	var untranslated bytes.Buffer
	err := h.ch.Component.Render(r.Context(), &untranslated)
	if err != nil {
		if h.ch.ErrorHandler != nil {
			h.ch.ErrorHandler(r, err).ServeHTTP(w, r)
			return
		}
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}

	translated, err := h.translations.TranslatePage(untranslated.Bytes())
	// fmt.Printf("[%s] Translated %q\n", h.translations.language, translated)
	if err != nil {
		color.Red("[%s] %s: Error while translating: %w", h.translations.language, r.URL.Path, err)
		return
	}

	_, err = io.WriteString(w, translated)
	if err != nil {
		color.Red("[%s] %s: Error while writing: %s", h.translations.language, r.URL.Path, err)
		return
	}

	h.translations.WriteUnusedMessages()
	h.translations.SavePO()
}

func (t *Translations) TranslatePage(source []byte) (string, error) {
	parsed, err := html.Parse(bytes.NewReader(source))
	if err != nil {
		return "", fmt.Errorf("while parsing output page HTML: %w", err)
	}

	return t.Translate(parsed), nil
}

// TranslateToLanguage translates the given html node to french or english, removing translation-related attributes
func (t *Translations) Translate(root *html.Node) string {
	// Open files
	doc := goquery.NewDocumentFromNode(root)
	doc.Find("i18n, [i18n]").Each(func(_ int, element *goquery.Selection) {
		element.RemoveAttr("i18n")
		msgContext, _ := element.Attr("i18n-context")
		element.RemoveAttr("i18n-context")
		if t.language != SourceLanguage {
			innerHTML, _ := element.Html()
			innerHTML = html.UnescapeString(innerHTML)
			innerHTML = strings.TrimSpace(innerHTML)
			if innerHTML == "" {
				return
			}
			translated, err := t.GetTranslation(innerHTML, msgContext)
			if err != nil {
				color.Yellow("[%s] Missing translation for %q", t.language, innerHTML)
				t.missingMessages = append(t.missingMessages, po.Message{
					MsgId:      innerHTML,
					MsgContext: msgContext,
				})
			} else {
				element.SetHtml(translated)
			}
		}
	})
	htmlString, _ := doc.Html()
	htmlString = strings.ReplaceAll(htmlString, "<i18n>", "")
	htmlString = strings.ReplaceAll(htmlString, "</i18n>", "")
	return htmlString
}

// LoadTranslations reads from i18n/[language].po to load translations
func LoadTranslations(languages []string) (TranslationsCatalog, error) {
	translations := make(TranslationsCatalog)
	for _, languageCode := range languages {
		translationsFilepath := fmt.Sprintf("i18n/%s.po", languageCode)
		poFile, err := po.LoadFile(translationsFilepath)
		if err != nil {
			color.Yellow("[%s] Couldn't load translations: %s", languageCode, err)
			err = WriteEmptyPOFile(languageCode)
			if err != nil {
				return nil, fmt.Errorf("while writing empty PO file: %w", err)
			}

			return LoadTranslations(languages)
		} else {
			translations[languageCode] = &Translations{
				poFile:          *poFile,
				seenMessages:    mapset.NewSet(),
				missingMessages: make([]po.Message, 0),
				language:        languageCode,
			}
			fmt.Printf("[%s] Loaded %d translations\n", languageCode, len(poFile.Messages))
		}
	}
	return translations, nil
}

func WriteEmptyPOFile(language string) error {
	poFile := po.File{
		Messages: []po.Message{},
		MimeHeader: po.Header{
			Language: language,
		},
	}

	t := Translations{
		language: language,
	}

	color.Cyan("[  ] Writing empty PO file for %s at %s", t.language, t.PoFilePath())
	os.MkdirAll(filepath.Dir(t.PoFilePath()), 0755)
	return poFile.Save(t.PoFilePath())
}

func (t Translations) WriteUnusedMessages() error {
	unused := make([]po.Message, 0)
	to := fmt.Sprintf("i18n/%s-unused-messages.yaml", t.language)
	for _, message := range t.poFile.Messages {
		if !t.seenMessages.Contains(message.MsgId + message.MsgContext) {
			unused = append(unused, message)
		}
	}

	if len(unused) == 0 {
		os.Remove(to)
	} else {
		file, err := os.OpenFile(to, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		os.WriteFile(to, []byte("# Generated at "+time.Now().String()+"\n"), 0644)
		for _, message := range unused {

			if message.MsgContext != "" {
				_, err = file.WriteString(fmt.Sprintf("- {msgid: %q, msgctxt: %q}\n", message.MsgId, message.MsgContext))
			} else {
				_, err = file.WriteString(fmt.Sprintf("- %q\n", message.MsgId))
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SavePO writes the .po file to the disk, with its potential modifications
// It removes duplicate messages beforehand
func (t *Translations) SavePO() {
	// TODO: sort file after saving, (po.File).Save is not stable... (creates unecessary diffs in git)
	// Remove unused messages with empty msgstrs
	uselessRemoved := make([]po.Message, 0)
	for _, msg := range t.poFile.Messages {
		if !t.seenMessages.Contains(msg.MsgId+msg.MsgContext) && msg.MsgStr == "" {
			t.seenMessages.Remove(msg.MsgId + msg.MsgContext)
			continue
		}
		uselessRemoved = append(uselessRemoved, msg)
	}
	t.poFile.Messages = uselessRemoved
	// Add missing messages
	t.poFile.Messages = append(t.poFile.Messages, t.missingMessages...)
	// Remove duplicate messages
	dedupedMessages := make([]po.Message, 0)
	for _, msg := range t.poFile.Messages {
		var isDupe bool
		for _, msg2 := range dedupedMessages {
			if msg.MsgId == msg2.MsgId && msg.MsgContext == msg2.MsgContext {
				isDupe = true
			}
		}
		if !isDupe {
			dedupedMessages = append(dedupedMessages, msg)
		}
	}
	t.poFile.Messages = dedupedMessages
	// Sort them to guarantee a stable write
	sort.Sort(ByMsgIdAndCtx(t.poFile.Messages))
	t.poFile.Save(t.PoFilePath())
}

// ByMsgIdAndCtx implement sorting gettext messages by their msgid+msgctxt
type ByMsgIdAndCtx []po.Message

func (b ByMsgIdAndCtx) Len() int {
	return len(b)
}

func (b ByMsgIdAndCtx) Less(i, j int) bool {
	return b[i].MsgId < b[j].MsgId || (b[i].MsgId == b[j].MsgId && b[i].MsgContext < b[j].MsgContext)
}

func (b ByMsgIdAndCtx) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// GetTranslation returns the msgstr corresponding to msgid and msgctxt from the .po file
// If not found, it returns an error
func (t Translations) GetTranslation(msgid string, msgctxt string) (string, error) {
	t.seenMessages.Add(msgid + msgctxt)
	for _, message := range t.poFile.Messages {
		if message.MsgId == msgid && message.MsgStr != "" && message.MsgContext == msgctxt {
			return message.MsgStr, nil
		}
	}
	return "", fmt.Errorf("cannot find msgstr in %s with msgid=%q and msgctx=%q", t.language, msgid, msgctxt)
}
