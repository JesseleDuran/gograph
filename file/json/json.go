package json

import (
  "encoding/json"
  "fmt"
  "log"
  "os"
)

// Write a .json file on the given path from the given value.
// The function returns the created file name.
func Write(name string, content interface{}) string {
  f, err := os.Create(name)
  if err != nil {
    fmt.Println(err)
    return ""
  }
  d2, _ := json.Marshal(content)
  n2, err := f.Write(d2)
  if err != nil {
    log.Println(err)
    f.Close()
    return ""
  }
  log.Println(n2, "bytes written successfully")
  err = f.Close()
  if err != nil {
    log.Println(err)
    return ""
  }
  return f.Name()
}
