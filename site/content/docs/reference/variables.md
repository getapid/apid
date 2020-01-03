+++
title = "variables"
description = "Variables can be used in templates in the YAML config"
template = "docs/article.html"
sort_by = weight
weight = 220
+++


{{ h2(text="Summary") }}

Variables are scoped either globally, to a transaction or to a step. Variables in a narrower scope have precedence over
those in a broader one. Variables are available in templates - `"{{ var.api_url }}"`. Make sure to add quotes to it
so that the YAML parser doesn't confuse it for a YAML object.

{{ h3(text="Regular variables") }}

These are declared either in the transaction, step or the root yaml document. There they are simply declared as a mapping
 from variable name to its value. Those will be available in templates
using the `var` prefix - `"{{ var.api_url }}"`

{{ h3(text="Exported variables") }}

Each step can export a set of variables. This is useful when you want to make a request and then use part of the response
in another request. For example, when you authenticate to get a token and then use this token in subsequent requests. See
[step](../step) about the exact syntax of exporting variables in a step. Exported variables will be available in
subsequent steps by using the step id that exported them; e.g `"{{ step_one.auth_token }}"`

{{ h3(text="Environment variables") }}

These will contain any environment variable that the APId CLI has inherited. Useful for injecting passwords or
other kinds of secrets. They will be available like so: `"{{ env.PASSWORD }}"`

{{ h3(text="Inside commands") }}

Variables are also available inside commands (`{% %}`). See [commands](../commands) for details.

{{ h2(text="Examples") }}

```yaml
variables:
  title: "A long time ago"
  subtitle: "in a {{ var.place }} far far away {{ env.DATABASE_USER }} accidentally dropped all tables"
  year: 2187
```