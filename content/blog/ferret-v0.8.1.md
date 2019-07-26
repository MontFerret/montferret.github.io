---
title: "Ferret v0.8.1"
subtitle: "Hot fixes"
draft: false
author: "Tim Voronov"
authorLink: "https://www.twitter.com/ziflex"
date: "2019-07-26"
---

With all big features, come small bugs. :)

With this new **[v0.8.1](https://github.com/MontFerret/ferret/releases/tag/v0.8.1)** release we are providing hot fixes for some of these small ones.

## Broken CLICK function
In 0.7 ``CLICK`` and ``CLICK_ALL`` functions were doing a check before clicking and was returning a boolean value indicating whether a target element was clicked.
In 0.8 this logic got broken due to internal refactorings and the check was not performed which led to broken scripts using this returned value.

## Scroll jumps
In 0.8, we introduct new interaction behavior with forms. Everytime you want to ``INPUT`` or ``CLICK`` Ferret would scroll into a target element before the operation. Unfortunately, we didn't add any checks *whether* the element is in a viewport, which caused weird page jumps.


Not that bad so far! If you find any issues, please **[report](https://github.com/MontFerret/ferret/issues)**.