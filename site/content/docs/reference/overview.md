+++
title = "Overview"
description = "an overview of the structure of the yaml configuration file"
template = "docs/article.html"
sort_by = weight
weight = 10
+++

{{ h3(text="Transactions and steps") }}

An apid configuration file consists of [transactions](../transactions) which in turn consist of [steps](../steps). Steps 
are the basic elements of the configuration. They specify how to make a request and then how to validate
its response. Transactions bundle steps together to help you represent meaningful stories.

{{ h3(text="Variables") }}

Apid allows you to have [variables](../variables) that will be replaced throughout your steps. Variables can be declared 
for the transaction or step scope or be globally available. They can also come from the environment, which
can be handy for things like injecting secrets and passwords.
