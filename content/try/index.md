---
title: "Try it!"
slug: "/try/"
type: "repl"
draft: false
---

{{< editor id="replEditor" sharable="true" >}}
LET doc = DOCUMENT('https://soundcloud.com/charts/top', {
    driver: 'cdp'
})

WAIT_ELEMENT(doc, '.chartTrack__details', 5000)

LET tracks = ELEMENTS(doc, '.chartTrack__details')

FOR track IN tracks
    RETURN {
        artist: TRIM(INNER_TEXT(track, '.chartTrack__username')),
        track: TRIM(INNER_TEXT(track, '.chartTrack__title'))
    }
{{</ editor >}}