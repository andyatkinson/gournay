package main

import (
  "fmt"
  "os"
  "net/http"
  "html/template"
  "crypto/md5"
  "encoding/hex"
  _ "github.com/lib/pq"
  "database/sql"
)

var db *sql.DB

func GetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
  var err error
  // need to set up a connection string ENV var, and a variable to control sslmode in heroku
  connectionString := os.Getenv("DATABASE_URL")
  db, err = sql.Open("postgres", connectionString)
  if err != nil {
    panic(err)
  }

  http.HandleFunc("/", newHandler)
  http.HandleFunc("/create", createHandler)
  http.HandleFunc("/find", findHandler)

  port := os.Getenv("PORT")
  if port == "" {
    port = "5000"
  }
  http.ListenAndServe(":" + port, nil)
  defer db.Close()
}

func newHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set( // set a header
    "Content-Type",
    "text/html",
  )

  var (
    url string
    hash string
  )
  data := make(map[string]string)
  rows, err := db.Query("SELECT url, hash FROM entries")
  if err != nil {
    panic(err)
  }
  for rows.Next() {
    err := rows.Scan(&url, &hash)
    if err != nil {
      panic(err)
    }
    data[hash] = url
  }
  defer rows.Close()

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
  fmt.Println("got url:", url)
  hash := GetMD5Hash(url)[0:5]
  stmt, err := db.Prepare("INSERT INTO entries (url, hash) VALUES ($1, $2)")
  if err != nil {
    panic(err)
  }
  res, err := stmt.Exec(url, hash)
  if err != nil {
    panic(err)
  }
  lastId, err := res.LastInsertId()
  fmt.Println("inserted records, last ID",lastId)
  http.Redirect(w, r, "/new", 301)
}

func findHandler(w http.ResponseWriter, r *http.Request) {
  hashQuery := r.FormValue("hash")
  fmt.Println("got hash query:", hashQuery)
  var url string
  err := db.QueryRow("SELECT url from entries WHERE hash = $1", hashQuery).Scan(&url)
  fmt.Println("Redirecting to....", url)
  if err != nil {
    panic(err)
  }
  http.Redirect(w, r, url, 301)
}
