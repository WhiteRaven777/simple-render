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
layout.html
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

view.html
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

#### safeCSS
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

#### safeHTML
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

#### safeHTMLAttr
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

#### safeJS
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

#### safeURL
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

#### slice
##### Syntax
```
slice ITEM...
```
##### Example
```
{{ print (slice 0 1 2)}}
-> [0 1 2]
```
