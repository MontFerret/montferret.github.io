# Tool documentation conventions

Documentation for first-party Ferret tools lives under `content/docs/tools`.
Use a shared vocabulary where tools expose the same concept, but keep
tool-specific pages and workflows distinct.

## Structure

- Use topic filenames such as `installation.md`, `configuration.md`, and
  `deployment.md` instead of shortened variants.
- Use `_index.md` as the concise tool landing page. It should explain what the
  tool is, what it is used for, and where the reader should go next.
- Add `overview.md` only when it contains substantial conceptual material that
  does not belong on the landing page.
- Do not add empty or redundant pages only to make tool directories match.

## Installation pages

When applicable, order sections as `Prebuilt binary`, `Install script` or
package manager, `Docker`, `Build from source`, `Verify installation`,
`Update`, `Uninstall`, and `Next steps`. Omit methods and lifecycle sections
that the tool does not support.

Reuse the existing terminal, code, related-page, and other Hugo shortcodes.
Use `data/versions.yaml` for version-sensitive release links, downloads,
images, and build requirements rather than copying version strings into pages.
