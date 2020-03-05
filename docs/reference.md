# Reference

The APId tests are configured using one or more `YAML` files. These `YAML` files define what checks and what steps each check has to perform.

For more information on how to setup your environment, please follow our [installation guide](https://github.com/getapid/apid-cli/tree/f70eeed52c3849135585cf5ef043d0e293d677ec/installation/README.md).

## Introduction

An APId configuration file consists of variables and transactions which in turn consist of steps. Variables are pieces of data that can be referenced in multiple transactions. Steps are the basic elements of the configuration. They specify how to make a request and then how to validate its response. Transactions bundle steps together to help you represent meaningful stories.

## Variables

APId has variables that can be inplaced throughout your steps and transactions. Variables can be declared for the transaction or step scope or be global. They can also come from the environment, which can be handy for things like secrets and passwords, or from a response from your API.

Variables are scoped either globally, to a transaction or to a step. Variables in a narrower scope have precedence over those in a broader one. Variables are available in templates - `"{{ var.api_url }}"`. Make sure to add quotes to it so that the YAML parser doesn't confuse it for a YAML object.

### Regular variables

These are declared either in the transaction, step or the root yaml document. There they are simply declared as a mapping from variable name to its value. Those will be available in templates using the `var` prefix - `"{{ var.api_url }}"`

### Exported variables

Each step can export a set of variables. This is useful when you want to make a request and then use part of the response in another request. For example, when you authenticate to get a token and then use this token in subsequent requests. See [step](reference.md#step) about the exact syntax of exporting variables in a step. Exported variables will be available in subsequent steps by using the step id that exported them; e.g `"{{ step_one.auth_token }}"`

### Environment variables

These will contain anything environment variable that the APId CLI has inherited. Useful for injecting passwords or other kinds of secrets. They will be available like so: `"{{ env.PASSWORD }}"`

```yaml
variables:
  title: 'A long time ago'
  subtitle: 'in a {{ var.place }} far far away {{ env.DATABASE_USER }} accidentally dropped all tables'
  year: 2187
```

## Commands

Commands are a familiar way to manipulate variables. They allow you to run shell commands and use their output. There are practically no limitations as to what you can do.

The syntax for commands is very similar to the syntax of variables, but instead of using `{{` and `}}` as as delimiters, it uses `{%` and `%}`. For example `{$ echo $ENV_VARIABLE $}`

### Executables

Commands are executed in the default shell \(defined in `$SHELL`\), or `/bin/sh` if none is set.

#### CLI

You can use whatever commands you want, obviously you need to have them set up on you dev machine / build server.

#### Docker

The default docker image of APId is using alpine as the base image, therefore, it has very few executables pre-installed. If you need a more versatile docker image, feel free to build your own.

### Using variables

You can use step and transaction variables from within commands. All the variables are exported for use in commands as `$VAR_CAPITALIZEDNAMEOFVARIABLE`, e.g if you want to use `"{{ step_one.auth_token }}"` in a command, you'd use \`

`. Another example might be, which will be available as`

\`. Note that dashes are replaced with underscores because most shells don't accept dashes inside variable names.

```yaml
steps:
  request:
    endpoint: '{{ var.api_url }}/avengers/{% curl https://dynamic-avengers-api.io/random-avenger-id %}'
```

## Transaction

A transaction is a list of [steps](reference.md#step) which are executed sequentially. If a step fails, the whole transaction fails.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| id | string | yes | A string to uniquely identify a transaction |
| variables | [`variables`](reference.md#variables) | no | Variables scoped to this transaction |
| steps | [`[]step`](reference.md#step) | yes | A list of steps to execute |

```yaml
id: 'transaction-one'
variables:
  api_url: 'https://jsonplaceholder.typicode.com'
steps:
  - id: 'todos-1'
    request:
      method: 'GET'
      endpoint: '{{ var.api_url }}/todos/1'
```

## Step

A step is a call to a single endpoint with optional validation of the response.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| id | string | yes | A string to uniquely identify a step within a transaction |
| variables | [`variables`](reference.md#variables) | no | Variables scoped to this step |
| request | [`request`](reference.md#request) | yes | The request to send |
| expect | [`expect`](reference.md#expect) | no | How to validate the response |
| export | [`export`](reference.md#export) | no | Data to export from this step as variables to be used in other steps |

```yaml
steps:
  - id: 'get user with id 1'
    variables:
      api_url: 'http://localhost:80'
    request:
      method: 'GET'
      endpoint: '{{ var.api_url }}/users/1'
    expect:
      code: 200
      body:
        type: 'json'
        content: |
          {
            "first_name": "Bobby",
            "last_name": "Hounslow",
            "address": {
                "postcode": "TW4 7AE",
                "country": "UK",
            }
          }
    export:
      auth_header: 'response.headers.X-APIDAUTH'
      auth_token: 'response.body.access_token'
```

## Request

Request specifies what request to make - which endpoint to go to, what body to use, etc.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| method | string | yes | The HTTP method of the request |
| endpoint | string | yes | The URL of the request |
| headers | mapping | no | Headers to attach to the request. Keys may repeat. If keys repeat, all the values are added to the header. |
| body | string | no | A string of the body of the request |

```yaml
request:
  method: "GET"
  endpoint: "http://https://jsonplaceholder.typicode.com/posts/1"
  headers:
    X-APIDAUTH: "{{ env.AUTH_TOKEN }}"
    Accept: "application/json"
    Accept: "application/ld+json"
    Accept:
      - "application/json"
      - "application/ld+json"
```

## Expect

Expect will define what we are expecting as a valid response from the API.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| code | int | no | The status code of the response |
| headers | mapping | no | What headers to expect |
| body | [body](reference.md#body) | no | What body to expect |

```yaml
expect:
  code: 200
  headers:
    Accept: 'application/json'
```

## Body

Body provides a bit more flexibility on what body to expect in this response.

When specifying the type of response, the `exact` value has the following behaviour:

If the `type` is `json` and `exact` is:

* `true`: will make sure the JSON content match recursively for every key and value.
* `false`: will check if the keys present in the `body` are also present in the responses body and if their respective values match.

On the other hand, if `type` is `string` and `exact` is:

* `true`: will perform an equals comparison.
* `false`: will check if the provided `body` is a substring of the responses body.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| type | enum | no | The type of the reponse, either `json` or `string` |
| content | string | no | What content of the body to expect |
| exact | bool | no | Is this the entire body, or a part of it |

```yaml
body:
  type: 'json'
  exact: true
  content: |
    {
      "name": "John Doe"
    }
```

