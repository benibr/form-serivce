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
  "log/slog"
)

// globals
var debug bool = false
var email []string

func getEnvConfig() {
  debug, _ = strconv.ParseBool(os.Getenv("FORM_SERVICE_DEBUG"))
  email    = append(email, os.Getenv("FORM_SERVICE_EMAIL"))
  //slog.Debug("Option set", "FORM_SERVICE_EMAIL", email)
}

func getPrettyFormData(v url.Values) (bytes.Buffer) {
  var buffer bytes.Buffer 
  for key, value := range(v) {
    slog.Debug("Found form data: ", fmt.Sprintf("%v", key), fmt.Sprintf("%v", value))
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
  slog.Info("Mail sent")
}

func getSubmission(w http.ResponseWriter, r *http.Request) {
  slog.Info("Received form submission for default form")
  formData := r.PostForm
  slog.Debug(fmt.Sprintf("%v",formData))
  r.ParseForm()
  // TODO: if validateInput(formData) {
  // TODO: text := templateText(formData)
  pretty := getPrettyFormData(formData)
  mailBody := pretty
  sendmail(mailBody)
  io.WriteString(w, "You message has been sent! Thank you for writing to us.")
}

func main() {
  logLevel := new(slog.LevelVar)
  logLevel.Set(slog.LevelInfo)
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
  slog.SetDefault(logger)
  slog.Info("form-service starting")

  getEnvConfig()
  if debug {
    logLevel.Set(slog.LevelDebug)
  }

  http.HandleFunc("/submit", getSubmission)
  http.ListenAndServe(":3333", nil)
}
