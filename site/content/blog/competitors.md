+++
title = "How is APId different from..."
description = "A compare and contrast with APId alternatives"
template = "blog/article.html"
slug = "comparison"
+++

{{ h2(text="Postman, Newman, Postwoman, Insomnia") }}

These are REST client tools that also allow you to organize requests in some sort of collections.
The Postman "line" of products also offers pre- and post- JS scripts to run with your request.
But any testing needed would have to be implemented in JS.

APId is different in that it allows you to do all of this without having to write a single line of code.
Request and response schemas are laid out in YAML files and all the testing is done by the tool.

Another difference is that APId definitions can be stored anywhere and don't need a special editor.
And you only need a text editor like VS Code or Vim to edit them.
They can also be versioned with your code and be treated the same way you treat your code - 
peer review it, tag it, put it in a pipeline, you name it.

{{ h2(text="MirageJS, WireMock") }}

These are API mocking tools. They help when you want to test your API in isolation from an external API.
You would mock that API so that you can remove anything that is out of your control from your test environment. 

APId is the tool that tests your API. It doesn't care how your API works or who it talks to. It makes sure
that the results make sense and are what you expect them to be.
