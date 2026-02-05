---
title: "Say 'Hello' to Lab"
subtitle: "A new way to test Ferret scripts"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2020-11-09"
---

{{< image src="images/lab-logo.png" size="256x256" transform="false" >}}

Hello friends,

Our team is excited to announce a new tool in our ecosystem - [The Lab](https://github.com/MontFerret/lab)!

The Lab is a test runner for Ferret that helps you to test how your scripts retrieve data or validate web UI.

Let's see how it works!

## Get started

The Lab can be installed either by downloading a [binary](https://github.com/MontFerret/lab/releases) or pulling a Docker image:

```bash
$ docker pull montferret/lab
```

If you are macOS/Linux user, you can install the latest binary with the shell script:

```bash
$ curl https://raw.githubusercontent.com/MontFerret/lab/master/install.sh | sudo sh
```

In this tutorial, we are going to use a local binary.

## First run

Let's create a simple script that retreives GitHub trending projects:

{{< code lang="fql" height="280px" >}}
LET github = DOCUMENT("https://github.com/trending")

FOR project IN ELEMENTS(github, ".Box-row")
    LET link = ELEMENT(project, "h1 a")
    LET href = link.attributes.href

    RETURN {
    	url: "https://github.com" + href,
    	name: SUBSTITUTE(link.innerText, " ", ""),
    	description: ELEMENT_EXISTS(project, "p") ? INNER_TEXT(project, "p") : "",
    	language: ELEMENT_EXISTS(project, '[itemprop="programmingLanguage"]') ? INNER_TEXT(project, '[itemprop="programmingLanguage"]') : "",
    	stars: TRIM(INNER_TEXT(project, 'a[href="'+href+'/stargazers"]'))
    }
{{</ code >}}

Let's save it in a file called ``trending.fql``.
Now we can the query with ``lab``:

```bash
$ lab trending.fql
9:04PM DBG using User-Agent user-agent=
9:04PM INF Passed Duration=1.32741042s File=/Users/ziflex/Work/src/github.com/MontFerret/_scripts/trending.fql Times=1
9:04PM INF Done Duration=1.32741042s Failed=0 Passed=1
```

``lab`` has successfully executed the given script! Even though it might not look very impressive, it is pretty useful when you write scripts to test a UI, like this:

{{< code lang="fql" height="250px" >}}
LET github = DOCUMENT("https://github.com/trending")
LET emailBtnSelector = 'a[href="/explore/email"]'
LET trendingBtnSelector = 'a[data-selected-links="trending_repositories /trending"]'

T::TRUE(ELEMENT_EXISTS(github, emailBtnSelector))
T::EQ(TRIM(INNER_TEXT(github, emailBtnSelector)), 'Get email updates')

T::TRUE(ELEMENT_EXISTS(github, trendingBtnSelector))
T::INCLUDE(ELEMENT(github, trendingBtnSelector).attributes.class, 'selected')

RETURN TRUE
{{</ code >}}

Let's save this query in a file called ``assertions.fql`` and execute it:
```bash
$ lab assertion.fql
9:06PM DBG using User-Agent user-agent=
9:06PM INF Passed Duration=415.113954ms File=/Users/ziflex/Work/src/github.com/MontFerret/_scripts/assertion.fql Times=1
9:06PM INF Done Duration=415.113954ms Failed=0 Passed=1
```
Great! Now we have 2 scripts that we can run together:

```bash
$ lab .
9:07PM DBG using User-Agent user-agent=
9:07PM INF Passed Duration=1.402902567s File=/Users/ziflex/Work/src/github.com/MontFerret/_scripts/assertion.fql Times=1
9:07PM DBG using User-Agent user-agent=
9:07PM INF Passed Duration=134.648972ms File=/Users/ziflex/Work/src/github.com/MontFerret/_scripts/trending.fql Times=1
9:07PM INF Done Duration=1.537551539s Failed=0 Passed=2
```

## Suite files

All that good, but why would we bother writing something that can be achieved by Ferret alone?

The answer is suite files. 

Till this moment, we've used plain FQL files which ``lab`` directly executed using ``ferret`` without possibility to validate correctness of returned data.  

Here come suite files:

```yaml
query:
  ref: trending.fql
assert:
  text: |
    RETURN T::LEN(@lab.data.query.result, 25)
```

Ok, what's happening here is that we have decoupled our query from the assertions. ``lab`` executes a query from the ``query`` section and then validates returned data by executing another query from the ``assert`` section.

The ``@lab`` parameter is a built-in parameter that holds runtime information, like query results and ``cdn`` endpoints, which we will discuss later.   


<div class="notification is-info">
  We are using ``ref`` keyword here to point to a external file, but we could also use ``text`` in there as we do in assertion section.
</div>

## Serving static files
One of the features that ``lab`` provides is static files server called ``cdn``. The ``cdn`` endpoint allows you to serve any static files or web apps during execution of the tests.

Let's create a plain text file and simple query that will download a it using the ``cdn`` endpoint.

```bash
$ mkdir content
$ echo "Hello world" > content/readme.txt
$ touch cdn-example.fql
```

{{< code lang="fql" height="90px" >}}
RETURN IO::NET::HTTP::GET(@lab.cdn.content)
{{</ code >}}

As you have alredy seen before, we are using the ``@lab`` param to get access to available cdn endpoints.
``content`` is a directory name that is used for static files.

```bash
$ lab --cdn=./content cdn-example.fql
⇨ http server started on [::]:57150
9:10PM INF Passed Duration=13.387547ms File=/Users/ziflex/Work/src/github.com/MontFerret/_scripts/cdn-example.fql Times=1
9:10PM INF Done Duration=13.387547ms Failed=0 Passed=1
```

In the logs before script execution, we see that ``lab`` launched a web server on port 57150 in order to serve static files.

<div class="notification is-info">
  You can serve as many folders as you want.
</div>

## Parallel execution
By default, ``lab`` executes queries sequentially, but you can make it run in parallel by setting ``concurrency`` parameter to the required level of parallelism:

```bash
$ lab --concurrency=10 .
```

## Docker
Aside from binary files, ``lab`` is also distributed as Docker images via [GitHub](https://github.com/orgs/MontFerret/packages/container/package/lab) and [Docker Hub](https://hub.docker.com/repository/docker/montferret/lab) that have both ``lab`` binaries and headless ``Chromium``:


```bash
$ docker run --rm -v $PWD:/test montferret/lab:1.2.0
```

## Summary

Let's summarize what we've learned.

``lab`` can be useful to test correctness of your scraping queries by running them and validating their outputs.

Also, ``lab`` can help in testing UI by using the new [assertion library](https://www.montferret.dev/docs/stdlib/testing/).

In some scenarios, it can be used as a batch script runner that can execute your scripts in parallel that will simplify your infrastructure.
