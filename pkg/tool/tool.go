package tool

import (
	"crypto/tls"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	"github.com/gorilla/mux"

	"github.com/charlieegan3/toolbelt/pkg/apis"
)

type Echo struct {
	path string
}

func (s *Echo) Name() string {
	return "echo"
}

func (s *Echo) FeatureSet() apis.FeatureSet {
	return apis.FeatureSet{
		HTTP:   true,
		Config: true,
	}
}

func (s *Echo) SetConfig(config map[string]any) error {

	var ok bool

	cfg := gabs.Wrap(config)

	path := "path"
	s.path, ok = cfg.Path(path).Data().(string)
	if !ok {
		return fmt.Errorf("config value %s not set", path)
	}

	return nil
}

func (s *Echo) Jobs() ([]apis.Job, error) { return []apis.Job{}, nil }

func (s *Echo) HTTPAttach(router *mux.Router) error {

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jsonHeaders := make(map[string][]string)
		for k, v := range r.Header {
			jsonHeaders[k] = v
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")

		err := enc.Encode(struct {
			Method           string               `json:"method"`
			Host             string               `json:"host"`
			Proto            string               `json:"proto"`
			URL              string               `json:"url"`
			RequestURI       string               `json:"request_uri"`
			RemoteAddr       string               `json:"remote_addr"`
			ContentLength    int64                `json:"content_length"`
			TransferEncoding []string             `json:"transfer_encoding"`
			TLS              *tls.ConnectionState `json:"tls"`
			Headers          map[string][]string  `json:"headers"`
		}{
			Method:        r.Method,
			Host:          r.Host,
			Proto:         r.Proto,
			URL:           r.URL.String(),
			ContentLength: r.ContentLength,
			RemoteAddr:    r.RemoteAddr,
			RequestURI:    r.RequestURI,
			TLS:           r.TLS,
			Headers:       jsonHeaders,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return nil
}
func (s *Echo) HTTPHost() string {
	return ""
}
func (s *Echo) HTTPPath() string { return s.path }

func (s *Echo) ExternalJobsFuncSet(f func(job apis.ExternalJob) error) {}

func (s *Echo) DatabaseSet(db *sql.DB) {}

func (s *Echo) DatabaseMigrations() (*embed.FS, string, error) {
	return nil, "", nil
}
