package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/antchfx/htmlquery"
)

// ICommentScraper defines the methods to scrap the comments related with a project
type ICommentScraper interface {
	ProjectRequest(url string) (*http.Response, error)
	CommentRequest(url string, resp *http.Response) (*http.Response, error)
	GraphRequest(commentableID, url string) (*http.Response, error)
	GetCommentableID(body io.ReadCloser) string
	GetComments(url string, body io.ReadCloser) ([]Comment, error)
}

// CommentScrapper implements all methods to scrap the comments
type CommentScrapper struct {
}

// NewCommentScraper returns a pointer to CommentScrapper
func NewCommentScraper() *CommentScrapper {
	return &CommentScrapper{}
}

func (s *CommentScrapper) ProjectRequest(url string) (*http.Response, error) {
	fmt.Println("--> ProjectRequest")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error:%#v\n", err)
		return nil, err
	}
	fmt.Println("<-- ProjectRequest")

	return resp, nil
}

func (s *CommentScrapper) CommentRequest(url string, resp *http.Response) (*http.Response, error) {
	fmt.Println("--> CommentRequest")
	client := &http.Client{}

	now := strconv.FormatInt(MakeTimestampMilli(), 10)
	commentsURL := fmt.Sprintf("%s/comments?_=%s", url, now)
	req, _ := http.NewRequest("GET", commentsURL, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9,fr;q=0.8,es-419;q=0.7,es;q=0.6,en-US;q=0.5,gl;q=0.4,pt;q=0.3")
	req.Header.Set("referer", fmt.Sprintf("%s/comments", url))

	fmt.Printf("--> URL: %s \n", commentsURL)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return nil, err
	}

	fmt.Println("<-- CommentRequest")
	return resp, nil
}

func (s *CommentScrapper) GetCommentableID(body io.ReadCloser) string {
	fmt.Println("--> CommentRequest")
	doc, err := htmlquery.Parse(body)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	node := htmlquery.FindOne(doc, "//div[@id='react-project-comments']")

	if node == nil || node.Attr == nil {
		fmt.Println("<-- GetCommentableID, node nil or attributes")
		return ""
	}

	for _, att := range node.Attr {
		if att.Key == "data-commentable_id" {
			fmt.Printf("<-- GetCommentableID:%s\n", att.Val)
			return att.Val
		}
	}
	fmt.Println("<-- GetCommentableID")
	return ""
}

func (s *CommentScrapper) GraphRequest(commentableID, url string) (*http.Response, error) {
	fmt.Println("--> GraphRequest")
	client := &http.Client{}

	q := GraphQuery
	gv := GraphVariable{CommentableID: commentableID}
	gp := GraphPayload{Variables: gv, Query: q}
	payload := make([]GraphPayload, 0)
	payload = append(payload, gp)

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		fmt.Printf("Error:%#v\n", err)
		fmt.Println("<-- GraphRequest")
		return nil, err
	}

	req, _ := http.NewRequest("POST", "https://www.kickstarter.com/graph", bytes.NewBuffer(jsonPayload))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB,en;q=0.9,fr;q=0.8,es-419;q=0.7,es;q=0.6,en-US;q=0.5,gl;q=0.4,pt;q=0.3")
	req.Header.Set("referer", fmt.Sprintf("%s/comments", url))
	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error making the request: %s\n", err)
		fmt.Println("<-- GraphRequest")
		return nil, err
	}

	fmt.Println("<-- GraphRequest")
	return resp, nil
}

func (s *CommentScrapper) GetComments(url string, body io.ReadCloser) ([]Comment, error) {
	fmt.Println("--> GetComments")

	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var commentResponses CommentResponses
	json.Unmarshal(bodyBytes, &commentResponses)

	cs := make([]Comment, 0)

	for _, r := range commentResponses {
		es := r.Data.Commentable.Comments.Edges

		for _, e := range es {
			n := e.Node
			c := Comment{ID: n.ID, Body: n.Body, CreatedAt: n.CreatedAt, URL: url}
			cs = append(cs, c)
		}
	}

	fmt.Println("<-- GetComments")

	return cs, nil
}
