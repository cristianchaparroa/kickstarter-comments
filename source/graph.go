package main

const (

	// GraphQuery is the query used to retrieve the comments from Kikstart
	GraphQuery = `query ($commentableId: ID!, $nextCursor: String, $previousCursor: String, $replyCursor: String, $first: Int, $last: Int) {
  commentable: node(id: $commentableId) {
    ... on Project {
      id
      url
      canComment
      commentsCount
      comments(first: $first, last: $last, after: $nextCursor, before: $previousCursor) {
        edges {
          node {
            ...CommentInfo
            ...CommentReplies
            __typename
          }
          __typename
        }
        pageInfo {
          startCursor
          hasNextPage
          hasPreviousPage
          endCursor
          __typename
        }
        __typename
      }
      __typename
    }
    __typename
  }
  me {
    id
    name
    imageUrl(width: 200)
    isKsrAdmin
    url
    __typename
  }
}

fragment CommentInfo on ProjectComment {
  id
  body
  createdAt
  parentId
  author {
    id
    imageUrl(width: 200)
    name
    url
    __typename
  }
  authorBadges
  canReport
  canDelete
  hasFlaggings
  deletedAuthor
  deleted
  authorCanceledPledge
  __typename
}

fragment CommentReplies on ProjectComment {
  replies(last: 3, before: $replyCursor) {
    totalCount
    nodes {
      ...CommentInfo
      __typename
    }
    pageInfo {
      startCursor
      hasPreviousPage
      __typename
    }
    __typename
  }
  __typename
}`
)

// GraphVariable ...
type GraphVariable struct {
	CommentableID string `json:"commentableId"`
}

// GraphPayload is the payload to retrieve the baker comments related to
// project
type GraphPayload struct {
	Query     string        `json:"query"`
	Variables GraphVariable `json:"variables"`
}
