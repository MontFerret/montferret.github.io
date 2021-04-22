---
title: "Say 'Hello' to new CLI"
subtitle: "Brand new Ferret CLI"
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2021-04-22"
---

Hello fellas,

We are happy to announce that we finally have a brand [new Ferret CLI](https://github.com/MontFerret/cli)!

{{< image src="images/cli-logo.svg" size="256x256" transform="false" >}}

As the project grows, it became clear that we need to separate Ferret Runtime and CLI, to let the projects evolve on their own.

The separation allowed us to review what we need from CLI and rebuild it from the ground up.

Thus, let's dive in and see what's in there.

# Get started

First of all, we need to install it:

```bash
$ curl https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | sudo sh
```

If you want to install it to a particular folder set it to ``LOCATION`` environment variable:

```bash
$ curl https://raw.githubusercontent.com/MontFerret/cli/master/install.sh | LOCATION=my-dir sh
```

# What's broken

Let's start with some breaking changes.

Fortunately, there is only one breaking change: now if you want to execute a script or enter to REPL, you need to call ``exec`` command:

```bash
$ ferret exec
Welcome to Ferret REPL 0.14.1

Please use `exit` or `Ctrl-D` to exit this program.
```

```bash
$ ferret exec my-script.fql
```

# What's new

## Auto-browser launch
With new CLI, you can easily run your script with automatically open local browser:

```bash
$ ferret exec --browser-open my-script.fql
```

``--browser-open`` will inform CLI to find locally installed Chrome or Chromium and open it before execution.

In case if you don't want your browser to be visible during execution, you can pass ``--browser-headless`` flag to open it in the headless mode:

```bash
$ ferret exec --browser-headless my-script.fql
```

<div class="notification is-warning">
  At this moment, CLI does not support browser installation, thus Chrome/Chromium must be installed on your machine before using these flags.
</div>

## Browser management

Also, you can briefly manage your local browser by openning and closing it.

```bash
$ ferret brower open
```

```bash
$ ferret brower open -d
89502
```

``-d`` flag indicates that the browser needs to be open in the background i.e. the proccess will exit once the browser is open and ready to be used and return the process id.

```bash
$ ferret brower close process-id
```

```bash
$ ferret brower close
```

``close`` without a given process id will **try** to find early opened browser with local debugging options.

<div class="notification is-info">
  In the future, we will add more options for browser management.
</div>

## Configuration

New CLI allows you to store overriden default flag values, by using configuration files.

The list of supported configuration keys you can get by using the following command:


```bash
$ ferret config ls

log-level: info
runtime: <nil>
browser-cookies: <nil>
browser-address: <nil>
browser-open: <nil>
browser-headless: <nil>
proxy: <nil>
user-agent: <nil>
```

For example, let's set ``browser-open`` to ``true`` so, every time we run our script, the browser would be opened:

```bash
$ ferret config set browser-open true
```

And then, when you launch the REPL, you will see your local browser will be automatically opened:

```bash
$ ferret exec
```

Once we exit the REPL, browser will be closed.

Also, you can get the current value of a particular flag:

```bash
$ ferret config get browser-open
true
```

# Summary

That's it, folks! 
Having a separate project for CLI helps us focus more on CLI-only features outside of the Ferret Runtime and improve its usability.

If you have any feature requests welcome to [GitHub Issues](https://github.com/MontFerret/cli/issues) or to discuss something join us on [GitHub Discuss board](https://github.com/MontFerret/cli/discussions)!