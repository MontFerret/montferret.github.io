---
title: "Ferret v0.8.2"
subtitle: "Hot fixes"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-08-05"
---

Another portion of hot fixes has arrived!

**[v0.8.2](https://github.com/MontFerret/ferret/releases/tag/v0.8.2)** release brings some critical fixes.

## Scrolling position is not centered
In 0.8, we introduced a new behavior for all user interactions. Before any actions like clicking ot typing, we scroll into a target element making it visible in a viewport. But the logic behind it, positioned elements on the top of the viewport, which led to some issues on some websites due to their markup - the elements were still invisible because they were under other elements.

## Unable to set custom logger fields
This is a small fix/improvement for internal logger. There was no way to set a specific id for an execution. It makes it difficult to track script execution in embedded scenarios when you have a request id.    
Now you can add as many fields as you want to the logger, thanks to [Zerologger](https://github.com/rs/zerolog).

## Invalid parent used for INNER_TEXT, INNER_HTML and friends
This is another issue that was triggered by internal refactoring. A new mechanism for retrieving inner HTML and text by selectors used document for all element look ups instead of a given parent element.

## Unable to set custom headers
There was an error occuring whenever you pass values to ``headers`` object.

As always, if you find any new issues, please **[report](https://github.com/MontFerret/ferret/issues)**.