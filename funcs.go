package render

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"html/template"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func dayFn() (o int) {
	return time.Now().UTC().Day()
}

func dateFn() (o string) {
	return time.Now().UTC().Format("2006-01-02")
}

func datetimeFn() (o string) {
	return time.Now().UTC().Format(time.RFC3339)
}

func defaultFn(i ...any) (o any) {
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
}

func dictFn(i ...any) (o map[any]any) {
	if len(i)%2 == 0 {
		o = make(map[any]any)

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
}

func evalFn(i any) (o any) {
	if v, err := toString(i); err == nil {
		if result, err := types.Eval(token.NewFileSet(), nil, token.NoPos, v); err == nil {
			switch result.Value.Kind() {
			case constant.Bool:
				return constant.BoolVal(result.Value)
			case constant.String:
				return constant.StringVal(result.Value)
			case constant.Int:
				if val, ok := constant.Int64Val(result.Value); ok {
					return val
				} else if val, ok := constant.Uint64Val(result.Value); ok {
					return val
				}
				return nil
			case constant.Float:
				val, _ := constant.Float64Val(result.Value)
				return val
			case constant.Complex:
				re, _ := constant.Float64Val(constant.Real(result.Value))
				im, _ := constant.Float64Val(constant.Imag(result.Value))
				return complex(re, im)
			default:
				return nil
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return
}

func findREFn(i ...any) (o []string) {
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
}

func inFn(i ...any) (o bool) {
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
}

func indexFn(i ...any) (o any) {
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
}

func lenFn(i any) (o int) {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		o = reflect.ValueOf(i).Len()
	default:
		if s, e := toString(i); e == nil {
			o = len(s)
		}
	}
	return
}

func lowerFn(i ...any) (o string) {
	var buf []string
	for _, tmp := range i {
		if s, err := toString(tmp); err == nil {
			buf = append(buf, strings.ToLower(s))
		}
	}
	return strings.Join(buf, " ")
}

func mapFn(i ...any) (o map[string]any) {
	if size := len(i); size%2 == 0 && size > 0 {
		o = make(map[string]any)
		var s string
		var e error
		for n := 0; size/2 > n; n++ {
			if s, e = toString(i[n*2]); e == nil {
				o[s] = i[n*2+1]
			}
		}
	}
	return
}

func monthFn() (o int) {
	return int(time.Now().UTC().Month())
}

func replaceFn(i ...any) (o string) {
	if len(i) == 3 {
		v0, err0 := toString(i[0])
		v1, err1 := toString(i[1])
		v2, err2 := toString(i[2])
		if err0 == nil && err1 == nil && err2 == nil {
			o = strings.Replace(v0, v1, v2, -1)
		}
	}
	return
}

func replaceREFn(i ...any) (o string) {
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
}

func safeCSSFn(i any) (o template.CSS) {
	if v, err := toString(i); err == nil {
		o = template.CSS(v)
	}
	return
}

func safeHTMLFn(i any) (o template.HTML) {
	if v, err := toString(i); err == nil {
		o = template.HTML(v)
	}
	return
}

func safeHTMLAttrFn(i any) (o template.HTMLAttr) {
	if v, err := toString(i); err == nil {
		o = template.HTMLAttr(v)
	}
	return
}

func safeJSFn(i any) (o template.JS) {
	if v, err := toString(i); err == nil {
		o = template.JS(v)
	}
	return
}

func safeURLFn(i any) (o template.URL) {
	if v, err := toString(i); err == nil {
		o = template.URL(v)
	}
	return
}

func sliceFn(i ...any) (o []any) {
	return i
}

func splitFn(i ...string) (o []string) {
	if len(i) >= 2 {
		// split STRING DELIMITER
		s, errS := toString(i[0]) // STRING
		d, errD := toString(i[1]) // DELIMITER

		if errS == nil && errD == nil {
			o = strings.Split(s, d)
		}
	}
	return
}

func timeFn() (o string) {
	return time.Now().UTC().Format("15:04:05")
}

func trimFn(i ...any) (o string) {
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
}

func trimLeftFn(i ...any) (o string) {
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
}

func trimRightFn(i ...any) (o string) {
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
}

func upperFn(i ...any) (o string) {
	var buf []string
	for _, tmp := range i {
		if s, err := toString(tmp); err == nil {
			buf = append(buf, strings.ToUpper(s))
		}
	}
	return strings.Join(buf, " ")
}

func yearFn() (o int) {
	return time.Now().UTC().Year()
}
