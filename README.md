# simple-render
This is a simple render library designed to simplify the rendering of JSON and HTML responses in web applications.
Additionally, it provides a set of custom template functions that enhance the flexibility and power of your HTML templates.

Currently, the following formats are supported:

* HTML
* JSON

By leveraging the `"html/template"` package's `Funcs()` method, this library extends the default template functionality with a variety of useful custom functions.
These functions are detailed in the **Functions** section below.
Below are examples and explanations on how to use the library effectively, including the available custom template functions.

# Installation
To install the package, use the following command:

```sh
go get github.com/WhiteRaven777/simple-render
```

# How To Use

## JSON
The following example demonstrates how to use the render.JSON function to send JSON responses.

```go
type Data struct {
Date string
Msg  string
}

data := Data {
Date: time.Now().Format("2006-01-02"),
Msg:  "message",
}

--- or ---

data := map[string]string{
"Date": time.Now().Format("2006-01-02"),
"Msg":  "message",
}
```

To send this data as a JSON response:

```go
func Sample(w http.ResponseWriter, r *http.Request) {
render.JSON(w, http.StatusOK, data)
}
```

This will produce the following JSON output:

```json
{
  "Date": "2019-01-17",
  "Msg": "message"
}
```

## HTML
The following example demonstrates how to use the Template struct and its HTML method to render HTML responses using templates.

### layout.html
The layout template defines the overall structure of the HTML document.

```html
{{- define "layout" }}
<!doctype html>
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
</html>
{{- end }}
```

### view.html
The view template defines the specific content for a particular page.

```html
{{ define "title" }}{{ .Title }}{{ end }}

{{ define "body" }}
<h1>{{ .Title }}</h1>
<article>
{{ .Body | safeHTML }}
</article>
<li><a href="{{ .Url | safeURL }}">{{ .Title }}</a></li>
{{ end }}
```

To render this template with data using the Template struct:

```go
func Sample(w http.ResponseWriter, r *http.Request) {
    tmpl := Template{
        Layout: "layout",
        Data: map[string]string{
            "Title": "title text",
            "Body":  "body text",
            "Url":   "https://github.com/WhiteRaven777/simple-render",
        },
        ExtraFuncMap: template.FuncMap{
            "customFunc": func() string { return "Custom Function" },
        },
    }
    tmpl.HTML(w, http.StatusOK)
}
```

In this example, the Template struct's HTML method renders the view.html template within the layout.html template, passing the Data map to populate the template variables. The ExtraFuncMap allows for the inclusion of custom functions within the template.

### Functions
This section provides a list of custom functions available for use within templates. Each function is designed to simplify template logic and enhance the flexibility of your templates.

#### ● day
Returns the current day of the month (UTC).

##### Syntax
```markdown
day
```

##### Example
```markdown
# Assuming today is 2019-02-14 17:18:19 (UTC)

{{ day }}
-> 14
```

#### ● date
Returns the current date in YYYY-MM-DD format (UTC).

##### Syntax
```markdown
date
```

##### Example
```markdown
# Assuming today is 2019-02-14 17:18:19 (UTC)

{{ date }}
-> 2019-02-14
```

#### ● datetime
Returns the current date and time in RFC3339 format (UTC).

##### Syntax
```markdown
datetime
```

##### Example
```markdown
# Assuming today is 2019-02-14 17:18:19 (UTC)

{{ datetime }}
-> 2019-02-14T17:18:19Z
```

#### ● default
Returns the input value if it is non-empty; otherwise, returns the default value.

##### Syntax
```markdown
default DEFAULT_VALUE INPUT
```

##### Example
```markdown
{{ default 2 1 }}
-> 1

{{ default "default" "" }}
-> default

{{ default 42 0 }}
-> 42
```

#### ● dict
Creates a dictionary (map) from a list of key-value pairs.

##### Syntax
```markdown
dict [KEY VALUE]...
```

##### Example
```markdown
{{ dict 0 "aaa" 1 "bbb" 2 "ccc" }}
-> map[0:"aaa", 1:"bbb", 2:"ccc"]

{{ dict "name" "John" "age" 30 }}
-> map["name":"John", "age":30]
```
Note: The dict function requires an even number of arguments.

