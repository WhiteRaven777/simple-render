package render

import (
	"fmt"
	"go/token"
	"go/types"
	"html/template"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var funcMap template.FuncMap

func toString(i interface{}) (string, error) {
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
	case template.HTML:
		return string(v), nil
	case template.URL:
		return string(v), nil
	case template.JS:
		return string(v), nil
	case template.CSS:
		return string(v), nil
	case template.HTMLAttr:
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
		"day": func() (o int) {
			return time.Now().UTC().Day()
		},
		"date": func() (o string) {
			return time.Now().UTC().Format("2006-01-02")
		},
		"datetime": func() (o string) {
			return time.Now().UTC().Format(time.RFC3339)
		},
		"eval": func(i interface{}) (o string) {
			if v, ok := i.(string); ok {
				if result, err := types.Eval(token.NewFileSet(), nil, token.NoPos, v); err == nil {
					o = result.Value.ExactString()
				} else {
					fmt.Println(err.Error())
				}
			}
			return
		},
		"in": func(i ...interface{}) (o bool) {
			if len(i) == 2 {
				if i[0] != nil {
					v := reflect.ValueOf(i[0])
					switch v.Type().Kind() {
					case reflect.Slice:
						var v0 string
						var err0 error
						v1, err1 := toString(i[1])
						if err1 == nil {
							for n := 0; n < v.Len(); n++ {
								v0, err0 = toString(v.Index(n))
								if err0 == nil && !o {
									o = strings.Contains(v0, v1)
									return
								}
							}
						}
					default:
						v0, err0 := toString(i[0])
						v1, err1 := toString(i[1])
						if err0 == nil && err1 == nil {
							o = strings.Contains(v0, v1)
						}
					}
				}
			}
			return
		},
		"map": func(i ...interface{}) (o map[string]interface{}) {
			if size := len(i); size%2 == 0 && size > 0 {
				o = make(map[string]interface{})
				var s string
				var e error
				for n := 0; size/2 > n; n++ {
					if s, e = toString(i[n*2]); e == nil {
						o[s] = i[n*2+1]
					}
				}
			}
			return
		},
		"month": func() (o int) {
			return int(time.Now().UTC().Month())
		},
		"replace": func(i ...interface{}) (o string) {
			if len(i) == 3 {
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				v2, err2 := toString(i[2])
				if err0 == nil && err1 == nil && err2 == nil {
					o = strings.Replace(v0, v1, v2, -1)
				}
			}
			return
		},
		"replaceRE": func(i ...interface{}) (o string) {
			if len(i) == 3 {
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				v2, err2 := toString(i[2])
				if err0 == nil && err1 == nil && err2 == nil {
					r := regexp.MustCompile(v0)
					o = r.ReplaceAllString(v2, v1)
				}
			}
			return
		},
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
		"time": func() (o string) {
			return time.Now().UTC().Format("15:04:05")
		},
		"year": func() (o int) {
			return time.Now().UTC().Year()
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
