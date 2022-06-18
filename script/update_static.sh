#!/bin/bash

# update minireset.css
# https://jgthms.com/minireset.css/
curl -L -o src/static/static/css/minireset.min.css  \
  https://raw.githubusercontent.com/jgthms/minireset.css/master/minireset.min.css

# update Alpine.js
# https://alpinejs.dev/
curl -L -o src/static/static/js/alpine.min.js  \
  https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js

# update htmx
# https://htmx.org/
curl -L -o src/static/static/js/htmx.min.js  \
  https://unpkg.com/htmx.org@1.x.x/dist/htmx.min.js
