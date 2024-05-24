package main

import (
  "fmt"
  "strconv"
  "os"
  "io"
  "bytes"
  "net/http"
)

// globals
var debug bool = false

func debugFormData(r *http.Request) {
  fmt.Printf("DEBUG: Method: %v\n", r.Method)
  fmt.Print("DEBUG: Full form object:")
  fmt.Printf(" %v\n\n", r)
}

func getEnvConfig() {
  debug, _ = strconv.ParseBool(os.Getenv("FORM_SERVICE_DEBUG"))
}

func getPrettyFormData(r *http.Request) (bytes.Buffer) {
  var buffer bytes.Buffer 
  for key, value := range(r.PostForm) {
    if debug {
      fmt.Print("Found form data: ")
      fmt.Printf("%v: %v\n", key, value)
    }
    buffer.WriteString(fmt.Sprintf("%v: %v\n", key, value))
  }
  return buffer
}

// route handler
func getSubmission(w http.ResponseWriter, r *http.Request) {

  fmt.Println("Received form submission\n")

  if debug {
    debugFormData(r)
  }
  r.ParseForm()
  io.WriteString(w, "Form submitted, I promise!!\n")
}

func main() {
  getEnvConfig()

  http.HandleFunc("/submit", getSubmission)
  http.ListenAndServe(":3333", nil)
}
