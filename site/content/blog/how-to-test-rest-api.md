+++
title = "How to end-to-end test an API"
description = "A quick overview of testing out a REST API"
template = "blog/article.html"
slug = "how-to-test-rest-api"
+++

Testing software is, and always will be, the best way to keep bugs and regression out of your product. REST-ful APIs are no exception to this rule. 

Writing an API is only half the battle. You need to test it as well, otherwise you'll never be sure if it works or not. There are two main categories of testing - white box and black box testing.

White box testing is when you're breaking down the program in pieces and test each of them individually. This category is split into several subcategories - unit and integration testing.

On the other hand, black box testing treats the system as a whole and tests it as is. For example, end-to-end testing is considered black box testing.

{{ h2(text="White box testing") }}

One of the two main ways to validate your API functions as expected is testing the source code directly.

{{ h2(text="Unit testing") }}

It has several different levels of testing. The first one is unit testing. This is where each individual encapsulated block of code, be it a function, or a class / component, is tested in isolation.

Unit testing can be very tricky if the source code is not designed with testability in mind. [This](https://www.guru99.com/unit-testing-guide.html) is a very good article on how to do unit testing the right way!

All popular languages out there have either a build in support for unit testing, or a framework to do so.

For example [unit tests](https://golang.org/pkg/testing/) in GO are first class citizens, where as in Java you need to use a third party framework, like [JUnit](https://junit.org/junit5/).

{{ h2(text="Integration testing") }}

One level above unit testing is integration testing. This type of test makes sure components in your system behave as expected when wired together.

Unlike unit testing, this stage does not test each component in isolation, on the contrary, it tests how it behaves with other components.

Imagine having a component that stores types of beer in a database. An unit test will test wether the component sends the right data to the database without actually sending it. whereas an integration test will send it and check if the data is stored correctly in the database.

Just like unit tests, integration tests can be written with the same frameworks and tools, making them, again, first class citizens.

{{ h2(text="Black box testing") }}

End-to-end testing is taking what an integration test is and expanding it one step further.

Instead of having only one component tested at a time, they test the whole package with all underlying dependencies.

This often times gets quite tricky and most frameworks do not support it out of the box due to different reasons. One of them is because unlike unit and integration testing, end-to-end tests can't be opinionated. Every app has a different set of requirements.

There are two ways you can write end-to-end tests.

The first one is the long and tedious way - write a separate app to test your app. This as you can imagine is quite time consuming and, well, error prone and difficult to maintain.

The other way is to use a declatarive fraemwork for end-to-end testing like [APId](https://www.getapid.com/). Unlike the previous option, using APId is not time consuming, error prone or difficult to maintain.

The only thing you need to provide are the test cases and a version of your stack to run the tests versus. Simple and easy. Learn more about APId [here](https://www.getapid.com/)

{{ h2(text="Conclusion") }}

RESTful API testing comprises of three different stages - unit, integration and end-to-end testing. Each of them covers a different set of functionality.

Together they help you keep your app healthy, add features with confidence and keep your code bug-free.