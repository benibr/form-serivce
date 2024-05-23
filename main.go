package main

import (
  "fmt"
  "io"
  "net/http"
)

func getSubmission(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Received form submission\n")
  io.WriteString(w, "Form submitted, I promise!!\n")
}

func main() {
  http.HandleFunc("/submit", getSubmission)

  http.ListenAndServe(":3333", nil)
}
