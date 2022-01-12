.POSIX:
.SUFFIXES:

images_png =            \
  img/logo_branded.png  \
  img/logo_black.png    \
  img/logo_white.png
images_webp = $(images_png:.png=.webp)

.PHONY: default
default: images

.PHONY: images
images: $(images_webp)

.SUFFIXES: .png .webp
.png.webp:
	cwebp -q 50 -o $@ $<
