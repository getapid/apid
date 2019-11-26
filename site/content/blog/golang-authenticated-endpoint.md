+++
title = "How to add JWT authentication to an API in Golang"
description = "A tutorial on securing API endpoints from unauthenticated access by verifying JWT tokens"
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

var inventory = []beer{
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
}

// ListBeers takes the inventory of beers, converts them to JSON and returns them
func ListBeers(w http.ResponseWriter, r *http.Request) {
  jsonBytes, _ := json.Marshal(inventory)
  w.Write(jsonBytes)
}

func main() {
  // wire the ListBeers function to handle requests sent to /beers
  http.HandleFunc("/beers", ListBeers)
  // start the HTTP server
  log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```

Simple enough. We started an HTTP server listening on port 8080.
We provided a handler function (`ListBeers`) for that server. Go is going to forward every request 
that comes on port 8080 to the `ListBeers` function we wrote. 
What `ListBeers` does is convert its inventory of beers to JSON so that they can be
sent back to whoever made the request.

As we already mentioned, the problem with this is that now everyone can send a request and read our list of beers.
But we don't want that. We want to make sure that only people we know can see our list of beers.

{{ h2(text="Verifying users") }}

The simplest way of verifying who a user is is by checking their username and password. However, sending these 
constantly increases the risk of them being compromised. Queue [**JWT**](https://jwt.io/) tokens.

{{ h3(text="JWT") }}

JWT tokens allow to verify someones identity without them showing you their password. The way they work is
by cryptographically signing information about the user. The user then sends this token along with every
request they make. This way we know that the request came from them.

A JWT token looks like this:

```text
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

It contains three parts separated by dots: a **header**, **payload**, and **signature**.

The **header** contains information about the cryptographic algorithm that is used to sign the token.

The **payload** is defined by whoever issues the token. This means that it can contain practically anything.

The **signature** is the combined payload and header encrypted with a secret key and then hashed. Each part
is then Base64-encoded and all the part are appended to each other so that it's easier to move around.
Then you get the string above.

If you decode the payload above and actually see what is contains:

```sh
echo 'eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ' | base64 --decode
```

In our case it contains only the user's name, their id in the subject(`sub`) field and the time the token was issued
in a Unix timestamp in seconds:

```json
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022
}
```

{{ h2(text="JWT tokens in Go") }}

JWT tokens may seem complex but we will see that actually integrating them into our beer API isn't that hard. 
First, we are going to look at issuing tokens.

{{ h3(text="Issuing JWT") }}

We will need a place to keep our list of users and their passwords so that we can verify them.
For now a `map` should do that job. We will have a single user:

```go
var users = map[string]string{
  "john.doe": "Pa55word",
}

type beer struct {
  Id    int     `json:"id"`
...
```

Next, we will add an endpoint that will receive user's credentials via Basic Auth and issue a JWT token.

```go
func Login(w http.ResponseWriter, r *http.Request) {
  // grab the basic auth from the request 
  username, password, ok := r.BasicAuth()
  if !ok {
    w.WriteHeader(http.Unauthorized)
  }  
}

func main() {
  // wire the ListBeers and Login functions
  http.HandleFunc("/beers", ListBeers)
  http.HandleFunc("/login", Login)
  // start the HTTP server
  log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```
