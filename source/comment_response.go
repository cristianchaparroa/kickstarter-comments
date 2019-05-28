package main

import "fmt"

// CommentResponses is the result to retrieve the comments related to the kickstart project
type CommentResponses []struct {
	Data struct {
		Commentable struct {
			ID            string `json:"id"`
			URL           string `json:"url"`
			CanComment    bool   `json:"canComment"`
			CommentsCount int    `json:"commentsCount"`
			Comments      struct {
				Edges []struct {
					Node struct {
						Typename  string      `json:"__typename"`
						ID        string      `json:"id"`
						Body      string      `json:"body"`
						CreatedAt int         `json:"createdAt"`
						ParentID  interface{} `json:"parentId"`
						Author    struct {
							ID       string `json:"id"`
							ImageURL string `json:"imageUrl"`
							Name     string `json:"name"`
							URL      string `json:"url"`
							Typename string `json:"__typename"`
						} `json:"author"`
						AuthorBadges         []string `json:"authorBadges"`
						CanReport            bool     `json:"canReport"`
						CanDelete            bool     `json:"canDelete"`
						HasFlaggings         bool     `json:"hasFlaggings"`
						DeletedAuthor        bool     `json:"deletedAuthor"`
						Deleted              bool     `json:"deleted"`
						AuthorCanceledPledge bool     `json:"authorCanceledPledge"`
						Replies              struct {
							TotalCount int           `json:"totalCount"`
							Nodes      []interface{} `json:"nodes"`
							PageInfo   struct {
								StartCursor     interface{} `json:"startCursor"`
								HasPreviousPage bool        `json:"hasPreviousPage"`
								Typename        string      `json:"__typename"`
							} `json:"pageInfo"`
							Typename string `json:"__typename"`
						} `json:"replies"`
					} `json:"node"`
					Typename string `json:"__typename"`
				} `json:"edges"`
				PageInfo struct {
					StartCursor     string `json:"startCursor"`
					HasNextPage     bool   `json:"hasNextPage"`
					HasPreviousPage bool   `json:"hasPreviousPage"`
					EndCursor       string `json:"endCursor"`
					Typename        string `json:"__typename"`
				} `json:"pageInfo"`
				Typename string `json:"__typename"`
			} `json:"comments"`
			Typename string `json:"__typename"`
		} `json:"commentable"`
		Me interface{} `json:"me"`
	} `json:"data"`
}

// Comment is the siplyfied result from response
type Comment struct {
	URL       string
	ID        string
	Body      string
	CreatedAt int
}

// ToSlice retrieves an slice of string with all the struct  fields
func (c Comment) ToSlice() []string {
	s := make([]string, 0)
	s = append(s, c.URL)
	s = append(s, c.Body)
	s = append(s, fmt.Sprintf("%v", c.CreatedAt))
	s = append(s, c.ID)
	return s
}