#### ● eval
Evaluates a mathematical or logical expression and returns the result.

##### Syntax
```markdown
eval FORMULAS
```

##### Example
```markdown
{{ eval "1 + 1" }}
-> 2

{{ eval "1 - 1 + 1" }}
-> 1

# Errors in the evaluation will return nil.

{{ eval "invalid expression" }}
-> nil
```

#### ● findRE
Finds all matches of a regular expression in the input string.

##### Syntax
```markdown
findRE PATTERN INPUT [LIMIT]
```

##### Example
```markdown
{{ findRE `/hoge/(\d+)` "/hoge/1234567890/hoge/987654321/abcdefghijk"   }}
-> [/hoge/1234567890 /hoge/987654321]

{{ findRE `/hoge/(\d+)` "/hoge/1234567890/hoge/987654321/abcdefghijk" 1 }}
-> [/hoge/1234567890]
```

#### ● in
Checks if the item is present within the set.

##### Syntax
```markdown
in SET ITEM
```

##### Example
```markdown
{{ if in "/example/aaa" "/example" }}True{{ else }}False{{ end }}
-> True

{{ if in "/sample/aaa" "/example" }}True{{ else }}False{{ end }}
-> False
```

#### ● index
Returns the element at the specified index or key from a collection (array, slice, or map).

##### Syntax
```markdown
index COLLECTION (INDEX|KEY)
```

##### Example
```markdown
{{ index (slice 3 4 5) 0 }}
-> 3

{{ index (dict 0 "aaa" 1 "bbb" 2 "ccc") 0 }}
-> aaa
```

#### ● len
Returns the length of the input. Supports arrays, maps, slices, and strings.

##### Syntax
```markdown
len INPUT
```

##### Example
```markdown
{{ if gt (len .Notification) 0 }}<div id="notification">
<div class="success">
<p>{{ .Notification | safeHTML }}</p>
</div>
</div>{{ end }}
```

#### ● lower
Converts the input string(s) to lowercase.

##### Syntax
```markdown
lower INPUT
```

##### Example
```markdown
{{ $arr := slice "A" "B" "C" }}
{{ range $arr }}{{ lower . }}{{ end }}
-> a b c
```

#### ● map
Creates a map from a list of key-value pairs.

##### Syntax
```markdown
map KEY VALUE [KEY VALUE]...
```

##### Example
```markdown
{{ $m := map "key1" 100 "key2" 200 "key3" 300 }}
{{ printf "%#v" $m }}
-> map[string]interface {}{"key1":100, "key2":200, "key3":300}

{{ $m := map "key1" 100 "key2" 200 "key3" 300 "key4"}}
{{ printf "%#v" $m }}
-> 
```

#### ● month
Returns the current month as an integer (UTC).

##### Syntax
```markdown
month
```

##### Example
```markdown
# Assuming today is 2019-02-14 17:18:19 (UTC)

{{ month }}
-> 2
```

#### ● replace
Replaces all occurrences of the old substring with the new substring in the input string.

##### Syntax
```markdown
replace INPUT OLD NEW
```

##### Example
```markdown
<span>{{ replace "Is this an apple?" "an apple" "a pen" }}</span>
-> <span>Is this a pen?</span>
```

#### ● replaceRE
Replaces all matches of the regular expression pattern with the replacement string in the input string.

##### Syntax
```markdown
replaceRE PATTERN REPLACEMENT INPUT
```

##### Example
```markdown
{{ replaceRE "^https?://([^/]+).*" "$1" "https://github.com/WhiteRaven777/simple-render" }}
-> github.com

{{ "https://github.com/WhiteRaven777/simple-render" | replaceRE "^https?://([^/]+).*" "$1" }}
-> github.com
```

#### ● safeCSS
Marks the input as safe CSS content to prevent escaping.

##### Syntax
```markdown
safeCSS INPUT
```

##### Example
```markdown

<p style="{{ .Style }}">...</p>
-> <p style="ZgotmplZ">...</p>

<p style="{{ .Style | safeCSS }}">...</p>
-> <p style="color: red;">...</p>
```

