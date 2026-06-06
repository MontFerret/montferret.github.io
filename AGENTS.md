# AGENTS.md

## Project overview

This repository contains the Ferret website and documentation.

The site is built with Hugo and includes:

- marketing pages
- documentation pages
- custom Hugo shortcodes
- Sass/CSS for the site and docs UI
- embedded FQL examples
- interactive editor blocks used in documentation

The primary goal of this repository is to explain Ferret clearly to developers. Prefer clarity, correctness, and maintainability over cleverness.

## Working principles

When making changes:

- Keep diffs focused and minimal.
- Preserve the existing visual language unless the task explicitly asks for a redesign.
- Prefer simple Hugo/Markdown/Sass changes over adding JavaScript.
- Do not introduce new dependencies unless clearly necessary.
- Do not rewrite large sections opportunistically.
- When editing docs, preserve the author’s tone: concise, direct, practical, and developer-focused.
- Favor concrete examples over abstract explanations.
- Avoid marketing fluff.

## Documentation style

Documentation should be written for developers who may be new to Ferret and FQL.

Use this style:

- Start sections with a short explanation of what the feature is.
- Show a minimal example early.
- Explain what the example does.
- Add edge cases or deeper details only after the basic idea is clear.
- Prefer short paragraphs.
- Prefer headings that describe the user’s task or concept.
- Avoid long bullet lists unless they improve scanning.
- Avoid vague phrases like “powerful”, “seamless”, “robust”, or “easy” unless supported by specifics.

Good documentation should answer:

1. What is this?
2. When would I use it?
3. What does the syntax look like?
4. What should I be careful about?
5. Where should I go next?

## FQL examples

FQL examples must be valid unless the surrounding text explicitly says they are intentionally invalid.

When adding or editing FQL examples:

- Use RETURN for terminal expressions.
- Use LET for immutable bindings.
- Use VAR only when reassignment is required.
- Do not show failing examples inside live editor blocks unless the page is specifically demonstrating errors.
- Keep examples small and focused.
- Prefer realistic examples over toy examples when the concept benefits from context.
- Do not invent syntax. Check existing grammar/docs before adding new syntax.
- Preserve Ferret terminology consistently.

When documenting runtime-specific behavior, make it clear whether the feature belongs to:

- the FQL language
- the Ferret runtime
- a module
- the CLI
- the browser runtime
- an embedding application

## Hugo conventions

Follow the existing Hugo structure and conventions.

When working with content:

- Preserve front matter fields such as title, sidebarTitle, weight, draft, description, and aliases.
- Keep page weights consistent with surrounding pages.
- Use existing shortcodes instead of inventing new page-local markup.
- Reuse existing tabs, terminal blocks, editor blocks, cards, and tiles where possible.
- Do not hard-code absolute URLs for internal docs links unless existing conventions require it.

When working with shortcodes:

- Keep shortcode APIs small and explicit.
- Preserve backward compatibility with existing content.
- Avoid making content authors repeat information that Hugo can infer.
- Prefer readable shortcode usage over highly generic abstractions.

## Sass/CSS conventions

When editing styles:

- Keep layout changes localized.
- Avoid global selector changes unless explicitly required.
- Check desktop and narrow layouts mentally before changing grid, sticky, overflow, or scrollbar behavior.
- Be careful with docs sidebar, table of contents, editor panels, and result panels.
- Prefer existing variables and mixins.
- Do not introduce unrelated color, spacing, or typography changes.
- Avoid fixing one page by breaking reusable docs layout behavior.

For interactive/editor UI:

- Loading indicators should not cause layout jumps.
- Result panels should handle overflow gracefully.
- Empty states should be useful but not visually noisy.
- Controls should remain understandable when loading, disabled, or empty.

## Accessibility and UX

When changing UI:

- Preserve keyboard accessibility.
- Keep focus states visible.
- Do not rely on color alone to communicate meaning.
- Avoid colors that make normal emphasis look like an error or warning.
- Maintain sufficient contrast.
- Prefer predictable layout over flashy animation.
- Avoid layout shifts during loading.

## Build and validation

This project uses Mage for project management. Prefer Mage targets over calling Hugo, npm, or shell scripts directly.

Common commands:

```bash
mage
mage serve
mage build
mage clean
mage install
mage generate
mage publish
```

Available Mage targets include:

* mage serve — cleans the output directory and starts the local Hugo server.
* mage build — cleans the output directory and runs the Hugo build.
* mage clean — removes the generated public directory.
* mage install — installs dependencies for the Ferret theme.
* mage generate — generates stdlib documentation from stdlib-docs-rep.yaml.
* mage publish — publishes the website using publish.sh.

Only run commands that exist in the repository. Check `Magefile.go` before assuming a command is available.

When validating changes:

* Use mage build for a production build.
* Use mage serve for local visual checks.
* Use mage generate only when stdlib docs or templates are affected.
* Use mage install only when theme dependencies need to be installed or refreshed.
* Do not call hugo, npm, or publish.sh directly unless there is a specific reason.
* Do not run mage publish unless the task explicitly asks to publish.

## Git and change discipline

Before editing:

- Inspect the relevant files first.
- Understand nearby conventions.
- Check for existing shortcodes, partials, layouts, or styles before creating new ones.

When done:

- Summarize what changed.
- Mention which checks were run.
- Mention any checks that could not be run.
- Call out risky areas or assumptions.

Do not:

- Reformat unrelated files.
- Rename pages or routes unless explicitly asked.
- Remove aliases without confirmation.
- Change generated files unless the repository expects them to be committed.
- Make broad content rewrites when the task asks for a targeted edit.

## Content priorities

The website should make Ferret feel:

- precise
- practical
- capable
- developer-friendly
- honest about runtime boundaries

Avoid making Ferret sound like a generic scraping tool or a vague automation platform.

Prefer language that explains the actual value:

- structured extraction
- portable FQL scripts
- runtime-defined capabilities
- browser and non-browser execution
- explicit error and timeout behavior
- embeddability
- testability

## When unsure

If a task is ambiguous:

- Make the smallest reasonable change.
- Preserve existing conventions.
- Leave a clear note about assumptions.
- Do not invent missing product behavior.
- Ask for clarification only when the ambiguity blocks safe progress.