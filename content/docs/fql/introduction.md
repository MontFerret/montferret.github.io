---
title: "Introduction"
weight: 1
draft: false
---

# What is Ferret Query Language?

The Ferret Query Language (FQL) is heavily inspired by the ArangoDB Query Language (AQL) and used as a starting point. But due to the domain specifics, there are some differences in how things work and future lagnuage changes are expected.

Even though, FQL is used to read data from the websites, it's considered as a general purpose query language. That means that all web related functionality is implemented as functions from the FQL Standard Library.

FQL is mainly a declarative language, meaning that a query expresses what result should be achieved but not how it should be achieved. 

Like AQL, FQL is similar to the Structured Query Language (SQL). But unlike AQL, FQL supports only reading, all modifications are available through specific functions only.

The syntax of FQL queries is different to SQL, even if some keywords overlap. Nevertheless, FQL should be easy to understand for anyone with an SQL background.