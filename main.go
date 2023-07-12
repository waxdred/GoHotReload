package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/waxdred/GoHotReload/models"
)

func main() {
	app := models.New().CheckingParse()
	if app.Errror() != nil {
		fmt.Println(app.Errror())
		os.Exit(-1)
	}
	err := app.Listen().Start().Errror()
	if err != nil {
		fmt.Println(err)
	}
}
