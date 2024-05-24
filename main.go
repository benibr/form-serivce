package main

import (
  "fmt"
  "strconv"
  "os"
  "io"
  "bytes"
  "net/http"
  "net/url"
  "net/smtp"
)

// globals
var debug bool = false
var email []string

func debugFormData(r *http.Request) {
  fmt.Printf("DEBUG: Method: %v\n", r.Method)
  fmt.Print("DEBUG: Full form object:")
  fmt.Printf(" %v\n\n", r)
}

func getEnvConfig() {
  debug, _ = strconv.ParseBool(os.Getenv("FORM_SERVICE_DEBUG"))
  email    = append(email, os.Getenv("FORM_SERVICE_EMAIL"))
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

func sendmail(body bytes.Buffer) {
  msg := []byte("To: " + email[0] + "\n" +
                "From: forms@madways.de\n" +
                "Subject: Form ausjefüddelt!" +
                "\n\n" +
                body.String())

  smtp.SendMail("mail.madways.de:25", nil, "forms@madways.de", email, msg)
  fmt.Println("Mail sent")
}

// route handler
func getSubmission(w http.ResponseWriter, r *http.Request) {
  // move this out to debug function
  if debug {
    debugFormData(r)
  }

  fmt.Println("Received form submission for default form")
  r.ParseForm()
  formData := r.PostForm
  // TODO: if validateInput(formData) {
  // TODO: text := templateText(formData)
  pretty := getPrettyFormData(formData)
  mailBody := pretty
  sendmail(mailBody)
  io.WriteString(w, "You message has been sent! Thank you for writing to us.")
}

func main() {
  getEnvConfig()

  http.HandleFunc("/submit", getSubmission)
  http.ListenAndServe(":3333", nil)
}
