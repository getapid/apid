# 🔭 APId documentation

APId is a framework that lets you write declarative, end-to-end collections of requests and make sure your API behaves the way you expect.

## ⬇️ Installation

APId comes in both binary packages and docker image. You can find the docker image [here](https://hub.docker.com/r/getapid/apid), while the binaries can be found [here](https://github.com/getapid/apid/releases)

Here's how to install the latest binary on UNIX based systems:

```sh
# make sure to substitute the URL with the correct platform for you
curl -L https://github.com/getapid/apid/releases/latest/download/apid-darwin-arm64 -o /tmp/apid
chmod +x /tmp/apid
sudo mv /tmp/apid /usr/local/bin/apid

# test if the installation was successful 
apid version
```

## ✅ A simple test

APId tests, or specs, are written in `jsonnet`. There are a number of built-in useful functions to make it easier to make and validate requests to your API.

```jsonnet
// contents of `example.jsonnet`

{
  simple_spec: spec([
    {
      name: "google homepage",
      request: {
        method: "GET",
        url: "https://www.google.com/"
      },
      expect: {
        code: 200
      }
    }
  ])
}
```

To run the test, issue

```bash
> apid check -s "example.jsonnet"

example::simple_spec
    google homepage
        + status code is 200

specs passed: 1
specs failed: 0
```

Success! You've just written your first APId test! If you change the `expect.code` from `200` to lets say `500` the test will fail and this will be the output:

```bash
> apid check -s "example.jsonnet"

example::simple_spec
    google homepage
        o status code: wanted 500, got 200  

specs passed: 0
specs failed: 1
```

## Structure

APId comes with a list of helpful functions that let you define what to expect from each response. Before looking into that, lets see what's the basic structure of a spec file.

Spec files are written in [jsonnet](https://jsonnet.org/). There is a helper function to define a spec, conveniently named `spec`.

A basic spec file returns a json object where each key is the name of the spec and it's value is a `spec`. For example

```jsonnet
{
  spec_name: spec([])
}
```

The `spec` helper function takes a list of steps as a parameter. Those steps are executed sequentially, letting you model how your users interact with your API. If one step fails the rest won't be executed.

## Steps

A step represents a single API call. It is defined as a json object, for example

```js
{
    name: 'a descriptive identifier for this step',
    request: {
        type: 'GET',
        url: 'https://www.google.com/',
        headers: {
            'haeder-name': 'header-value'
        },
        body: {
            'json body': 'body can also be just a simple string'
        }
    }
    expect: {
        code: 200,
        headers: {
            'header-name': 'header-value'
        },
        body: 'expect a string body'
    }
}
```

This is pretty self explanatory, the only non-obvious thing might be that the `body` in both `request` and `expect` can be of any value - object, array, string, float, etc.

## Validation

###  Matcher Translation

Matchers are a very versatile way of checking what value you got back. There is a list of matchers below, but before we get to them lets see how they work. APId transforms all keys and values in `expect.body`, `expect.headers` and `expect.body` blocks to matchers. This means that 

```js
{
    code: 200
}
```

is the same as writing

```js
{
    // more what float is below
    code: float(200)
}
```

APId implicitly transforms raw values to matchers the following way

| JSON type | Matcher  |
| --------- | -------- |
| Object    | `json`   |
| Number    | `float`  |
| String    | `string` |
| Array     | `array`  |

If you want to enforce checks for a specific type you can manually specify which checker to use.

### Matchers

Here are all the matchers you can use and what parameters they take. The matchers are provided in the form `function(param: type = default_value)`

> Please note all matchers in this table are of type `matcher`

| Matcher                                              | Description                                                                                                                       | Example                                  |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------- |
| `any()`                                              | Matches any value. Use this when you want to check for the existence of a key or value                                            | `any()`                                  |
| `string(value: string, case_sensitive: bool = true)` | Matches a string                                                                                                                  | `string("a string value")`               |
| `int(value: int)`                                    | Checks if the value is an int with the provided value                                                                             | `int(200)`                               |
| `float(value: float)`                                | Checks if the value is an float with the provided value                                                                           | `float(88.36)`                           |
| `regex(regex: string)`                               | Checks if a string matches the provided regex                                                                                     | `regex("\\w+")`                          |
| `json(object: map, subset: bool = false)`            | Checks if the value is an object. When subset is false, the received value can have extra keys not present in the provided object | `json({ some: "value" })`                |
| `array(array: map, subset: bool = false)`            | Checks if the value is an array. When subset is false, the received value can have extra values not present in the provided array | `array(["value", { another: "value" }])` |
| `len(length: int)`                                   | Checks if the length of the value matches. Can be used on `string`, `object` and `array`, otherwise fails                         | `len(3)`                                 |
| `range(from: float, to: float)`                      | Checks if the value is more than or equal to `from` and less than or equal to `to`                                                | `range(3.0, 8.0)`                        |

### Boolean matchers

There are two extra matchers provided for complex situations. These are the `and` and `or` matchers.

> Please note all matchers in this table are of type `matcher` allowing you to nest them indefinitely

| Matcher                    | Description                                                  | Example                              |
| -------------------------- | ------------------------------------------------------------ | ------------------------------------ |
| `and(matchers: []matcher)` | Checks if the value matches all provided matchers            | `all([ type.int, range(3.0, 8.0) ])` |
| `or(matchers: []matcher)`  | Checks if the value matches any one of the provided matchers | `or([ type.int, range(3.0, 8.0) ])`  |

Writing complex matchers

```js
// With the boolean matchers you can write something like
body: {
    key: and([
        type.object,
        len(3),
        or([
            {
                nested_key: regex("\\w+")
            },
            {
                nested_key: type.int
            }
        ])
    ])
}
```

The example above would pass only if the value of `key` is an object with three keys, one of which has a key with value `nested_key` and is either matching `\w+` or is an int. This might not be the best use of complex matchers, but it shows you how powerful they are.

### Key matchers

JSON keys are strings. In most cases it's more than enough to do an `equals` comparison, but in some cases you might want to check if there is a key that matches a specific regex for example. To define a key matcher the only thing you need to do is encapsulate the matchers you want in a `key()`.

| Matcher                 | Description                                  | Example             |
| ----------------------- | -------------------------------------------- | ------------------- |
| `key(matcher: matcher)` | Checks if a key matches the provided matcher | `key(regex("\w+"))` |

A matcher is any valid matcher, though some don't make sense to be used here e.g. you can't have an object as a key, but `key(json({ key: "value" }))` is a valid matcher. It will always fail, but won't cause compile issues.

An example of a complex key matcher would be

```js
{
    body: {
        key(
            or([
                regex("\\w+"),
                regex("\\d+"),
            ])
        ): "the value of that key"
    }
}
```

### Typechecks

If you don't care about the value, you can just check if a certain filed is of a certain type. APId provides a `types` object that has basic type matchers. For example

> Please note type checkers are not functions, but constants instead!

```js
{
    // check that the value is a float with value `467` (automatically casts ints to floats when checking)
    body: {
        'key': 467
    }
}
```

```js
{
    // will check that the value is an integer with value `467`
    body: {
        'key': 467
    }
}
```

```js
{
    // will check that the value is an integer and ignore the value
    body: {
        'key': type.int
    }
}
```

Here is a list of the type matchers available 

| Matcher       |
| ------------- |
| `type.int`    |
| `type.float`  |
| `type.string` |
| `type.object` |
| `type.array`  |

## Patterns

Jsonnet is a very powerful language which can be utilised to make your life easier.

### Split variables from tests

For example you can extract any variables in a separate file

```js
// vars.libsonnet
{
  url: 'http://localhost:8080',
}

// test.jsonnet
{
    name: 'request',
    request: {
        url: vars.url,
    },
    expect: {
        code: 200,
    },
},
```

### Store matchers in variables

You can extract your matchers in a local variable to make the test easier to read

```js
// test.jsonnet
local key_matcher = key(
    or([
        regex("\\w+"),
        regex("\\d+"),
    ])
);

{
    body: {
        [key_matcher]: "the value of that key" // note the [] around the key, ref: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Object_initializer#computed_property_names
    }
}
```
