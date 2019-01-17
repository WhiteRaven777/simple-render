package render

import (
	"html/template"
	"io"
	"net/http"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{
		"safeCSS": func(i interface{}) (o template.CSS) {
			if v, ok := i.(string); ok {
				o = template.CSS(v)
			}
			return
		},
		"safeHTML": func(i interface{}) (o template.HTML) {
			if v, ok := i.(string); ok {
				o = template.HTML(v)
			}
			return
		},
		"safeHTMLAttr": func(i interface{}) (o template.HTMLAttr) {
			if v, ok := i.(string); ok {
				o = template.HTMLAttr(v)
			}
			return
		},
		"safeJS": func(i interface{}) (o template.JS) {
			if v, ok := i.(string); ok {
				o = template.JS(v)
			}
			return
		},
		"safeURL": func(i interface{}) (o template.URL) {
			if v, ok := i.(string); ok {
				o = template.URL(v)
			}
			return
		},
		"slice": func(i ...interface{}) (o []interface{}) {
			return i
		},
	}
}

func HTML(w io.Writer, status int, data interface{}, layout string, ext ...string) {
	files := append([]string{layout}, ext...)
	for i := range files {
		files[i] = "templates/" + files[i] + ".html"
	}

	tmpl := template.Must(template.New(layout + ".html").Funcs(funcMap).ParseFiles(files...))
	if hw, ok := w.(http.ResponseWriter); ok {
		if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
			http.Error(hw, err.Error(), http.StatusInternalServerError)
		} else {
			hw.Header().Set("Context-Type", "text/html")
			hw.WriteHeader(status)
		}
	}
}
