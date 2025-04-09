package handlers

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/meowrain/dockersearch/internal/models"
	"github.com/meowrain/dockersearch/internal/query"
)

type IndexPageData struct {
	Title string
}

type SearchPageData struct {
	Title      string
	Query      string
	NumResults int
	Results    []models.Result
	Page       int
	NumPages   int
	Limit      int
}

func IndexHandler(tmpl *template.Template, staticFS embed.FS, w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
	switch r.URL.Path {
	case "/", "/index.html":
		data := IndexPageData{
			Title: "Docker 镜像搜索",
		}
		err := tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case "/robots.txt":
		writeResponse(w, http.StatusOK, "text/plain", "User-agent: *\nDisallow: /")
	case "/search":
		q := r.URL.Query().Get("q")
		if q == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		page := 1
		limit := 25

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			fmt.Sscanf(pageStr, "%d", &page)
		}
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			fmt.Sscanf(limitStr, "%d", &limit)
		}

		result := query.QueryImageInfo(q, page, limit)

		data := SearchPageData{
			Title:      "Docker 搜索结果",
			Query:      q,
			NumResults: result.NumResults,
			Results:    result.Results,
			Page:       result.Page,
			NumPages:   result.NumPages,
			Limit:      limit,
		}

		err := tmpl.ExecuteTemplate(w, "search.html", data)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		if strings.HasPrefix(r.URL.Path, "/static/") {
			serveStaticFile(staticFS, w, r)
			return
		}
		writeResponse(w, http.StatusNotFound, "text/plain", "404 Not Found")
	}
}

func serveStaticFile(staticFS embed.FS, w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	content, err := staticFS.ReadFile("static/" + filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	contentType := "text/plain"
	switch {
	case strings.HasSuffix(filePath, ".css"):
		contentType = "text/css"
	case strings.HasSuffix(filePath, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(filePath, ".png"):
		contentType = "image/png"
	case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
		contentType = "image/jpeg"
	case strings.HasSuffix(filePath, ".svg"):
		contentType = "image/svg+xml"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(content)
}
