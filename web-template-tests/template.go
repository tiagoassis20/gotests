package main

import (
	"os"
	"text/template"
)

type sages struct {
	Title string
	Sages []string
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("list.gohtml"))

}

func main() {
	sage := sages{
		Title: "Great Sages of the Would",
		Sages: []string{"Abraão", "dalai hilama", "buda", "Martin luter king", "jun tsu"},
	}

	tpl.ExecuteTemplate(os.Stdout, "list.gohtml", sage)
}
