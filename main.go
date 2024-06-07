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

// exit codes
const ecMissingParam = 3

// globals
var debug bool = false
var email []string
var logLevel = new(slog.LevelVar)

func getEnvConfig() {
  //debug
  debug, _ = strconv.ParseBool(os.Getenv("FORM_SERVICE_DEBUG"))
  if debug {
    logLevel.Set(slog.LevelDebug)
    slog.Debug("Option set", "FORM_SERVICE_DEBUG", debug)
  }

  //email
  var env_email string = os.Getenv("FORM_SERVICE_EMAIL")
  if len(env_email) == 0 {
    slog.Error("Env var FORM_SERVICE_EMAIL not set, cannot continue!")
    os.Exit(ecMissingParam)
  } else {
    email = append(email, env_email)
  }
  slog.Debug("Option set", "FORM_SERVICE_EMAIL", email)
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

  err := smtp.SendMail("mail.madways.de:25", nil, "forms@madways.de", email, msg)
  if err != nil {
    slog.Error(fmt.Sprintf("Cannot send mail to %v", email))
    slog.Error(fmt.Sprintf("%v", err))
  } else {
    slog.Info(fmt.Sprintf("Mail sent to %v", email))
  }
}

func getSubmission(response http.ResponseWriter, request *http.Request) {
  slog.Info("Received form submission for default form")
  formData := request.PostForm
  slog.Debug(fmt.Sprintf("%v",formData))
  request.ParseForm()
  // TODO: if validateInput(formData) {
  // TODO: text := templateText(formData)
  pretty := getPrettyFormData(formData)
  mailBody := pretty
  sendmail(mailBody)
  io.WriteString(response, "You message has been sent! Thank you for writing to us.")
}

func main() {
  logLevel.Set(slog.LevelInfo)
  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
  slog.SetDefault(logger)
  slog.Info("form-service starting")

  getEnvConfig()

  http.HandleFunc("/submit", getSubmission)
  err := http.ListenAndServe(":3333", nil)
  if err != nil {
    slog.Error("Cannot open TCP listener on port 3333")
    slog.Error(fmt.Sprintf("%v", err))
  }
}
