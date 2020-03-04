# commands

## Summary

Commands are a familiar way to manipulate variables. They allow you to run shell commands and use their output. There are practically no limitations as to what you can do.

The syntax for commands is very similar to the syntax of variables, but instead of using `{{` and `}}` as as delimiters, it uses \`

`. For example`

\`

## Executables

Commands are executed in the default shell \(defined in `$SHELL`\), or `/bin/sh` if none is set.

### CLI

You can use whatever commands you want, obviously you need to have them set up on you dev machine / build server.

### Docker

The default docker image of APId is using alpine as the base image, thus it has very few executables pre-installed. If you need a more versatile docker image, feel free to build your own.

## Using variables

In certain cases, one might want to access [variables](reference/variables/README.md) from within commands. All the variables are exported for use in commands as `$VAR_CAPITALIZEDNAMEOFVARIABLE`, e.g if you want to use `"{{ step_one.auth_token }}"` in a command, you'd use \`

`. Another example might be, which will be available as`

\`.

## Examples

```yaml
steps:
  request:
    endpoint: '{{ var.api_url }}/avengers/{% curl https://dynamic-avengers-api.io/random-avenger-id %}'
```
