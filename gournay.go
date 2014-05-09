package main

import (
  "fmt"
  "net/http"
  "os"
  "io"
)

func main() {
  http.HandleFunc("/new/", newHandler)
  http.HandleFunc("/create/", createHandler)
  http.HandleFunc("/show/", showHandler)
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
      <form action="/create/" method="POST">
        <label>New title</label>
        <input type="text" name="title"/>
        <input type="submit" value="Save" />
      </form>
    </body>
    </html>`,
  )
}

func createHandler(w http.ResponseWriter, r *http.Request) {
  title := r.FormValue("title")
  fmt.Println("title", title)
  http.Redirect(w, r, "/show/", http.StatusFound)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
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
    <p>Some data goes here</p>
    </body>
    </html>`,
  )
}
