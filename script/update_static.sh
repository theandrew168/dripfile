#!/bin/bash

# update Alpine.js
# https://alpinejs.dev/
curl -L -o static/js/alpine.min.js  \
  https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js

# update htmx
# https://htmx.org/
curl -L -o static/js/htmx.min.js  \
  https://unpkg.com/htmx.org@1.x.x/dist/htmx.min.js
