package main

import (
	"keypass/lib/controllers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	controllers.RunApp()
}
