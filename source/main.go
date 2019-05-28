package main

import (
	"fmt"
	"sync"
	"time"
)

// Custom user agent.
const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/53.0.2785.143 " +
		"Safari/537.36"
)

// Process contains the information related to the current process
type Process struct {
	Comments []Comment
	Message  string
	Index    int
}

func main() {
	cm := NewCommentsManager("./dataset/projects.csv")
	dff := cm.GetDataframe()
	urls := cm.GetURLs(dff)

	processAll(urls, cm)
}

func processAll(urls []string, cm *CommentsManager) {
	// 4. Fetch the url's concurrency
	sem := make(chan struct{}, 100)
	ch := make(chan Process)

	//comments := make(chan []Comment)

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(urls))

	/*
		for i := 0; i < 1; i++ {
			url := urls[i]*/
	for i, url := range urls {

		go func(index int, u string, ch chan<- Process) {

			// if there are already 100 goroutines running(buffered channel),
			// below send will block and a new file wont be open
			sem <- struct{}{}

			// once this goroutine finishes, empty the buffer by one
			// so the next process may start (another goroutine blocked on
			// above send will now be able to execute the statement and continue)
			defer func() { <-sem }()

			// wg.Done must be deferred after a read from sem so that
			// it executes before the above read
			defer wg.Done()

			fetch(index, u, ch)
			// handle open file
		}(i, url, ch)
	}

	cms := make([]Comment, 0)
	for range urls {
		p := <-ch
		cms = append(cms, p.Comments...)
		fmt.Printf("%v %s \n", p.Index, p.Message)
	}

	wg.Wait()
	close(sem)
	close(ch)

	cm.WriteCSV(cms)

	secs := time.Since(start).Seconds()
	fmt.Printf("%.2fs  elapsed\n", secs)
}

// fetch explore a specific url an extract the comments from there.
func fetch(i int, url string, ch chan<- Process) {
	fmt.Printf("--> fetch: %s\n", url)
	start := time.Now()
	sp := NewCommentScraper()
	resp, err := sp.ProjectRequest(url)

	if err != nil {
		m := fmt.Sprintf("Error to make a project request %s %v", url, err)
		p := Process{Message: m}
		ch <- p
		return
	}

	defer resp.Body.Close()

	cResp, err := sp.CommentRequest(url, resp)

	if err != nil {
		m := fmt.Sprintf("while get the comments from comments %s %v", url, err)
		p := Process{Message: m, Index: i}
		ch <- p
		return
	}

	defer cResp.Body.Close()

	commentableID := sp.GetCommentableID(cResp.Body)
	fmt.Println(commentableID)

	gResp, err := sp.GraphRequest(commentableID, url)

	if err != nil {
		m := fmt.Sprintf("while get the comments from graph %s %v", url, err)
		p := Process{Message: m, Index: i}
		ch <- p
		return
	}
	defer gResp.Body.Close()

	cs, _ := sp.GetComments(url, gResp.Body)

	secs := time.Since(start).Seconds()
	m := fmt.Sprintf("%.2fs   %s", secs, url)
	p := Process{Message: m, Comments: cs, Index: i}
	ch <- p
	fmt.Println("<-- fetch")
}
