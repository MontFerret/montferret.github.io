---
title: "Say 'Hello' to Ferret"
subtitle: "A modern web scraping tool"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-02-08"
---

As a developer, you've probably created a couple of web scrapers to get some data from the Internet during your career. That data could be images for an AI project or hotel prices from booking websites for a new travel mobile app or something else. But something that needs data.

It's a relatively easy task to write one or two web scrapers if the amount of required data is small. But it's a completely different level of complexity when you need to scrape big amount of data and, moreover, do it frequently.

This is what happened to me. A while ago, I started a side project which was heavily data driven. The project required to gather data about cosmetological brands and its products, the ingredients of these products and their effect on human's skin. Not huge, but pretty large data set.
The first impulse was to write bunch of scrapers using some scraping framework. And this is what was done. After a while, I started facing with issues I'm going to talk about next and came up with a solution you will see after.

## The problem
There are many problems that you need to solve building a solid and scalable scraping system. Thankfully, many of them can be solved by picking a good library or framework. But there are few, that I could not solve and they were pretty important to me:

- Every time you start, you have to write some amount of boilerplate code which might be annoying. Even with a good scraping library or framework - you have to write some glue code.
- Web pages tend to evolve and change their markup. It leads to frequent code updates and re-deploys.
- Nowadays, more and more websites use dynamic page rendering and old plain HTTP GET request is not a solution any more. It gets worse if the pages require some user interaction to get the data you need.

Of course, these problems are solvable, but wouldn't it be nice if we didn't have to solve them again and again every time we write scraping code?
So, I started looking a for a tool that would solve the issues for me, but unfortunately, I could not find any. Instead I created one.

## The solution
<figure class="image is-128x128" style="margin: 0 auto;">
    <img src="https://cdn-images-1.medium.com/max/1600/0*06mYhTdfqmr4QI_K" />
</figure>

<p class="subtitle is-7 has-text-centered has-text-grey" style="margin-bottom: 2em;">
    Logo by Chris Friel
</p>

And this is how Ferret was born, a modern web scraping tool with a declarative query language.    

Ferret is a DSL with a runtime that represents pretty powerful and extensible declarative query language.

It provides very high level solution that lets users focus on data they are looking for instead of writing all that glue code.

It's embeddable. Even though, it comes with a CLI, it was designed to be easily embedded into applications from the beginning.

And of course, it's super extensible. You can register not only your own functions, but also create custom types that Ferret runtime will recognize.

The overall idea is simple, you write a query using the declarative query language focusing only on data you need, and Ferret handles the rest of the work for you.

<div style="width: 700px;margin: 0 auto;">
    <iframe style="width: 700px; height:525px" src="https://www.youtube.com/embed/g9gWv1PODi8" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>

### Show me the code
Let's see how it works.

For example, to scrape a list of trending repositories on GitHub, you can use the following query:

{{< code fql >}}
LET page = DOCUMENT("https://github.com/trending")

FOR row IN ELEMENTS(page, "ol.repo-list li")
    LET name = INNER_TEXT(row, "div:nth-child(1)")
    LET description = INNER_TEXT(row, "div:nth-child(3)")

    RETURN { name, description }
{{</ code >}}

There are not many things happening above, but they solve 2 of our 3 problems: boilerplate code and frequent markup changes. We will discuss how to solve the 3rd later.

Let's take a closer look.

The first problem gets solved automatically by the language itself - Ferret abstracts away all infrastructural glue code from a user. All heavy lifting is done under the hood. All a user needs to do is to describe what data they need to scrape, and what shape of the data should be as an output. Declarative languages allow you to express yourself in a very readable way.

The second problem gets solved by the fact that FQL is an interpreted language. Script updates do not require you to redeploy your host server. You can easily change it on your production server without any harm.

Now, let's discuss what is happening in our script in order to get better understanding of the tool.

First of all, we tell Ferret what page to load by calling ``DOCUMENT`` function. Ferret loads the page and creates its representation in memory as an object. Then we assigned the returned document object to a variable. Yes, in FQL it's possible to define variables which is an extremely useful feature.

