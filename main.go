package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

const (
// searchURL = "https://ekb.docdoc.ru/doctor/Bayazitova_Diana"
// searchURL = "https://prodoctorov.ru/ekaterinburg/vrach/919155-bayazitova/"
// platform  = "prodoctorov"
)

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "8000", "Port to listen on")
	client   *req.Client
	reviews  []Review
)

func init() {
	client = req.C().ImpersonateChrome()
}

type Review struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	Message string `json:"message"`
	Source  string `json:"source"`
}

type ApiResponse struct {
	Error int             `json:"error"`
	Data  ApiResponseData `json:"data"`
}

type ApiResponseData struct {
	Code    int      `json:"error,omitempty"`
	Desc    string   `json:"desc,omitempty"`
	Reviews []Review `json:"reviews,omitempty"`
}

func main() {

	// Define routes
	go http.HandleFunc("/api/v1/getReviews", PostHandler)

	// Start the server
	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, nil))

	//fmt.Println(resp)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var resp ApiResponse
	var formError bool

	if r.Method == "POST" {

		r.ParseForm()

		platform := r.PostFormValue("platform")
		doctorUrl := r.PostFormValue("doctorUrl")

		if len(strings.TrimSpace(platform)) == 0 || len(strings.TrimSpace(doctorUrl)) == 0 {
			//log.Fatal(err)
			formError = true
			resp = makeErrorResponse(400, "Неверное значение параметров")

		}

		if !formError {

			reviews := runParser(platform, doctorUrl)

			if reviews != nil {
				resp = makeSuccessResponse(reviews)

			} else {
				resp = makeErrorResponse(400, "Ошибка получения данных")
			}
		}

		json.NewEncoder(w).Encode(resp)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func makeSuccessResponse(reviews []Review) ApiResponse {
	return ApiResponse{
		Error: 0,
		Data: ApiResponseData{
			Reviews: reviews,
		},
	}
}

func makeErrorResponse(code int, desc string) ApiResponse {
	return ApiResponse{
		Error: 1,
		Data: ApiResponseData{
			Code: code,
			Desc: desc,
		},
	}
}

func runParser(platform string, doctorUrl string) []Review {
	reviews = nil

	resp, err := client.R().Get(doctorUrl)

	if err != nil {
		//log.Fatal(err)
		return reviews
	}

	if !resp.IsSuccessState() {
		fmt.Println("bad response status:", resp.Status)
		return reviews
	}

	defer resp.Body.Close()

	return findReviews(resp.Body, platform)
}

func findReviews(result io.Reader, platform string) []Review {

	if platform == "prodoctorov" {
		return parseProdoctorov(result)
	}

	if platform == "sberzdorovie" {
		return parseSberzdorovie(result)
	}

	return reviews
}

func saveResult(result io.Reader) {

	out, err := os.Create("result.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, result)
}

func parseProdoctorov(result io.Reader) []Review {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(result)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".b-review-card").Each(func(i int, s *goquery.Selection) {

		rid, _ := s.Find("div[itemprop='reviewBody']").Attr("data")
		name := s.Find(".b-review-card__author-link").Text()
		date := s.Find("div[itemprop='datePublished']").Text()
		message := s.Find(".b-review-card__comment").Text()

		r := Review{
			ID:      rid,
			Name:    trimAllSpace(name),
			Date:    trimAllSpace(date),
			Message: trimAllSpace(message),
			Source:  "prodoctorov",
		}

		reviews = append(reviews, r)
	})

	return reviews
}

func parseSberzdorovie(result io.Reader) []Review {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(result)
	if err != nil {
		log.Fatal(err)
	}

	//var reviews Review
	res := doc.Find("script[id='__NEXT_DATA__']").Text()

	// data in JSON format which
	// is to be decoded
	Data := []byte(res)

	var parsed map[string]interface{}

	err = json.Unmarshal(Data, &parsed)

	if err != nil {
		fmt.Println(err)
	}

	// Find the review items
	reviewsFromSber := parsed["props"].(map[string]interface{})["pageProps"].(map[string]interface{})["preloadedState"].(map[string]interface{})["doctorPage"].(map[string]interface{})["doctor"].(map[string]interface{})["reviewsForSeo"].([]interface{})

	for _, item := range reviewsFromSber {
		review := item.(map[string]interface{})

		rid := fmt.Sprintf("%.0f", review["id"])

		r := Review{
			ID:      rid,
			Name:    review["name"].(string),
			Date:    review["date"].(string),
			Message: review["text"].(string),
			Source:  "sberzdorovie",
		}

		reviews = append(reviews, r)
	}

	return reviews
}

func trimAllSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
