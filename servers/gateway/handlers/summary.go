package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	// adds an http header to the response
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// gets the url query string paramter value from the request

	form := r.FormValue("url")

	if len(form) == 0 {

		// respond with a http.StatusBadRequest error
		text := http.StatusText(http.StatusBadRequest)
		w.Write([]byte(text))
	}
	w.Header().Add("Content-Type", "application/json")

	htmlStream, err := fetchHTML(form)
	if err == nil {
		pageSum, err := extractSummary(form, htmlStream)

		json.NewEncoder(w).Encode(pageSum)

		if err != nil {
			panic(err)
		}
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	resp, err := http.Get(pageURL)

	if err == nil {
		if resp.StatusCode >= 400 {
			err = errors.New("Bad status code")
			return nil, err
		}

		ctype := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(ctype, "text/html") {
			err = errors.New("Non html page")
			return nil, err
		}
		return resp.Body, nil

	} else {
		return nil, errors.New("404 page not found")
	}
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	pageSum := &PageSummary{}
	if htmlStream == nil {
		return nil, nil
	}
	tokenizer := html.NewTokenizer(htmlStream)
	imgs := []*PreviewImage{}
	foundImage := false
	for {
		tokenType := tokenizer.Next()
		// handles error token / end of file edge case
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}
		// handles title
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()

			if "link" == token.Data {
				img := &PreviewImage{}
				attributes := token.Attr
				for i := 0; i < len(attributes); i++ {
					if attributes[i].Key == "href" {
						link := attributes[i].Val

						u, err := url.Parse(link)
						if err != nil {
							log.Fatal(err)
						}
						base, err := url.Parse(pageURL)
						if err != nil {
							log.Fatal(err)
						}

						img.URL = base.ResolveReference(u).String()

					}
					if attributes[i].Key == "type" {
						img.Type = attributes[i].Val
					}
					if attributes[i].Key == "sizes" && attributes[i].Val != "any" {
						sizes := strings.Split(attributes[i].Val, "x")
						img.Height, _ = strconv.Atoi(sizes[0])
						if len(sizes) == 2 {
							img.Width, _ = strconv.Atoi(sizes[1])
						}
					}
				}
				pageSum.Icon = img
			}
			if "meta" == token.Data {
				attributes := token.Attr
				for i := 0; i < len(attributes); i++ {
					if attributes[i].Val == "keywords" {
						keywordSlice := strings.Split(attributes[i+1].Val, ",")
						for j := 0; j < len(keywordSlice); j++ {
							keywordSlice[j] = strings.TrimSpace(keywordSlice[j])
						}
						pageSum.Keywords = keywordSlice
					}
					if attributes[i].Val == "og:title" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.Title = attributes[j].Val
								println(attributes[j].Val)
							}
						}
					}

					if attributes[i].Val == "og:type" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.Type = attributes[j].Val
							}
						}
					}

					if attributes[i].Val == "og:url" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.URL = attributes[j].Val
							}
						}
					}

					if attributes[i].Val == "og:site_name" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.SiteName = attributes[j].Val
							}
						}
					}

					if attributes[i].Val == "og:description" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.Description = attributes[j].Val
							}
						}
					} else if attributes[i].Val == "description" && pageSum.Description == "" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.Description = attributes[j].Val
							}
						}
					}

					if attributes[i].Val == "author" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								pageSum.Author = attributes[j].Val
							}
						}
					}

					if attributes[i].Val == "og:image" {
						link := ""
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								link = attributes[j].Val
							}
						}
						foundImage = true
						img := &PreviewImage{}

						u, err := url.Parse(link)
						if err != nil {
							log.Fatal(err)
						}
						base, err := url.Parse(pageURL)
						if err != nil {
							log.Fatal(err)
						}

						img.URL = base.ResolveReference(u).String()

						imgs = append(imgs, img)
					}
					if attributes[i].Val == "og:image:secure_url" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								imgs[len(imgs)-1].SecureURL = attributes[j].Val
							}
						}
					}
					if attributes[i].Val == "og:image:type" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								imgs[len(imgs)-1].Type = attributes[j].Val
							}
						}
					}
					if attributes[i].Val == "og:image:width" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								num, err := strconv.Atoi(attributes[j].Val)
								if err == nil {
									imgs[len(imgs)-1].Width = num
								}
							}
						}
					}
					if attributes[i].Val == "og:image:height" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								num, err := strconv.Atoi(attributes[j].Val)
								if err == nil {
									imgs[len(imgs)-1].Height = num
								}
							}
						}
					}
					if attributes[i].Val == "og:image:alt" {
						for j := 0; j < len(attributes); j++ {
							if attributes[j].Key == "content" {
								imgs[len(imgs)-1].Alt = attributes[j].Val
							}
						}
					}
				}
			}
			if "title" == token.Data {
				tokenType = tokenizer.Next()

				if tokenType == html.TextToken {
					if pageSum.Title == "" {
						pageSum.Title = tokenizer.Token().Data
					}
				}
			}

		}

	}
	pageSum.Images = imgs

	if foundImage == false {
		pageSum.Images = nil
	}
	htmlStream.Close()
	return pageSum, nil
}
