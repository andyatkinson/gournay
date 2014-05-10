package main

import (
  "fmt"
  "net/http"
  "os"
  "html/template"
  "crypto/md5"
  "encoding/hex"
)

var storage map[string]string

func GetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
  storage = make(map[string]string)

  http.HandleFunc("/", newHandler)
  http.HandleFunc("/create", createHandler)
  http.HandleFunc("/find", findHandler)
  err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
  if err != nil {
    panic(err)
  }
}

func newHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set( // set a header
    "Content-Type",
    "text/html",
  )
  data := storage
  tmpl, err := template.New("form page").Parse(
    `<!DOCTYPE html>
    <html>
    <head>
      <title>Gournay URL shrinker</title>
      <!-- Latest compiled and minified CSS -->
      <link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
    </head>
    <body>
    <div class="container" style="padding:50px 0;">
      <h2>Gournay URL shrinker</h2>
      <a href="https://github.com/andyatkinson/gournay">Source on github</a>
      <form role="form" action="/create" method="POST">
        <div class="form-group">
          <label for="url-input">URL</label>
          <input type="text" class="form-control" id="url-input" name="url" placeholder="http://cnn.com" />
        </div>
        <button type="submit" class="btn btn-primary">Shrink URL</button>
      </form>

      <div class="links" style="padding: 20px 0;">
        <ul>
        {{range $index, $element := .}}
          <li>
            <form role="form" action="/find" method="POST">
              <div class="form-group">
                <input type="hidden" name="hash" value="{{$index}}" />
                <button type="submit" class="btn btn-link">{{$index}}</button>
              </div>
            </form>
          </li>
        {{end}}
        </ul>
      </div>
    </div>
    </body>
    </html>`)
    err = tmpl.Execute(w, data)
    if err != nil { panic(err) }
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  url := r.FormValue("url")
  hash := GetMD5Hash(url)[0:5]
  storage[hash] = url
  http.Redirect(w, r, "/new", 301)
}

func findHandler(w http.ResponseWriter, r *http.Request) {
  hashQuery := r.FormValue("hash")
  result := storage[hashQuery]
  http.Redirect(w, r, result, 301)
}
