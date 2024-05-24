package main

import (
  "fmt"
  "strconv"
  "os"
  "io"
  "bytes"
  "net/http"
  "net/url"
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

func getPrettyFormData(v url.Values) (bytes.Buffer) {
  var buffer bytes.Buffer 
  for key, value := range(v) {
    if debug {
      // move this to debug funtion
      fmt.Print("Found form data: ")
      fmt.Printf("%v: %v\n", key, value)
    }
    buffer.WriteString(fmt.Sprintf("%v: %v\n", key, value))
  }
  return buffer
}

func sendmail(mb bytes.Buffer) {
  fmt.Println("Mail sent:")
  fmt.Printf("%s", mb)
}

// route handler
func getSubmission(w http.ResponseWriter, r *http.Request) {
  // move this out to debug function
  if debug {
    debugFormData(r)
  }

  fmt.Println("Received form submission for default form\n")
  r.ParseForm()
  formData := r.PostForm
  // if validateInput(formData) {
  // text := templateText(formData)
  pretty := getPrettyFormData(formData)
  mailBody := pretty
  sendmail(mailBody)
  io.WriteString(w, "\n")
}

func main() {
  getEnvConfig()

  http.HandleFunc("/submit", getSubmission)
  http.ListenAndServe(":3333", nil)
}
