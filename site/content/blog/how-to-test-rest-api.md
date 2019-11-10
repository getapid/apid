+++
title = "How to end-to-end test an API"
description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc eu feugiat sapien. Aenean ligula nunc, laoreet id sem in, interdum bibendum felis. Donec vel dui neque. Praesent ac sem ut justo volutpat rutrum a imperdiet tellus."
template = "blog/article.html"
slug = "how-to-test-rest-api"
+++

{{ h2(text="REST API from ten thousand feet") }}

Most modern softwares today use what is called are front-end and back-end systems. Typically the front-end is used to display data and handle user input, where as the back-end is mostly used to store the said data and manipulate it in certain ways. 


There are many ways data can flow from the front- to the back-end, however the web world nowadays is primarily dominated by HTTP, a hyper text transfer protocol, or simply put a protocol on how to structure and send data over the internet. A REST-ful API in turn is another "protocol", or an ideology, on how to structure data one level ontop of HTTP. 


It uses strict naming conventions and resource locations to bring order to the chaos. For more information on what those conventions are click [here](https://en.wikipedia.org/wiki/Representational_state_transfer).

{{ h2(text="Types of testing") }}

In most ways, a RESTful API is no different than testing an ordinary software product. 

{{ h3(text="Unit testing") }}

It has several different levels of testing. The first one is unit testing. This is where each individual encapsulated block of code, be it a function, or a class / component, is tested in isolation. 

Unit testing can be very tricky if the source code is not designed with testability in mind. [This](https://www.guru99.com/unit-testing-guide.html) is a very good article on how to do unit testing the right way!

All popular languages out there have either a build in support for unit testing, or a framework to do so.

For example [unit tests](https://golang.org/pkg/testing/) in GO are first class citizens, where as in Java you need to use a third party framework, like [JUnit](https://junit.org/junit5/).

{{ h3(text="Integration testing") }}

One level above unit testing is integration testing. This type of test makes sure components in your system behave as expected when wired together. 

Unlike unit testing, this stage does not test each component in isolation, on the contrary, it tests how it behaves with other components.

Imagine having a component that stores types of beer in a database. An unit test will test wether the component sends the right data to the database without actually sending it. whereas an integration test will send it and check if the data is stored correctly in the database.

Just like unit tests, integration tests can be written with the same frameworks and tools, making them, again, first class citizens.

{{ h3(text="End-to-end testing") }}

End-to-end testing is taking what an integration test is and expanding it one step further. 

Instead of having only one component tested at a time, they test the whole package with all underlying dependencies.

This often times gets quite tricky and most frameworks do not support it out of the box due to different reasons. One of them is because unlike unit and integration testing, end-to-end tests can't be opinionated. Every app has a different set of requirements.

There are two ways you can write end-to-end tests. 

The first one is the long and tedious way - write a separate app to test your app. This as you can imagine is quite time consuming and, well, error prone and difficult to maintain. 

The other way is to use a declatarive fraemwork for end-to-end testing like [APId](/). Unlike the previous option, using APId is not time consuming, error prone or difficult to maintain. 

The only thing you need to provide are the test cases and a version of your stack to run the tests versus. Simple and easy. Learn more about APId [here](/)

{{ h2(text="Conclusion") }}

RESTful API testing comprises of three different stages - unit, integration and end-to-end testing. Each of them covers a different set of functionality.

Together they help you keep your app healthy, add features with confidence and keep your code bug-free.