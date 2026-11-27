---
title: "Ferret Comes to JetBrains IDEs"
subtitle: "Native FQL development in IntelliJ IDEA and GoLand."
draft: false
author: "Tim Voronov"
authorLink: "https://github.com/ziflex"
date: "2026-11-26"
---

Hello, friends!

A few months ago, when I introduced the official Visual Studio Code extension, I mentioned that it wouldn't be the last editor integration for Ferret.

Today, I'm happy to share that the first release of the official Ferret plugin for JetBrains IDEs is available!

The plugin brings FQL development directly into IntelliJ IDEA and GoLand, with support for:

- Language intelligence and diagnostics
- Code completion and formatting
- Native Run and Debug configurations
- Breakpoints and stepping
- Variable inspection, expression evaluation, and watches

Just like the Visual Studio Code extension, it comes with a bundled Ferret daemon (`ferretd`), so there's no need to install it separately.

## Beyond local execution

One feature I'm particularly excited about is the ability to run and debug FQL programs against runtimes hosted by other applications through [Wire](https://ferretlang.org/docs/embedding/wire/).

As I've mentioned in previous articles, Ferret is a library first. Applications can embed the engine, register their own functions and modules, expose domain-specific capabilities, and define the environment in which FQL programs execute.

This creates an interesting challenge for development tooling. If your application registers its own functions, running an FQL program against a standalone Ferret runtime in your IDE won't necessarily work. That runtime simply doesn't know anything about your application or the capabilities it provides.

With Wire, however, the IDE can connect to the runtime hosted by your application and use it directly. You don't need to reproduce your application's runtime configuration inside the IDE or maintain a separate development environment for your FQL programs.

Let's see what that looks like.

### Exposing your application's runtime

First, your application needs to expose its Ferret runtime through a Wire server.

For this example, we'll use the existing [Wire example application](https://github.com/MontFerret/montferret.github.io/tree/main/examples/wire), which registers a custom `app::greet` function.

You can download and run it with:

```bash
$ git clone https://github.com/MontFerret/montferret.github.io.git
$ cd montferret.github.io/examples/wire
$ go run ./host
```

The application starts a Wire server listening on `127.0.0.1:50051`.

If you're already embedding Ferret in your own application, take a look at the [Wire documentation](https://ferretlang.org/docs/embedding/wire/) for instructions on exposing your existing runtime.

### Configuring the IDE

The Ferret daemon supports named runtime configurations. This allows you to define multiple execution environments and switch between them without changing your FQL programs.

Create a `runtimes.json` file somewhere in your project:

```json
{
  "defaultRuntime": "local",
  "runtimes": [
    {
      "name": "application",
      "label": "My Application",
      "provider": "wire",
      "endpoint": "tcp://127.0.0.1:50051",
      "insecure": true,
      "connectTimeout": "5s"
    }
  ]
}
```

Here, we're keeping the built-in local runtime as the default while introducing another runtime called `application`, which connects to our application's Wire server.

The `insecure` option enables plaintext communication for this local development example. For deployments beyond local development, see the [Wire security guide](https://ferretlang.org/docs/embedding/wire-security/).

In your JetBrains IDE:

1. Open **Settings → Languages & Frameworks → Ferret**.
2. Set **Runtime configuration file** to the newly created `runtimes.json`.
3. Apply the settings.
4. Expand **Runtime catalog** to see the configured runtimes.

If you change the configuration after it has been loaded, use **Reload runtimes** once your active Run and Debug sessions have finished.

### Running and debugging FQL

Now, let's create a simple FQL program:

{{< code lang="fql" >}}
let greeting = app::greet(@name)

return greeting
{{</ code >}}

The `app::greet` function isn't part of Ferret's standard library. It exists only in our example application.

Create a new Ferret Run/Debug configuration and select `application` in the **Runtime** field. Set the source file to your FQL program and provide the following bind parameters:

```json
{
  "name": "Ada"
}
```

Leave the working directory empty and run the program.

The result should be:

```json
"Hello, Ada!"
```

Now set a breakpoint on the `RETURN` statement and start a debugging session using the same configuration.

You can inspect variables, evaluate expressions, add watches, and step through the program using the native JetBrains debugger, even though the actual FQL execution takes place inside another application.

Your application still owns its runtime, functions, modules, services, and configuration. The IDE simply connects to it through Wire and starts a new execution or debugging session. Stopping that session doesn't shut down your application.

And, of course, you can switch back to the built-in local runtime whenever you need it.

## One foundation, multiple development tools

This functionality isn't exclusive to JetBrains. It's also available in the latest release of the Visual Studio Code extension.

Both integrations build on the same Ferret daemon, Universal API, and Wire infrastructure. The editors provide their native development experience, while Ferret remains responsible for the language and its execution.

It's particularly satisfying to see the architectural ideas behind the [Universal API](https://ferretlang.org/blog/ferret-v2-uapi/) and Wire coming together in something that developers can actually use.

## Give it a try

This is the initial v0.1.0 release, with support for IntelliJ IDEA and GoLand 2026.2 and 2026.3.

You can download the plugin from [GitHub](https://github.com/MontFerret/editorium/releases/tag/jetbrains/v0.1.0) and find installation and usage instructions in the [documentation](https://ferretlang.org/docs/tools/jetbrains/).

Both integrations are developed as part of [Editorium](https://github.com/MontFerret/editorium), the home of Ferret's editor integrations.

As always, feedback, suggestions, and bug reports are very welcome.

Happy ferreting! 🦦
