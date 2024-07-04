package render

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

var funcMap template.FuncMap

func toString(i any) (string, error) {
	switch v := i.(type) {
	case string:
		return v, nil
	case bool:
		return strconv.FormatBool(v), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case uint:
		return strconv.FormatInt(int64(v), 10), nil
	case uint64:
		return strconv.FormatInt(int64(v), 10), nil
	case uint32:
		return strconv.FormatInt(int64(v), 10), nil
	case uint16:
		return strconv.FormatInt(int64(v), 10), nil
	case uint8:
		return strconv.FormatInt(int64(v), 10), nil
	case []byte:
		return string(v), nil
	case template.CSS:
		return string(v), nil
	case template.HTML:
		return string(v), nil
	case template.HTMLAttr:
		return string(v), nil
	case template.JS:
		return string(v), nil
	case template.JSStr:
		return string(v), nil
	case template.URL:
		return string(v), nil
	case template.Srcset:
		return string(v), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	default:
		return "", fmt.Errorf("cast error; value: %#v, type: %T", i, i)
	}
}

func init() {
	funcMap = template.FuncMap{
		"day":          dayFn,
		"date":         dateFn,
		"datetime":     datetimeFn,
		"default":      defaultFn,
		"dict":         dictFn,
		"eval":         evalFn,
		"findRE":       findREFn,
		"in":           inFn,
		"index":        indexFn,
		"len":          lenFn,
		"lower":        lowerFn,
		"map":          mapFn,
		"month":        monthFn,
		"replace":      replaceFn,
		"replaceRE":    replaceREFn,
		"safeCSS":      safeCSSFn,
		"safeHTML":     safeHTMLFn,
		"safeHTMLAttr": safeHTMLAttrFn,
		"safeJS":       safeJSFn,
		"safeURL":      safeURLFn,
		"slice":        sliceFn,
		"split":        splitFn,
		"time":         timeFn,
		"trim":         trimFn,
		"trimLeft":     trimLeftFn,
		"trimRight":    trimRightFn,
		"upper":        upperFn,
		"year":         yearFn,
	}
}

// HTML renders the template as HTML to the provided io.Writer.
//
// Deprecated: Use Template.HTML method instead.
func HTML(w io.Writer, status int, data any, layout string, ext ...string) {
	files := append([]string{layout}, ext...)
	for i := range files {
		files[i] = "templates/" + files[i] + ".html"
	}

	tmpl := template.Must(template.New(layout + ".html").Funcs(funcMap).ParseFiles(files...))
	if hw, ok := w.(http.ResponseWriter); ok {
		if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
			http.Error(hw, err.Error(), http.StatusInternalServerError)
		} else {
			hw.Header().Set("Content-Type", "text/html")
			hw.WriteHeader(status)
		}
	}
}

type Template struct {
	OsFiles   []*os.File
	HttpFiles []http.File
	Layout    string
	Data      any
}

// HTML renders the template as HTML to the provided io.Writer.
//
// This method takes an io.Writer (typically an http.ResponseWriter) and an HTTP status code.
// It renders the template with the given data and writes the resulting HTML to the writer.
// The HTTP status code is set on the http.ResponseWriter to indicate the response status.
//
// Parameters:
// - w: io.Writer to which the rendered HTML will be written. This is usually an http.ResponseWriter.
// - status: HTTP status code to set on the http.ResponseWriter.
//
// Example usage:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//	    tmpl := Template{
//	        Layout: "layout.html",
//	        Data:   myData,
//	    }
//	    tmpl.HTML(w, http.StatusOK)
//	}
func (t Template) HTML(w io.Writer, status int) {
	tmpl := template.New(t.Layout).Funcs(funcMap)

	for i := range t.OsFiles {
		buf, _ := io.ReadAll(t.OsFiles[i])
		if t, err := tmpl.Parse(string(buf)); err == nil {
			tmpl = t
		}
	}

	for i := range t.HttpFiles {
		buf, _ := io.ReadAll(t.HttpFiles[i])
		if t, err := tmpl.Parse(string(buf)); err == nil {
			tmpl = t
		}
	}

	if hw, ok := w.(http.ResponseWriter); ok {
		if err := tmpl.ExecuteTemplate(w, t.Layout, t.Data); err != nil {
			http.Error(hw, err.Error(), http.StatusInternalServerError)
		} else {
			hw.Header().Set("Content-Type", "text/html")
			hw.WriteHeader(status)
		}
	}
}
