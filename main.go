package main

import (
	"dockersearch/templates"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

const (
	baseURL            = "https://ai-proxy.hidewiki.top/docker/v1/search?"
	imageDetailBaseURL = "https://ai-proxy.hidewiki.top/hubdocker/v2/repositories/"
)

type SearchResponse struct {
	NumPages   int      `json:"num_pages"`
	NumResults int      `json:"num_results"`
	PageSize   int      `json:"page_size"`
	Page       int      `json:"page"`
	Query      string   `json:"query"`
	Results    []Result `json:"results"`
}

type Result struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PullCount   int    `json:"pull_count"`
	StarCount   int    `json:"star_count"`
	IsTrusted   bool   `json:"is_trusted"`
	IsOfficial  bool   `json:"is_official"`
	IsAutomated bool   `json:"is_automated"`
}

var tmpl *template.Template

func initTemplates() {
	// 获取模板函数映射
	funcMap := templates.GetFuncMap()

	// 使用嵌入的模板文件
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html"))
	log.Println("Templates loaded from embedded filesystem")
}

type IndexPageData struct {
	Title string
}

type SearchPageData struct {
	Title      string
	Query      string
	NumResults int
	Results    []Result
	Page       int
	NumPages   int
	Limit      int
}

func handler(w http.ResponseWriter, r *http.Request) {
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

		result := queryImageInfo(q, page, limit)

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
			serveStaticFile(w, r)
			return
		}
		writeResponse(w, http.StatusNotFound, "text/plain", "404 Not Found")
	}
}

func serveStaticFile(w http.ResponseWriter, r *http.Request) {
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

func limitString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func writeResponse(w http.ResponseWriter, statusCode int, contentType, body string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func formatPullCount(count int) string {
	if count >= 1000000000 {
		return fmt.Sprintf("%.1fB", float64(count)/1000000000)
	} else if count >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	} else if count >= 1000 {
		return fmt.Sprintf("%.1fK", float64(count)/1000)
	}
	return fmt.Sprintf("%d", count)
}

func queryImageInfo(imageName string, page int, limit int) *SearchResponse {
	log.Println("search Image:", imageName, "page:", page, "limit:", limit)
	queryParams := url.Values{}
	queryParams.Add("n", fmt.Sprintf("%d", limit))
	queryParams.Add("page", fmt.Sprintf("%d", page))
	queryParams.Add("type", "image")
	queryParams.Add("q", imageName)

	targetURL := baseURL + queryParams.Encode()
	log.Println("query URL: ", targetURL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:137.0) Gecko/20100101 Firefox/137.0")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: Request failed with status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return &searchResponse
}

func displayPrettyTable(results *SearchResponse) {
	bold := color.New(color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n%s '%s' 的搜索结果（第 %s/%s 页，共 %s 条）\n\n",
		bold("Docker Hub"),
		bold(results.Query),
		yellow(fmt.Sprintf("%d", results.Page)),
		yellow(fmt.Sprintf("%d", results.NumPages)),
		yellow(fmt.Sprintf("%d", results.NumResults)))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NAME", "DESCRIPTION", "STARS", "OFFICIAL", "AUTOMATED", "PULLS"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, img := range results.Results {
		official := ""
		if img.IsOfficial {
			official = green("[OK]")
		}
		automated := ""
		if img.IsAutomated {
			automated = cyan("[OK]")
		}
		desc := limitString(img.Description, 60)
		stars := fmt.Sprintf("%d", img.StarCount)
		if img.StarCount > 1000 {
			stars = yellow(stars)
		}
		pulls := formatPullCount(img.PullCount)
		table.Append([]string{
			img.Name,
			desc,
			stars,
			official,
			automated,
			pulls,
		})
	}
	table.Render()
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "server" {
			fmt.Println("Starting Docker Search HTTP server on :8096")
			initTemplates()
			http.HandleFunc("/", handler)
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

			searchResponse := queryImageInfo(imageName, page, limit)
			displayPrettyTable(searchResponse)

			if searchResponse.NumPages > 1 {
				fmt.Printf("\n页码: %d/%d  ", searchResponse.Page, searchResponse.NumPages)
				fmt.Printf("(使用 --page <页码> 参数查看其他页)\n")
			}
		}
	} else {
		fmt.Println("Usage: dockersearch <ImageName> [--page <number>] [--limit <number>] or dockersearch server")
	}
}
