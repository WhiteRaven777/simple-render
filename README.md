# simple-render

This is a simple render.

Currently the following formats are supported.
* HTML
* JSON

# How To Use

## JSON
```
type Data Struct {
    Date string
    Msg  string
}

data := Data {
    Date: time.Now().Formt("2006-01-02"),
    Msg:  "message",
}

--- or ---

data := map[string]string{
    "Date": time.Now().Formt("2006-01-02"),
    "Msg":  "message",
}

---

func Sample(w http.ResponseWriter, r *http.Request) {
    render.JSON(w, http.StatusOK, data)
}

    |
    V
{
  "Date": "2019-01-17",
  "Msg": "message"
}
```

## HTML
### layout.html
```
{{ define "layout" }}<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{ template "title" . }}</title>
</head>
<body>
<div class="content">
{{ template "body" . }}
</div>
</body>
</html>{{ end }}
```

### view.html
```
{{ define "title" }}{{ .Title }}{{ end }}

{{ define "body" }}
<h1>{{ .Title }}</h1>
<article>
{{ .Body | safeHTML }}
</article>
<li><a href="{{ .Url | safeURL }}">{{ .Title }}</a></li>
{{ end }}
```

```
func Sample(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{
        "Title": "title text",
        "Body":  "body text",
        "Url":   "https://github.com/WhiteRaven777/simple-render",
    }
    render.HTML(w, http.StatusOK, data, "base", "view")
}
```

### Functions

#### ● day
##### Syntax
```
day
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ day }}
-> 14
```

#### ● date
##### Syntax
```
date
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ date }}
-> 2019-02-14
```

#### ● datetime
##### Syntax
```
datetime
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ datetime }}
-> 2019-02-14T17:18:19Z

# format: RFC3339
```

#### ● default
##### Syntax
```
default DEFAULT_VALUE INPUT
```
##### Example
```
{{ 1 | default 2}}
-> 2

{{ "" | default "default"}}
-> default
```

#### ● dict
##### Syntax
```
dict [KEY VALUE]...
```
##### Example
```
{{ dict 0 "aaa" 1 "bbb" 2 "ccc" }}
-> map[0:"aaa",1:"bbb",2:"ccc"]
```

#### ● eval
##### Syntax
```
eval FORMULAS
```
##### Example
```
{{ eval "1 + 1" }}
-> 2

{{ eval "1 - 1 + 1" }}
-> 1
```

#### ● findRE
##### Syntax
```
findRE PATTERN INPUT [LIMIT]
```
##### Example
```
{{ findRE `/hoge/(\d+)` "/hoge/1234567890/hoge/987654321/abcdefghijk"   }}
-> [/hoge/1234567890 /hoge/987654321]

{{ findRE `/hoge/(\d+)` "/hoge/1234567890/hoge/987654321/abcdefghijk" 1 }}
-> [/hoge/1234567890]
```

#### ● in
##### Syntax
```
in SET ITEM
```
##### Example
```
{{ if in "/example/aaa" "/example" }}True{{ else }}False{{ end }}
-> True

{{ if in "/sample/aaa" "/example" }}True{{ else }}False{{ end }}
-> False
```

#### ● index
##### Syntax
```
index CORRECTION (INDEX|KEY)
```
##### Example
```
{{ index (slice 3 4 5) 0 }}
-> 3

{{ index (dict 0 "aaa" 1 "bbb" 2 "ccc") 0 }}
-> aaa
```

#### ● len
##### Syntax
```
len INPUT
```
##### Example
```
{{ if gt (len .Notification) 0 }}<div id="notification">
    <div class="success">
        <p>{{ .Notification | safeHTML }}</p>
    </div>
</div>{{ end }}
```

