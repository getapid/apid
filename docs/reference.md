# Reference

The APId tests are configured using one or more `YAML` files. These `YAML` files define what checks and what steps each check has to perform.

For more information on how to setup your environemnt, please follow our [installation guide](https://github.com/getapid/apid-cli/tree/c89ab8593880509f8996788a4c2628616cf89aeb/installation/README.md).

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

### Examples

```yaml
variables:
  title: 'A long time ago'
  subtitle: 'in a {{ var.place }} far far away {{ env.DATABASE_USER }} accidentally dropped all tables'
  year: 2187
```

## Commands

Commands are a familiar way to manipulate variables. They allow you to run shell commands and use their output. There are practically no limitations as to what you can do.

The syntax for commands is very similar to the syntax of variables, but instead of using `{{` and `}}` as as delimiters, it uses `{%` and `%}`. For example `{% echo $ENV_VARIABLE %}`

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

### Examples

```yaml
steps:
  request:
    endpoint: '{{ var.api_url }}/avengers/{% curl https://dynamic-avengers-api.io/random-avenger-id %}'
```

