package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logging function to log string by prefix with 2 params (prefix string, logcontent string)
func Logging(_prefix string, _s string) {
	_log := fmt.Sprintf("%s: %s\n", time.Unix(time.Now().Unix(), 0), _s)
	_dirname := "logs"
	if _, err := os.Stat(_dirname); os.IsNotExist(err) {
		os.Mkdir(_dirname, 0755)
	}
	_fname := fmt.Sprintf("%s/%s-%s.log", _dirname, _prefix, time.Now().Format("200601"))
	f, err := os.OpenFile(_fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(_log); err != nil {
		log.Println(err)
	}
}

func ErrorLogging(_prefix string, _s string) {
	Logging("error", fmt.Sprintf("[%s] %s", _prefix, _s))
}
