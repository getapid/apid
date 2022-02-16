# User Focused API testing

## 🔭 What is APId?

APId is a framework that lets you write declarative, end-to-end collections of requests and make sure your API behaves the way you expect.

## ⬇️ Installation

APId comes in both binary packages and docker image. You can find the docker image [here](https://github.com/getapid/apid/pkgs/container/apid), while the binaries can be found [here](https://github.com/getapid/apid/releases)

Here's how to install the latest binary on UNIX based systems:

```sh
# make sure to substitute the URL with the correct platform for you
curl -L https://github.com/getapid/apid/releases/latest/download/apid-darwin-arm64 -o /tmp/apid
chmod +x /tmp/apid
sudo mv /tmp/apid /usr/local/bin/apid

# test if the installation was successful 
apid version
```

### ✅ A simple test

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

And to run the test

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

For more examples please check the [`examples`](examples) folder in this repository.

## 📚 Documentation

You can find the most up to date documentation [`here`](docs)

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


## ⁉️ Help

If you have any questions or ideas please feel free to head over to the [discussion](https://github.com/getapid/apid/discussions) tab in this repository and ask away!

### 💻 CLI

1. Head to the [latest release](https://github.com/getapid/apid/releases/latest) and select the binary for your operating system.
2. Once downloaded open the archive and place the executable in a directory on your `$PATH`.

## 👽 Contributing

To contribute to APId, please see [CONTRIBUTING](CONTRIBUTING.md).
