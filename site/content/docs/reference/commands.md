+++
title = "commands"
description = "Using shell commands in templates"
template = "docs/article.html"
sort_by = weight
weight = 30
+++


{{ h2(text="Summary") }}

Commands are a familiar way to manipulate variables. They allow you to run shell commands and use their output. 
There are practically no limitations as to what you can do.

The syntax for commands is very similar to the syntax of variables, but instead of using `{{` and `}}` as as delimiters, 
it uses `{%` and `%}`. For example `{% echo $ENV_VARIABLE %}`

{{ h2(text="Executables") }}

Commands are executed in the default shell (defined in `$SHELL`), or `/bin/sh` if none is set.

{{ h3(text="CLI") }}

You can use whatever commands you want, obviously you need to have them set up on you dev machine / build server.

{{ h3(text="Docker") }}

The default docker image of APId is using alpine as the base image, therefore, it has very few executables pre-installed.
If you need a more versatile docker image, feel free to build your own.

{{ h2(text="Using variables") }}

You can use step and transaction [variables](./variables) from within commands. All the variables are exported 
for use in commands as `$VAR_CAPITALIZEDNAMEOFVARIABLE`, e.g if you want to use `"{{ step_one.auth_token }}"` in a command, you'd
use `{% echo $STEP_ONE_AUTH_TOKEN %}`. Another example might be `{{ var.my-name }}`, which will be available as 
`{% VAR_MY_NAME %}`. Note that dashes are replaced with underscores because most shells don't accept dashes inside variable names.

{{ h2(text="Examples") }}

```yaml
steps:
  request:
    endpoint: "{{ var.api_url }}/avengers/{% curl https://dynamic-avengers-api.io/random-avenger-id %}"
```