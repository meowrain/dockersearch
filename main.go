package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/meowrain/dockersearch/internal/display"
	"github.com/meowrain/dockersearch/internal/display/handlers"
	"github.com/meowrain/dockersearch/internal/query"
	"github.com/sirupsen/logrus"

	"github.com/meowrain/dockersearch/utils"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

var tmpl *template.Template

func initLogger() {
	logrus.SetLevel(logrus.DebugLevel)
}
func initTemplates() {
	// 获取模板函数映射
	funcMap := utils.GetFuncMap()

	// 使用嵌入的模板文件
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html"))
	log.Println("Templates loaded from embedded filesystem")
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			fmt.Println("Starting Docker Search HTTP server on :8096")
			initLogger()
			initTemplates()
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				handlers.IndexHandler(tmpl, staticFS, w, r)
			})
			log.Fatal(http.ListenAndServe(":8096", nil))
		} else {
			imageName := os.Args[1]

			page := 1
			limit := 25

			for i := 2; i < len(os.Args); i++ {
				switch os.Args[i] {
				case "--page", "-p":
					if i+1 < len(os.Args) {
						fmt.Sscanf(os.Args[i+1], "%d", &page)
						i++
					}
				case "--limit", "-l":
					if i+1 < len(os.Args) {
						fmt.Sscanf(os.Args[i+1], "%d", &limit)
						i++
					}
				}
			}

			searchResponse := query.QueryImageInfo(imageName, page, limit)
			display.DisplayPrettyTable(searchResponse)

			if searchResponse.NumPages > 1 {
				fmt.Printf("\n页码: %d/%d  ", searchResponse.Page, searchResponse.NumPages)
				fmt.Printf("(使用 --page <页码> 参数查看其他页)\n")
			}
		}
	} else {
		fmt.Println("Usage: dockersearch <ImageName> [--page <number>] [--limit <number>] or dockersearch server")
	}
}
