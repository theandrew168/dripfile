#!/bin/bash

# update tachyons.css
# https://tachyons.io/
curl -L -o internal/web/static/css/tachyons.min.css  \
  https://unpkg.com/tachyons@4.12.0/css/tachyons.min.css

# update Alpine.js
# https://alpinejs.dev/
curl -L -o internal/web/static/js/alpine.min.js  \
  https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js

# update htmx
# https://htmx.org/
curl -L -o internal/web/static/js/htmx.min.js  \
  https://unpkg.com/htmx.org@1.x.x/dist/htmx.min.js