#### ● safeHTML
Marks the input as safe HTML content to prevent escaping.

##### Syntax
```markdown
safeHTML INPUT
```

##### Example
```markdown
# Link = `<a href="https://example.com">sample</a>`

{{ .Link }}
-> &lt;a href=&#34;https://example.com&#34;&gt;sample&#34;/a&gt;

{{ .Link | safeHTML }}
-> <a href="https://example.com">sample</a>
```

#### ● safeHTMLAttr
Marks the input as a safe HTML attribute to prevent escaping.

##### Syntax
```markdown
safeHTMLAttr INPUT
```

##### Example
```markdown
Url = "https://example.com"

<a href="{{ .Url }}">
-> <a href="#ZgotmplZ">

<a {{ printf "href=%q" .Url | safeHTMLAttr }}>
-> <a href="https://example.com">
```

#### ● safeJS
Marks the input as safe JavaScript content to prevent escaping.

##### Syntax
```markdown
safeJS INPUT
```

##### Example
```markdown
Hash = "abc123"

<script>var form_{{ .Params.hash }}</script>
-> <script>var form_"abc123"</script>

<script>var form_{{ .Params.hash | safeJS }}</script>
-> <script>var form_abc123</script>
```

#### ● safeURL
Marks the input as a safe URL to prevent escaping.

##### Syntax
```markdown
safeURL INPUT
```

##### Example
```markdown
Url = "https://example.com"

<a href="{{ .Url }}">
-> <a href="#ZgotmplZ">

<a href="{{ .Url | safeURL }}">
-> <a href="https://example.com">
```

#### ● slice
Creates a slice from the input elements.

##### Syntax
```markdown
slice ITEM...
```

##### Example
```markdown
{{ print (slice 0 1 2) }}
-> [0 1 2]
```

#### ● split
Splits the input string by the specified delimiter and returns a slice of substrings.

##### Syntax
```markdown
split STRING DELIMITER
```

##### Example
```markdown
{{ split "SAMPLE-TEXT" "-" }}
-> [SAMPLE TEXT]

{{ split "AAA+BBB-CCC+DDD" "+" }}
-> [AAA BBB-CCC DDD]
```

#### ● time
Returns the current time in HH:MM:SS format (UTC).

##### Syntax
```markdown
time
```

##### Example
```markdown

Assuming the current time is 17:18:19 (UTC)
{{ time }}
-> 17:18:19
```

#### ● trim
Trims the specified characters from both ends of the input string. If no characters are specified, it trims whitespace by default.

##### Syntax
```markdown
trim [CHARACTERS]
```

##### Example
```markdown
{{ trim "!!?? abcdef ??!!" "!?" }}
-> " abcdef "

{{ trim " abcdef " }}
-> "abcdef"
```

#### ● trimLeft
Trims the specified characters from the left side of the input string. If no characters are specified, it trims whitespace by default.

##### Syntax
```markdown
trimLeft [CHARACTERS]
```

##### Example
```markdown
{{ trimLeft "!!?? abcdef ??!!" "!?" }}
-> " abcdef ??!!"

{{ trimLeft " abcdef " }}
-> "abcdef "
```

#### ● trimRight
Trims the specified characters from the right side of the input string. If no characters are specified, it trims whitespace by default.

##### Syntax
```markdown
trimRight [CHARACTERS]
```

##### Example
```markdown
{{ trimRight "!!?? abcdef ??!!" "!?" }}
-> "!!?? abcdef "

{{ trimRight " abcdef " }}
-> " abcdef"
```

#### ● upper
Converts the input string(s) to uppercase.

##### Syntax
```markdown
upper INPUT
```

##### Example
```markdown
{{ $arr := slice "a" "b" "c" }}
{{ range $arr }}{{ upper . }}{{ end }}
-> A B C
```

#### ● year
Returns the current year as an integer (UTC).

##### Syntax
```markdown
year
```

##### Example
```markdown
# Assuming today is 2019-02-14 17:18:19 (UTC)

{{ year }}
-> 2019
```

# License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/WhiteRaven777/simple-render/blob/master/LICENSE) file for details.
