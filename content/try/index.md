---
title: "Try it!"
slug: "/try/"
type: "repl"
draft: false
---

{{< editor id="replEditor" sharable="true" >}}
// Open the SoundCloud Top Charts page using a browser-based driver (CDP)
// This allows Ferret to execute JavaScript and work with dynamic content
LET doc = DOCUMENT('https://soundcloud.com/charts/top', {
    driver: 'cdp'
})

// Wait until at least one chart tile is present on the page
// This is important because SoundCloud loads content asynchronously
WAIT_ELEMENT(doc, '.audibleTile', 5000)

// Select all track tiles from the page
LET tracks = ELEMENTS(doc, '.audibleTile')

// Iterate over each track tile and extract useful data
FOR track IN tracks
    RETURN {
        // Chart position / description text shown on the tile
        chart: TRIM(INNER_TEXT(track, '.playableTile__descriptionContainer')),
        
        // Build an absolute URL to the track page
        // The link on the page is relative, so we prepend the SoundCloud domain
        link: "https://soundcloud.com" +
                  TRIM(ELEMENT(track, '.playableTile__artworkLink')?.attributes.href)
        }
{{</ editor >}}