package render

import (
	"fmt"
	"go/token"
	"go/types"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		"default": func(i ...interface{}) (o interface{}) {
			switch {
			case len(i) == 1:
				o = i[0]
			case len(i) >= 2:
				d, x := i[0], i[1]

				v := reflect.ValueOf(x)
				if !v.IsValid() {
					o = d
					break
				}

				var exist bool
				switch v.Kind() {
				case reflect.Bool:
					exist = true
				case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
					exist = v.Len() != 0
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					exist = v.Int() != 0
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					exist = v.Uint() != 0
				case reflect.Float32, reflect.Float64:
					exist = v.Float() != 0
				case reflect.Complex64, reflect.Complex128:
					exist = v.Complex() != 0
				case reflect.Struct:
					switch actual := x.(type) {
					case time.Time:
						exist = !actual.IsZero()
					default:
						exist = true
					}
				default:
					exist = !v.IsNil()
				}

				if exist {
					o = x
				} else {
					o = d
				}
			default:
				o = ""
			}
			return
		},
		"dict": func(i ...interface{}) (o map[interface{}]interface{}) {
			if len(i)%2 == 0 {
				o = make(map[interface{}]interface{})

				for n := 0; n < len(i); n += 2 {
					v0 := reflect.ValueOf(i[n])
					switch v0.Type().Kind() {
					case reflect.Slice, reflect.Array:
						for m, l := 0, v0.Len(); m < l; m += 2 {
							o[v0.Index(m).Interface()] = v0.Index(m + 1).Interface()
						}
					case reflect.Map:
						for _, k := range v0.MapKeys() {
							o[k.Interface()] = v0.MapIndex(k).Interface()
						}
					default:
						o[i[n]] = i[n+1]
					}
				}
			}
			return
		},
		"eval": func(i interface{}) (o string) {
			if v, err := toString(i); err == nil {
				if result, err := types.Eval(token.NewFileSet(), nil, token.NoPos, v); err == nil {
					o = result.Value.ExactString()
				} else {
					fmt.Println(err.Error())
				}
			}
			return
		},
		"findRE": func(i ...interface{}) (o []string) {
			if size := len(i); size == 2 {
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				if err0 == nil && err1 == nil {
					o = regexp.MustCompile(v0).FindAllString(v1, -1)
				}
			} else if size == 3 {
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				v2, err2 := toString(i[2])
				if err0 == nil && err1 == nil && err2 == nil {
					if limit, err := strconv.Atoi(v2); err == nil {
						o = regexp.MustCompile(v0).FindAllString(v1, limit)
					} else {
						fmt.Println(err.Error())
					}
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
		"index": func(i ...interface{}) (o interface{}) {
			if len(i) >= 2 {
				// index CORRECTION (INDEX|KEY)
				cv := reflect.ValueOf(i[0]) // CORRECTION
				iv := reflect.ValueOf(i[1]) // (INDEX|KEY)

				switch cv.Type().Kind() {
				case reflect.Slice, reflect.Array:
					o = cv.Index(int(iv.Int()))
				case reflect.Map:
					o = cv.MapIndex(iv)
				}
			}
			return
		},
		"len": func(i interface{}) (o int) {
			switch reflect.TypeOf(i).Kind() {
			case reflect.Array, reflect.Map, reflect.Slice:
				o = reflect.ValueOf(i).Len()
			default:
				if s, e := toString(i); e == nil {
					o = len(s)
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
			if v, err := toString(i); err == nil {
				o = template.CSS(v)
			}
			return
		},
		"safeHTML": func(i interface{}) (o template.HTML) {
			if v, err := toString(i); err == nil {
				o = template.HTML(v)
			}
			return
		},
		"safeHTMLAttr": func(i interface{}) (o template.HTMLAttr) {
			if v, err := toString(i); err == nil {
				o = template.HTMLAttr(v)
			}
			return
		},
		"safeJS": func(i interface{}) (o template.JS) {
			if v, err := toString(i); err == nil {
				o = template.JS(v)
			}
			return
		},
		"safeURL": func(i interface{}) (o template.URL) {
			if v, err := toString(i); err == nil {
				o = template.URL(v)
			}
			return
		},
		"slice": func(i ...interface{}) (o []interface{}) {
			return i
		},
		"split": func(i ...string) (o []string) {
			if len(i) >= 2 {
				// split STRING DELIMITER
				s, errS := toString(i[0]) // STRING
				d, errD := toString(i[1]) // DELIMITER

				if errS == nil && errD == nil {
					o = strings.Split(s, d)
				}
			}
			return
		},
		"time": func() (o string) {
			return time.Now().UTC().Format("15:04:05")
		},
		"trim": func(i ...interface{}) (o string) {
			switch len(i) {
			case 0:
				// none
			case 1:
				if v, err := toString(i[0]); err == nil {
					o = strings.Trim(v, " ")
				}
			default:
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				if err0 == nil && err1 == nil {
					o = strings.Trim(v0, v1)
				}
			}
			return
		},
		"trimLeft": func(i ...interface{}) (o string) {
			switch len(i) {
			case 0:
				// none
			case 1:
				if v, err := toString(i[0]); err == nil {
					o = strings.TrimLeft(v, " ")
				}
			default:
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				if err0 == nil && err1 == nil {
					o = strings.TrimLeft(v0, v1)
				}
			}
			return
		},
		"trimRight": func(i ...interface{}) (o string) {
			switch len(i) {
			case 0:
				// none
			case 1:
				if v, err := toString(i[0]); err == nil {
					o = strings.TrimRight(v, " ")
				}
			default:
				v0, err0 := toString(i[0])
				v1, err1 := toString(i[1])
				if err0 == nil && err1 == nil {
					o = strings.TrimRight(v0, v1)
				}
			}
			return
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
			hw.Header().Set("Content-Type", "text/html")
			hw.WriteHeader(status)
		}
	}
}

type Template struct {
	OsFiles   []*os.File
	HttpFiles []http.File
	Layout    string
	Data      interface{}
}

func (t Template) HTML(w io.Writer, status int) {
	tmpl := template.New(t.Layout).Funcs(funcMap)

	for i := range t.OsFiles {
		buf, _ := ioutil.ReadAll(t.OsFiles[i])
		if t, err := tmpl.Parse(string(buf)); err == nil {
			tmpl = t
		}
	}

	for i := range t.HttpFiles {
		buf, _ := ioutil.ReadAll(t.HttpFiles[i])
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
