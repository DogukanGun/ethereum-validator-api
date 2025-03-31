#!/bin/bash

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "ImageMagick is required but not installed. Please install it first."
    exit 1
fi

# Generate PNG versions
convert eth-icon.svg -resize 16x16 favicon-16x16.png
convert eth-icon.svg -resize 32x32 favicon-32x32.png
convert eth-icon.svg -resize 180x180 apple-touch-icon.png
convert eth-icon.svg -resize 192x192 android-chrome-192x192.png
convert eth-icon.svg -resize 512x512 android-chrome-512x512.png

# Generate ICO file (contains both 16x16 and 32x32)
convert favicon-16x16.png favicon-32x32.png favicon.ico

# Generate OG image with text
convert eth-icon.svg -resize 1200x630 \
    -gravity center \
    -background "#1E293B" \
    -extent 1200x630 \
    og-image.png

# Copy SVG for Safari pinned tab
cp eth-icon.svg safari-pinned-tab.svg

echo "All favicon files have been generated!" 