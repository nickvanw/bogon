package bogon

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof" // imported for side effects
	"os"

	"github.com/gorilla/mux"
)

func (c *Client) adminSocket(fn string) error {
	_ = os.Remove(fn)
	l, err := net.ListenUnix("unix", &net.UnixAddr{fn, "unix"})
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	r.HandleFunc("/state", c.dumpState)
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	return http.Serve(l, r)
}

func (c *Client) dumpState(w http.ResponseWriter, r *http.Request) {
	data := bytes.NewBuffer(nil)
	if err := json.NewEncoder(data).Encode(c.state); err != nil {
		http.Error(w, "unable to serialize state", 500)
		return
	}
	io.Copy(w, data)
}
