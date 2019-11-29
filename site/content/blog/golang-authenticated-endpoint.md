+++
title = "Adding JWT authentication to a Go API"
description = "A tutorial on securing API endpoints from unauthenticated access by issuing and verifying JWT tokens as well as showing a way of streamlining their testing"
template = "blog/article.html"
slug = "authenticated-endpoint-in-golang"
+++

The security of an API is often ignored in favor of completing the core functionality of our application.
Maybe because we don't have time, maybe because we aren't exactly sure how or maybe because we didn't
remember to do it. In this post we will talk about adding authentication to your API.

Authentication is the process of proving that you are who you say you are (which is different from authorization).
It is often that you want authentication on your API endpoints. 
For example, not everyone should be able to use your API. Or it will come in quite handy if you know exactly who is
using it. This blog post is going to show you how to add an industry-proven authentication mechanism to an API in Go.

Before starting I should mention that this article assumes you are familiar with Go.
If not, you can head over to the [tour of Go](https://tour.golang.org/list) first.

{{ h2(text="A simple web API") }}

We are going to start off with a simple server serving HTTP. We are using HTTP because it's simpler to implement, but
you should always use HTTPS (a really nice [talk on that by Eric Chiang](https://www.youtube.com/watch?v=VwPQKS9Njv0)).
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
  // wire the ListBeers function to handle requests sent to /beer
  http.HandleFunc("/beer", ListBeers)
  // start the HTTP server
  log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```

Simple enough. We define the data type of a `beer` and an `inventory` of two beers.
We start an HTTP server listening on port 8080.
We provide a handler function (`ListBeers`) for that server. Go is going to forward every request 
that comes on port 8080 to the `ListBeers` function we wrote. 
What `ListBeers` does is convert our inventory of beers to JSON so that they can be
sent back to whoever made the request.

You can test it out with `curl` in the terminal: 
```sh
curl localhost:8080/beer
```

This should return (formatted here):
```json
[
  {
    "id": 1,
    "name": "Stella Artois",
    "type": "pilsner",
    "price": 2.5
  },
  {
    "id": 2,
    "name": "Guinness",
    "type": "Irish dry stout",
    "price": 3.5
  }
]
```

The problem with this is that now everyone can send a request and read our list of beers.
We want to make sure that only people we know can see our list of beers.

{{ h2(text="Verifying users") }}

The simplest way of verifying who a user is is by checking their username and password. However, sending these 
constantly increases the risk of them being compromised. And if you needed any additional information about the user
such as their full name or location, you need a separate secure store for that. Queue [**JWT**](https://jwt.io/) tokens.

{{ h3(text="JWT") }}

JWT tokens are short-lived pieces of text that prove the identity of a user
(or provide an assertion, if we had to use proper terms).
Usually an authority such as a central authentication server issues them in exchange for username and password.
In our tutorial we are going to play the role of the authentication server.
JWT tokens allow to verify someone's identity without them showing their password. The way they work is
by cryptographically signing information about the user with a secret key.
The user then sends this token along with every request they make.
After that we can verify the legitimacy of the token and know that the request came from that specific user.

A JWT token looks like this:

```text
eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzQ4ODY0MDcsInVzZXJuYW1lIjoiam9obi5kb2UifQ.OP2rJebLf0Ds7l29M8wY2IcYynFJORJpePa0fUonCPsUkVlonsGWMLHNxtg-S-hpA27EVstruiZxuMGBC6OVzQ
```

It contains three parts separated by dots:

1. A **header** containing information about the cryptographic algorithm that is used to sign the token.

2. **Payload** which is defined by whoever issues the token. 
It typically contains the expiration date of the token, when it was issued, who issued it, username of
the user, etc. But it can contain practically anything.

3. A **signature**; this is the combined payload and header encrypted with the secret key and then hashed.

Each part is then Base64-encoded and all the part are joined by dots so that it's easier to move them around.
Then you get the jumbled up string above.

But if you decode the payload above, you can see what is contains:

```sh
echo -n 'eyJleHAiOjE1NzQ4ODY0MDcsInVzZXJuYW1lIjoiam9obi5kb2UifQ' | base64 --decode
```

In our case it's only the user's name and the expiry time of the token as a Unix timestamp:

```json
{
  "exp": 1574886407,
  "username": "john.doe"
}
```

{{ h2(text="JWT tokens in Go") }}

JWT tokens may seem complex but we will see that integrating them into our beer API isn't going to be difficult.
We are going to use a package for dealing with JWT called [jwt-go](https://github.com/dgrijalva/jwt-go).

First, we are going to look at issuing tokens.

{{ h3(text="Issuing JWT") }}

We will need a place to keep our list of users and their passwords so that we can look them up when a user comes for a token.
For now a `map` should do the job:

```go
var users = map[string]string{
  "john.doe": "Pa55word",
}
```

We also need to define a `struct` type which can hold the information we want in the token about the user.
In JWT terms this is called the claims of the token. We will keep our claims simple. They will
contain only the username and some of the
[standard JWT claims](https://auth0.com/docs/tokens/jwt-claims#reserved-claims) (such as expiry time).

```go
import "github.com/dgrijalva/jwt-go"

// ...

type claims struct {
  // we embed the JWT standard claims for fields like expiry time and audience
  jwt.StandardClaims
  Username string `json:"username"`
}
```

Next, we will add a login endpoint that will receive user's credentials via Basic Auth and issue a JWT token.
I will first show the code first and then go through what it does.

```go
var secretKey = []byte("very-secret-key")

func Login(w http.ResponseWriter, r *http.Request) {
  // grab the basic auth from the request
  username, providedPassword, ok := r.BasicAuth()
  if !ok {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // verify that the user exists and the provided password is correct
  actualPassword, ok := users[username]
  if !ok || actualPassword != providedPassword {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  // create the claims that we are going to issue to our user
  // this will be the payload of our JWT token
  userClaims := claims{
    StandardClaims: jwt.StandardClaims{
      // we give an expiration of 1 hour;
      // the shorter, the more secure
      ExpiresAt: time.Now().Add(time.Hour).Unix(),
    },
    Username: username,
  }

  // create the header and payload of the token
  token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)
  // sign the token with our secret key
  tokenString, err := token.SignedString(secretKey)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // wrap the token in a response so that we can return valid JSON
  response := struct {
    Token string `json:"access_token"`
  }{tokenString}

  // serialize the response into JSON bytes
  serializedToken, _ := json.Marshal(response)

  // send back the token
  w.Write(serializedToken)
}

func main() {
  http.HandleFunc("/beer", ListBeers)
  // add the login endpoint to our server
  http.HandleFunc("/login", Login)
  log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```
<br>

1. We get the username and password that the user sent. If they didn't send any, we reply with 
status 400 (Bad Request).
2. Next, we check if the user actually exists in our records and if their password matches what we have.
If not, we reply with status 401 (Unauthorized)
3. We proceed to construct the claims we will include in the token. We keep them simple by specifying only the expiry
time of the token and the user's username.
4. The next step is to create and sign the JWT token. We create the token by providing the claims and 
the algorithm we want to use for signing it. Then we pass in our secret key and let `jwt-go` do its magic.
5. Now that we have the token we need to pack it in a consumable way. We do the same as we did for 
our beer list.

You can try it out again with `curl`:
```sh
curl localhost:8080/login -u john.doe:Pa55word
```

Which should return something along the lines of:
```json
{"access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzQ4ODYwMjEsInVzZXJuYW1lIjoiam9obi5kb2UifQ.2UTAODiz_zogbQOuoUWkdGWh_DikfGfzw10Z-B6znZbRCl-8uWfsFCHiMSlGqHXFW2E_nzdx4vq0QgRaqya2OQ"}
```

{{ h3(text="Validating JWT") }}

We have a way of issuing tokens but our `/beer` endpoint is still open to everyone.
We will require the user to send their token in the request so that we can verify they
have access to view our list of beers. The convention for sending JWT tokens is that the token
is sent as a `Bearer` token in an `Authorization` header. Sending a `Bearer` token just means that
you need to prepend `Bearer ` to the token.

```go
func ListBeers(w http.ResponseWriter, r *http.Request) {
  // extract the token from the request's headers
  header := r.Header.Get("Authorization")
  header = strings.TrimPrefix(header, "Bearer ")

  // parse the header they sent using our secret key
  token, err := jwt.Parse(header, func(*jwt.Token) (interface{}, error) {
    return secretKey, nil
  })
  if err != nil || !token.Valid {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  bytes, _ := json.Marshal(beers)
  w.Write(bytes)
}
```

Validation has less steps than issuing a token. 

1. Get the request header.
2. Remove the `Bearer ` prefix
3. Pass is to the `jwt-go` package along with our secret key.

Out we get a parsed token. Because we used some of the standard claims, `jwt-go` is able to check
the expiry time of the token for us. If the parsing went OK and if the token hasn't expired yet,
we return the list of beers.

Note that we don't use all the claims we included in our token (namely the username). This is because 
we verify the token using our secret key. It is a secret key, so we know that it's only us who could have issued the token.
And we already check the user during login, so everything is good.
If, for example, the token has a long expiry time and our list of users may change during that time,
then it might be worth it also checking if the username in the token is valid.

You can give it a go yourself. Just make sure to replace `<token>` with a valid token you have from your `/login`:
```sh
curl localhost:8080/beer -H 'Authorization:Bearer <token>'
```

We have successfully added authentication to our API. You can now easily access control your API and be
certain that it is accessed only by the right users.

{{ h3(text="Adding authentication on multiple endpoints") }}

We have the code to check a user token. However, it is mixed with the code that lists beers.
If you want to extend that to many endpoints, you will need to 
either move it to a function that each handler calls or copy it in every handler. None of these are really convenient.
We are going to implement a middleware to solve this. The middleware is just a wrapper handler that 
verifies a token and then delegates to another handler if the token is valid. 

```go
func Authenticated(handler http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // extract the token from the request's headers
    header := r.Header.Get("Authorization")
    header = strings.TrimPrefix(header, "Bearer ")

    // parse the header they sent using our secret key
    token, err := jwt.Parse(header, func(*jwt.Token) (interface{}, error) {
      return secretKey, nil
    })
    if err != nil || !token.Valid {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }

    handler(w, r)
  }
}

func ListBeers(w http.ResponseWriter, r *http.Request) {
  bytes, _ := json.Marshal(beers)
  w.Write(bytes)
}
```

We write a function that takes a function and returns a function.
`Authenticated` takes a handler (`handler`) and wraps it in another handler function.
The wrapper first checks the token from the request and then calls the wrapped handler function (`handler`). 
`Authenticated` then returns the wrapper.

You may notice that the function for listing beers is back to how it
was in the very first example. Now it doesn't have check itself the user that is making the request. 
The middleware ensures that. We also update the `main` to use the middleware:

```go
func main() {
  // wire the ListBeers function to handle requests sent to /beers
  http.HandleFunc("/beer", Authenticated(ListBeers)) 
  // ...
}
```

Using this middleware you can wrap any handler function you have. You will not have to copy code or clutter your business logic
(listing beers in our case) with authentication logic.

You can give it another go with `curl` to make sure it works.

{{ h2(text="Making sure it always works") }}

We built a secure API, open only to certain users. However, now the process of verifying that the API works
isn't that trivial. You need to send a request, get a token, copy the token, send it in a header, then finally verify
that the response is correct. This may become tedious when extending your API.

Another way of testing your API is using [APId](../../). With APId you only need a `.yaml` file with the
endpoints you want to test, what requests you send and what responses you expect.
The [APId CLI tool](../../download) can run those tests any time you need.
To test the authenticated `/beer` endpoint our `apid.yaml` will contain the following (with explanation below):

```yaml
transactions:
  - id: "authenticated-list-beers"
    variables:
      api_url: "http://localhost:8080"
    steps:
      - id: "auth"
        request:
          method: "POST"
          endpoint: "{{ var.api_url }}/login"
          headers:
            Authorization: "Basic {% echo -n 'john.doe:Pa55word' | base64 %}"
        export:
          auth_token: "response.body.access_token"

      - id: "list-beers"
        request:
          method: "GET"
          endpoint: "{{ var.api_url }}/beer"
          headers:
            Authorization: "Bearer {{ auth.auth_token }}"
        expect:
          code: 200
          body:
            exact: false
            type: "json"
            content: |
              [
                {
                  "name": "",
                  "type": "",
                  "price": "",
                  "id": ""
                }
              ]
```

I am going to go through in details the contents explaining what each block does.

This is a single transaction with two steps: one to authenticate (`auth`) and another one to list beers (`list-beers`).

We set a variable for the URL of our API so that we don't have to repeat it in every step (`api_url: "http://localhost:8080"`)

The first step is sending the request for the JWT token. The endpoint that we want to hit would be `http://localhost:8080/login`
but we use a template with the variable we set above. Then we set the Basic Auth
header with the following: `Authorization: "Basic {% echo -n 'john.doe:Pa55word' | base64 %}"`.
This runs a command to base64-encode the credentials so that they are in the expected format.
Then we say that we want to export the auth token from the response. Exporting will set another variable that we can
use for the beers endpoint.

The seconds step makes a request to `http://localhost:8080/beer`. It uses the token we exported just like it would any
other variable. This is the endpoint we want to verify that works. So this step also has an `expect` block.
In the `expect` block we say that we expect status 200 (OK) back from the server.
We also say that we expect the body of the response to be JSON. We don't care about the particular values of the response, we are interested in its structure.
As long as it returns a JSON array of beers, everything is OK. We do this with `exact: false` and giving the
structure of what we want: an array of items, each of which has a `"name"`, `"type"`, `"price"` and `"id"`.
The values of these don't matter so we leave them as empty strings.

You can now save this in a file named `apid.yaml` and run APId to see if it works:
```sh
apid check
```

You can now add this to CI pipelines or schedule its runs to ensure your product is working as expected at all times.

For more details on the syntax of the APId YAML config you can have a glimpse of our [documentation](../../docs).

{{ h2(text="Conclusion") }}

We learned a useful way of adding authentication to our APIs that is considered an industry standard.
We did it in such a way that you can easily modify any existing APIs without touching the core logic
of your endpoints. Then we looked at a convenient way of testing the whole setup with APId.
