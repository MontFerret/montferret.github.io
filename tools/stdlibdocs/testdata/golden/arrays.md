---
aliases:
  - "/docs/stdlib/arrays/"
description: "Global functions for working with arrays."
draft: false
functionMenu:
  - label: "append"
    anchor: "global-append"
  - label: "flatten"
    anchor: "global-flatten"
sidebarTitle: "Arrays"
stdlibGenerated: true
stdlibKind: "Category"
title: "Arrays"
type: "docs"
weight: 10
---
<h1>Arrays</h1>

<p>Global functions for working with arrays.</p>

<nav class="stdlib-category-index" aria-label="Arrays functions">
  <ul class="stdlib-category-functions">
    <li><a href="#global-append"><code>append</code></a></li>
    <li><a href="#global-flatten"><code>flatten</code></a></li>
  </ul>
</nav>

<div class="stdlib-category-reference">
<section class="stdlib-api-function" aria-labelledby="global-append">
  <h2 id="global-append">
  <a class="stdlib-api-entity-link" href="#global-append" aria-label="append" title="Link to append">
    <code>append</code><span aria-hidden="true">#</span>
  </a>
</h2>
  <section class="stdlib-api-signature" aria-labelledby="global-append-signature-fixed-2">
  <h3 id="global-append-signature-fixed-2">
    <a class="stdlib-api-entity-link" href="#global-append-signature-fixed-2">2-parameter signature<span aria-hidden="true">#</span></a>
  </h3>
  <p class="stdlib-api-signature-code"><code>append(array, value)</code></p>
  <p>append adds a value to an array.</p>
  <dl class="stdlib-api-details">
    <div>
      <dt>Parameters</dt>
      <dd><ul><li><span class="stdlib-api-value-heading"><code>array</code><code class="stdlib-api-value-type">Any[]</code></span><span>Target array.</span></li><li><span class="stdlib-api-value-heading"><code>value</code><code class="stdlib-api-value-type">Any</code></span><span>Value to append.</span></li></ul></dd>
    </div><div><dt>Returns</dt><dd class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>Any[]</code></span><span>Updated array.</span></dd></div>
  </dl>
</section>
</section>
<section class="stdlib-api-function" aria-labelledby="global-flatten">
  <h2 id="global-flatten">
  <a class="stdlib-api-entity-link" href="#global-flatten" aria-label="flatten" title="Link to flatten">
    <code>flatten</code><span aria-hidden="true">#</span>
  </a>
</h2>
  <section class="stdlib-api-signature" aria-labelledby="global-flatten-signature-fixed-1">
  <h3 id="global-flatten-signature-fixed-1">
    <a class="stdlib-api-entity-link" href="#global-flatten-signature-fixed-1">1-parameter signature<span aria-hidden="true">#</span></a>
  </h3>
  <p class="stdlib-api-signature-code"><code>flatten(arr)</code></p>
  <p>flatten turns an array of arrays into a flat array.</p>
  <dl class="stdlib-api-details">
    <div>
      <dt>Parameters</dt>
      <dd><ul><li><span class="stdlib-api-value-heading"><code>arr</code><code class="stdlib-api-value-type">Any[]</code></span><span>Target array.</span></li></ul></dd>
    </div><div><dt>Returns</dt><dd class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>Any[]</code></span><span>Flat array.</span></dd></div>
  </dl>
</section>
  <section class="stdlib-api-signature" aria-labelledby="global-flatten-signature-fixed-2">
  <h3 id="global-flatten-signature-fixed-2">
    <a class="stdlib-api-entity-link" href="#global-flatten-signature-fixed-2">2-parameter signature<span aria-hidden="true">#</span></a>
  </h3>
  <p class="stdlib-api-signature-code"><code>flatten(arr, depth)</code></p>
  <p>flatten turns an array of arrays into a flat array.</p>
  <dl class="stdlib-api-details">
    <div>
      <dt>Parameters</dt>
      <dd><ul><li><span class="stdlib-api-value-heading"><code>arr</code><code class="stdlib-api-value-type">Any[]</code></span><span>Target array.</span></li><li><span class="stdlib-api-value-heading"><code>depth</code><code class="stdlib-api-value-type">Int</code></span><span>Depth level.</span></li></ul></dd>
    </div><div><dt>Returns</dt><dd class="stdlib-api-value"><span class="stdlib-api-value-heading"><code>Any[]</code></span><span>Flat array.</span></dd></div>
  </dl>
</section>
</section>
</div>