#### ● map
##### Syntax
```
map KEY VALUE [KEY VALUE]...
```
##### Example
```
{{ $m := map "key1" 100 "key2" 200 "key3" 300 }}
{{ printf "%#v" $m }}
-> map[string]interface {}{"key1":100, "key2":200, "key3":300}

{{ $m := map "key1" 100 "key2" 200 "key3" 300 "key4"}}
{{ printf "%#v" $m }}
-> 
```

#### ● month
##### Syntax
```
month
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ month }}
-> 2
```

#### ● replace
##### Syntax
```
replace INPUT OLD NEW
```
##### Example
```
<span>{{ replace "Is this an apple?" "an apple" "a pen" }}</span>
-> <span>Is this a pen?</span>
```

#### ● replaceRE
##### Syntax
```
replaceRE PATTERN REPLACEMENT INPUT
```
##### Example
```
{{ replaceRE "^https?://([^/]+).*" "$1" "https://github.com/WhiteRaven777/simple-render" }}
-> github.com

{{ "https://github.com/WhiteRaven777/simple-render" | replaceRE "^https?://([^/]+).*" "$1" }}
-> github.com
```

#### ● safeCSS
##### Syntax
```
safeCSS INPUT
```
##### Example
```
<p style="{{ .Style }}">...</p>
-> <p style="ZgotmplZ">...</p>

<p style="{{ .Style | safeCSS }}">...</p>
-> <p style="color: red;">...</p>
```

#### ● safeHTML
##### Syntax
```
safeHTML INPUT
```
##### Example
```
Link = `<a href="https://example.com">sample</a>`

{{ .Link }}
-> &lt;a href=&#34;https://example.com&#34;&gt;sample&#34;/a&gt;

{{ .Link | safeHTML }}
-> <a href="https://example.com">sample</a>
```

#### ● safeHTMLAttr
##### Syntax
```
safeHTMLAttr INPUT
```
##### Example
```
Url = "https://example.com"

<a href="{{ .Url }}">
-> <a href="#ZgotmplZ">

<a {{ printf "href=%q" .Url | safeHTMLAttr }}>
-> <a href="https://example.com">
```

#### ● safeJS
##### Syntax
```
safeJS INPUT
```
##### Example
```
Hash = "abc123"

<script>var form_{{ .Params.hash }}</script>
-> <script>var form_"abc123"</script>

<script>var form_{{ .Params.hash | safeJS }}</script>
-> <script>var form_abc123</script>
```

#### ● safeURL
##### Syntax
```
safeURL INPUT
```
##### Example
```
Url = "https://example.com"

<a href="{{ .Url }}">
-> <a href="#ZgotmplZ">

<a href="{{ .Url | safeURL }}">
-> <a href="https://example.com">
```

#### ● slice
##### Syntax
```
slice ITEM...
```
##### Example
```
{{ print (slice 0 1 2)}}
-> [0 1 2]
```

#### ● split
##### Syntax
```
split STRING DELIMITER
```
##### Example
```
{{ slice "SAMPLE-TEXT" "-" }}
-> [SAMPLE TEXT]

{{ slice "AAA+BBB-CCC+DDD" "+" }}
-> [AAA BBB-CCC DDD]
```

#### ● time
##### Syntax
```
time
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ time }}
-> 17:18:19
```

#### ● trim, trimLeft, trimRight
##### Syntax
```
trim
trimLeft
trimRight
```
##### Example
```
{{ trim "!!?? abcdef ??!!" "!?" }}
-> " abcdef "

{{ trim " abcdef " }}
-> "abcdef"

{{ trimLeft "!!?? abcdef ??!!" "!?" }}
-> " abcdef ??!!"

{{ trimLeft " abcdef " }}
-> "abcdef "

{{ trimRight "!!?? abcdef ??!!" "!?" }}
-> "!!?? abcdef "

{{ trimRight " abcdef " }}
-> " abcdef"

```

#### ● year
##### Syntax
```
year
```
##### Example
```
# today: 2019-02-14 17:18:19 (UTC)

{{ year }}
-> 2019
```