Once the document is loaded, we use the ``ELEMENTS`` function to find all list entries and iterate over them using FQL keywords ``FOR IN``, the construction which is familiar to many developers. Notice, for the search we use regular CSS selector like we would use in a browser.

Inside the ``FOR`` loop body we get inner text of elements of a row using the helper function ``INNER_TEXT``. And again we use variables and CSS selector for simplicity and more efficiency.

After getting inner text of all elements we are interested in, we use ``RETURN`` statement with an object literal. Those who are familiar with JavaScript will recognize the syntax - an object literal with property shorthand. In other words, on each iteration we create an object with 2 properties: name and description. The ``RETURN`` statement works as aggregator, and the result of the query is an array of objects.

### Dynamic pages
The example above solves only 2 problems, and we still cannot scrape pages that are rendered dynamically.

What if we need to get a list of top songs from SoundCloud? If we execute the following query, we will get an empty array:

{{< code fql >}}
LET doc = DOCUMENT("https://soundcloud.com/charts/top") 
LET tracks = ELEMENTS(doc, ".chartTrack__details")

FOR track IN tracks
    RETURN {
        artist: TRIM(INNER_TEXT(track, ".chartTrack__username")),
        track: TRIM(INNER_TEXT(track, ".chartTrack__title"))
    }
{{</ code >}}

In order to fix the query we need to do 2 things:

