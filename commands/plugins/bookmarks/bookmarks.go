package bookmarks

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nickvanw/bogon/commands"
)

const (
	titlePrefix   = "bookmark_"
	commandPrefix = "."
)

// Handler contains the information necessary to store and retrieve bookmarks
// as well as reserved matches for new bookmarks
type Handler struct {
	datastore Storage
	reserved  map[string]*regexp.Regexp
	tpl       *template.Template
}

// New returns a new bookmark handler, if a nil Storage is passed it will
// store them in memory
func New(d Storage) (*Handler, error) {
	if d == nil {
		d = NewMemStorage()
	}
	t, err := template.New("t").Parse(tpl)
	if err != nil {
		return nil, err
	}
	return &Handler{
		datastore: d,
		tpl:       t,
		reserved:  map[string]*regexp.Regexp{},
	}, nil
}

// Block records regular expressions that will not be allowed
// for bookmark names
func (h *Handler) Block(c map[string]*regexp.Regexp) {
	for k, v := range c {
		h.reserved[k] = v
	}
}

// Exports returns a list of plugins to export for use in an IRC bot
func (h *Handler) Exports() []commands.RegisterFunc {
	rawHandler := func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
		return titlePrefix + "raw", nil, h.rawHandler, commands.Options{Raw: true}
	}
	newHandler := func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
		r := regexp.MustCompile("(?i)^\\.b(ook)?m(ark)?$")
		return titlePrefix + "new", r, h.newHandler, commands.Options{}
	}
	delHandler := func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
		r := regexp.MustCompile("(?i)^\\.delbm$")
		return titlePrefix + "del", r, h.delHandler, commands.Options{}
	}
	return []commands.RegisterFunc{rawHandler, newHandler, delHandler}
}

// Serve is a blocking method that creates and starts a web server at the
// specified address that serves the entire list of bookmarks being tracked
func (h *Handler) Serve(addr string) error {
	r := mux.NewRouter()
	r.HandleFunc("/", h.webHandler)
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}
	return server.ListenAndServe()
}

func (h *Handler) webHandler(w http.ResponseWriter, r *http.Request) {
	data, err := h.datastore.Dump()
	if err != nil {
		http.Error(w, "Unable to get a bookmark list", 500)
	}
	h.tpl.Execute(w, data)
}

func (h *Handler) rawHandler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) == 0 || msg.Params[0][0] != commandPrefix[0] {
		return ""
	}
	key := strings.ToLower(msg.Params[0][1:])
	if data, ok := h.datastore.Lookup(key); ok {
		return fmt.Sprintf("[%s]: %s", key, data)
	}
	return ""
}

func (h *Handler) newHandler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 3 || len(msg.Params[1]) < 2 || h.isBanned(msg.Params[1]) {
		return "Usage: .bm [key] [bookmark message], where key is 3 or more letters and not a command"

	}
	messsage := strings.Join(msg.Params[2:], " ")
	key := strings.ToLower(msg.Params[1])
	if err := h.datastore.New(key, messsage); err != nil {
		return "Sorry, I hit an error inserting that"

	}
	return fmt.Sprintf("Successfully added %q", msg.Params[1])
}

func (h *Handler) delHandler(msg commands.Message, ret commands.MessageFunc) string {
	if len(msg.Params) < 2 {
		return "Usage: .delbm [key]"

	}
	key := strings.ToLower(msg.Params[1])
	if err := h.datastore.Remove(key); err != nil {
		return "I was unable to remove that"

	}
	return fmt.Sprintf("Bookmark %q removed", msg.Params[1])
}

func (h *Handler) isBanned(key string) bool {
	key = commandPrefix + key
	for _, v := range h.reserved {
		if v != nil && v.MatchString(strings.ToLower(key)) {
			return true
		}
	}
	return false
}
