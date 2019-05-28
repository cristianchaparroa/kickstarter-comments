package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	dataframe "github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
)

// ICommentsManager defines the mehtods to load the dataset
type ICommentsManager interface {
	GetDataframe() dataframe.DataFrame

	GetURLs() []string

	WriteCSV(cs []Comment)
}

// CommentsManager is in charge to load the project dataset
type CommentsManager struct {
	DatasetPath string
}

// NewCommentsManager generates a pointer to CommentsManager
func NewCommentsManager(path string) *CommentsManager {
	return &CommentsManager{DatasetPath: path}
}

// GetDataframe retrieves the records with comments
func (l *CommentsManager) GetDataframe() dataframe.DataFrame {
	content, _ := ioutil.ReadFile(l.DatasetPath)
	ioContent := strings.NewReader(string(content))

	df := dataframe.ReadCSV(ioContent)

	rows, _ := df.Dims()
	fmt.Printf("Number of projects: %#v\n", rows)

	// 2. Filter the projects with comments.
	filter := dataframe.F{Colname: "Comments", Comparator: series.Greater, Comparando: 0}
	dff := df.Filter(filter)

	return dff
}

// GetURLs retrieve urls from dataframe
func (l *CommentsManager) GetURLs(dff dataframe.DataFrame) []string {
	rows, _ := dff.Dims()
	urls := make([]string, 0)
	urlIndex := 2
	for i := 0; i < rows; i++ {
		url := dff.Elem(i, urlIndex)
		urls = append(urls, url.String())
	}
	return urls
}

// WriteCSV write a csv file with comments
func (l *CommentsManager) WriteCSV(cs []Comment) {
	fmt.Println("--> WriteCSV")
	data := make([][]string, 0)

	hs := []string{"project_url", "comment", "created_at", "comment_id"}

	data = append(data, hs)

	for _, c := range cs {
		s := c.ToSlice()
		data = append(data, s)
	}

	fileName := "comments.csv"
	file, err := os.Create(fileName)

	if err != nil {
		fmt.Printf("Error trying to create the file .csv, :%s", err)
		return
	}

	defer file.Close()

	w := csv.NewWriter(file)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		fmt.Printf("error writing csv:%s\n", err)
	}

	fmt.Println("<-- WriteCSV")
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
