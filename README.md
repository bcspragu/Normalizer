# Normalizer
Finds the smallest image in a directory, and makes all the images have those dimensions. Does things concurrently.

Makes a number of (bad) assumptions:
- All the images are the same ratio (or it doesn't matter if they get squished)
- The smallest image is < 10,000 pixels across in either direction