- Run Chrome or Chromium either in [Docker](https://hub.docker.com/r/alpeware/chrome-headless-stable) or with *--remote-debugging-port=9222* argument.
- Add a boolean argument to the ``DOCUMENT`` function — this informs Ferret that the page document is dynamic, which causes to use different driver for that.

Here is an updated version of the query:

{{< code fql >}}
LET doc = DOCUMENT("https://soundcloud.com/charts/top") 
LET doc = DOCUMENT('https://soundcloud.com/charts/top', true)

WAIT_ELEMENT(doc, '.chartTrack__details', 5000)

LET tracks = ELEMENTS(doc, '.chartTrack__details')

FOR track IN tracks
    RETURN {
        artist: TRIM(INNER_TEXT(track, '.chartTrack__username')),
        track: TRIM(INNER_TEXT(track, '.chartTrack__title'))
    }
{{</ code >}}

You may notice, that there is a new line of code: ``WAIT_ELEMENT(doc, ‘.chartTrack__details’, 5000)``. Since everything is dynamically rendered, we need to be sure that the data we need is in place. That’s why we block the execution until the table with tracks appears on the page or the operation times out. The signature is pretty similar to what we have seen above with an extra timeout argument.

Once the data is loaded and rendered, we are ready to iterate over all table rows and parse them.

As you can see, the way how we had written the query didn’t change much, but at the same time we need to take into consideration some nuances of dynamic pages and give Ferret some hints how to handle them.

### Page interactions
Ok, the above examples may not look very impressive and cannot be a real reason to build such a relatively complex solution like Ferret.

In order to turn the Internet into the Database of the Web we need to be able to interact with the Web like a real user. Clicks, data inputs or just mouse events — all these regular actions need to be supported.

And Ferret supports them.

Let’s imagine, that we are building a website that aggregates product prices from all popular sales platforms like Amazon and eBay.

In order to get a list of products and their information we need to implement the following algorithm:

- Type a search criteria into an input box
- Click “Search” button
- Wait when the page gets loaded and the data gets rendered
- Iterate over the results
- Move to a next page
- Repeat step #3

{{< code fql >}}
LET amazon = DOCUMENT('https://www.amazon.com/', true)

INPUT(amazon, '#twotabsearchtextbox', @criteria)
CLICK(amazon, '.nav-search-submit input[type="submit"]')
WAIT_NAVIGATION(amazon)

LET resultListSelector = '#s-results-list-atf'
LET resultItemSelector = '.s-result-item.celwidget'
LET nextBtnSelector = '#pagnNextLink'
LET vendorSelector1 = 'div > div:nth-child(3) > div:nth-child(2) > span:nth-child(2)'
LET vendorSelector2 = 'div > div:nth-child(5) > div:nth-child(2) > span:nth-child(2)'
LET priceWholeSelector = 'span.sx-price-whole'
LET priceFracSelector = 'sup.sx-price-fractional'
LET pages = TO_INT(INNER_TEXT(amazon, '#pagn > span.pagnDisabled'))

LET result = (
    FOR pageNum IN 1..pages
        LET clicked = pageNum == 1 ? false : CLICK(amazon, nextBtnSelector)
        LET wait = clicked ? WAIT_NAVIGATION(amazon) : false
        LET waitSelector = wait ? WAIT_ELEMENT(amazon, resultListSelector) : false

        LET items = (
            FOR el IN ELEMENTS(amazon, resultItemSelector)
                LET priceWholeTxt = INNER_TEXT(el, priceWholeSelector)
                LET priceFracTxt = INNER_TEXT(el, priceFracSelector)
		        LET price = TO_FLOAT(priceWholeTxt + "." + priceFracTxt)
                LET vendor = ELEMENT_EXISTS(el, vendorSelector1) ? INNER_TEXT(el, vendorSelector1) : INNER_TEXT(el, vendorSelector2)

                RETURN {
                    title: INNER_TEXT(el, 'h2'),
                    vendor,
                    price
                }
        )

        RETURN items
)

RETURN FLATTEN(result)
{{</ code >}}

At this point, most of the constructions of the language should be familiar, so I will just highlight what has not been covered before.

First of all, it’s data input and submitting the form by clicking the search button. Notice, we use parameters (words starting with @ symbol) in this script which. Their values are define either via CLI or a context in embedded mode (I will cover how to embed Ferret to you application in next blog post).

Once the results page gets loaded, we iterate over a range of page numbers, using another parameter. Ideally, you need to extrapolate amount of pages and use this number, but I use the parameter here for simplicity.

On each iteration, we navigate to the target page (except 1st page) by clicking the pagination button and again wait till it gets loaded. After that we just parse the result table.

As you can see, most of the actions are similar to what we do as regular users. Also, it’s a good demonstration of how web crawlers can be easily implemented using FQL.

Here you can find more examples and use cases. PR are welcome.

## Who needs it?
Very reasonable and popular question.

Well, first of all I need it.

As for others, I can only suggest that this tool **might be** useful for data scientists, ML/AI researchers, QA specialists and anyone who needs to scrape a website without setting up a new Python/Node/etc project.

## What next?
The project is still growing and one of the main milestones is to stabilize API and provide rich functionality out of the box. Besides the Ferret itself, there are some other projects that are planned to be built:

<h4>
    <a href="https://github.com/MontFerret/ferret-server">
    Ferret Server
    </a>
</h4>
As you may noticed, Ferret is just a library that comes with a handy CLI. But it does not yet provide a production grade solution for complex scraping. That’s what Ferret Server aims to solve — provide a web scraping platform. The platform that would allows you to persists scripts and scraped data in a database, executes scripts manually or with a scheduler, notify 3rd party systems and etc. The project has started but it’s in very early stage.

<h4>
    <a href="https://github.com/MontFerret/ferret/issues?utf8=%E2%9C%93&q=is%3Aopen+label%3Aarea%2Fdrivers+Create+in%3Atitle+">
    More HTML drivers
    </a>
</h4>
Having more HTML drivers will let QA and developers easily test their web apps in different web browsers using same query.

<h4>
    <a href="https://github.com/MontFerret/ferret-registry">
    Ferret Registry
    </a>
</h4>
See it as NPM registry for FQL scripts where people can share and explore scripts which will make data scraping even easier!

<h4>
    Jupyter kernel
</h4>
I’m not really well familiar with Jupyter, but I know that many data scientists use it, so I figured that it could be useful to integrate Ferret with it.

So, if you have some spare time and want to contribute to an open source project — you are very welcome!