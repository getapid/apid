+++
title = "How to add authentication to a endpoint in Golang"
description = "A tutorial on securing an endpoint from unauthenticated access"
template = "blog/article.html"
slug = "authenticated-endpoint-in-golang"
+++

Authentication is the process of proving that you are who you say you are (which is different from authorization).
It is often that you want authentication on your endpoints. 
Maybe not everyone should be able to use your API or maybe you really care who it is that is using it.
These are cases where you need to make sure you protect your API from unauthenticated access.
I am going to show you how to do this in Golang.

Before starting I should mention that this article assumes you are familiar with Go.
If not, you can head over to the [tour of Go](https://tour.golang.org/list).

{{ h2(text="A simple web API") }}

We are going to start off with a simple server serving HTTP (we are using HTTP because it's simpler to implement;
you should always use HTTPS; a really nice [talk on that by Eric Chiang](https://www.youtube.com/watch?v=VwPQKS9Njv0)).
Because we at APId like beer, that server is going to serve different brands of beer (no pun intended).
Without further ado, here's the code:

```go
type beer struct {
  Id    int     `json:"id"`
  Name  string  `json:"name"`
  Type  string  `json:"type"`
  Price float32 `json:"price"`
}

type handler struct {
  inventory []beer
}

func (s handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  bytes, _ := json.Marshal(s.inventory)
  w.Write(bytes)
}

func main() {
  // create the handler with an inventory of two beers
  beerServer := handler{
    inventory: []beer{
      {
        Id:    1,
        Name:  "Stella Artois",
        Type:  "pilsner",
        Price: 2.50,
      },
      {
        Id:    2,
        Name:  "Guinness",
        Type:  "Irish dry stout",
        Price: 3.50,
      },
    },
  }
  
  // start the HTTP server
  log.Fatal(http.ListenAndServe("localhost:8080", beerServer))
}
```

Simple enough. We started an HTTP server listening on port 8080.
We provided a handler (`beerServer`) for that server. Go is going to forward every request 
that comes on port 8080 to the `ServeHTTP` function we wrote. 
What `ServeHTTP` does is convert its inventory of beers to a slice of bytes so that they can be
sent back to whoever made the request.

{{ h3(text="Open to everyone") }}

As we already mentioned, the problem with this is that now everyone can send a request and read our list of beers.
But we don't want that. We want to make sure that only people we know can see our beers.
