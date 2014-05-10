package main

import (
  "fmt"
  "net/http"
  "os"
  "io"
  "html/template"
  "crypto/md5"
  "encoding/hex"
)

type Data struct {
  Url string
  Hash string
}

var storage map[string]string

func main() {
  storage = make(map[string]string)

  http.HandleFunc("/", newHandler)
  http.HandleFunc("/create", createHandler)
  http.HandleFunc("/find", findHandler)
  fmt.Println("listening...")
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
  io.WriteString(
    w,
    `<doctype html>
    <html>
    <head><title>Hello World</title></head>
    <body>
      <form action="/create" method="POST">
        <label>Create hash for URL</label>
        <input type="text" name="url"/>
        <input type="submit" value="Save" />
      </form>

      <form action="/find" method="POST">
        <label>Look up hash</label>
        <input type="text" name="hash"/>
        <input type="submit" value="Find" />
      </form>
    </body>
    </html>`,
  )
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  url := r.FormValue("url")
  hash := GetMD5Hash(url)[0:5]
  storage[hash] = url
  data := Data{url, hash}
  fmt.Println("storage", storage)
  fmt.Println("data", data)
  tmpl, err := template.New("urlResp").Parse(`<!doctype html>
    <html><body>
    <p>{{.Hash}}</p>
    <a href="/">Back</a>
    </body></html>
    </html>`)
  if err != nil { panic(err) }
  err = tmpl.Execute(w, data)
  if err != nil { panic(err) }
}

func GetMD5Hash(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}

func findHandler(w http.ResponseWriter, r *http.Request) {
  hashQuery := r.FormValue("hash")
  result := storage[hashQuery]
  fmt.Println("result", result)
  data := Data{result, hashQuery}
  fmt.Println("data", data)
  tmpl, err := template.New("result").Parse(`<!doctype html>
    <html><body>
    <p>{{.Url}}</p>
    <a href="/">Back</a>
    </body></html>
    </html>`)
  if err != nil { panic(err) }
  err = tmpl.Execute(w, data)
  if err != nil { panic(err) }
}
