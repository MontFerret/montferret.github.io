---
title: "Introduction"
weight: 1
draft: false
---

# What is Ferret?
Ferret project is an ambitious initiative trying to bring the universal platform for writing scrapers without any hassle. It aims to simplify data extraction from the web for UI testing, machine learning, analytics and more.    
      
Ferret allows you to focus on the data by abstracting away the technical details and complexity of underlying technologies using its own declarative language. 

Ferret is extremely portable, extensible and fast.

<hr />

# Motivation
Nowadays data is everything and who owns data - owns the world.

[I](https://github.com/ziflex) have worked on multiple data-driven projects where data was an essential part of a system and I realized how cumbersome writing tons of scrapers is.   
After some time looking for a tool that would let me to not write a code, but just express what data I need, decided to come up with my own solution.

<hr />

# Inspiration
FQL (Ferret Query Language) is heavily inspired by AQL (ArangoDB Query Language).
But due to the domain specifics, there are some differences in how things work.

<hr />

# How it works

<img src="/img/design.png"  />

Ferret consists of the following main parts:
- FQL parser, compiler and runtime
- Standard library
- User functions registry
- HTML drivers and types (in-memory and Chrome Devtools Protocol based)
- CLI

In order to use CDP (Chrome Devtools Protocol based) driver, you should have running Chrome/Chromium with open debugging port (9222), so that Ferret could connect to it.

For in-memory page parsing and processing, no external dependencies are needed.