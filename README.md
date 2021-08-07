# NewsAppapi
Inshorts API is a basic version of HTTP/JSON api,built using Golang.

## It is capable of the following operations-

### 1.Create an article

POST request using JSON request body. URL: ‘/articles’

### 2.Get an article using id

GET request with Id as url parameter. URL: ‘/articles/{id here}’

### 3.List all articles

GET request.URL: ‘/articles’

### 4.Search for an Article (search in title, subtitle, content)

GET request with the query parameter ‘q’.URL:‘/articles/search?q={search term here}’
  
## Additionally it supports-

- Pagination using offset,limit through URL(/articles?offset={num}&limit={num}
- Validation of input url query.
- Response from list all articles is ordered by timestamp to get latest articles first.

## Data Types
Article looks like this:
```json
{
        "id": "someid",
        "title": "Title of the article",
        "subtitle": "Subtitle",
        "content": "Article/News body",
        "creationtimestamp": "2021-08-07T15:52:43.239166+05:30"
    }
```




