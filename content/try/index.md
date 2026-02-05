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

WAIT_ELEMENT(doc, '.audibleTile', 5000)

LET tracks = ELEMENTS(doc, '.audibleTile')

FOR track IN tracks
RETURN {
    chart: TRIM(INNER_TEXT(track, '.playableTile__descriptionContainer')),
    link: "https://soundcloud.com" + TRIM(ELEMENT(track, '.playableTile__artworkLink')?.attributes.href)
}
{{</ editor >}}