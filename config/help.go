package config

var helpText = `Usage:
  bitmap <command> [arguments]

The commands are:
  header    prints bitmap file header information
  apply     applies processing to the image and saves it to the file
`

var headerHelpText = `Usage:
  bitmap header <source_file>

Description:
  Prints bitmap file header information
`

var applyHelpText = `Usage:
  bitmap apply [options] <source_file> <output_file>

The options are:
  --help      prints program usage information
  --mirror    mirrors the image
  --filter    applies a filter to the image
  --rotate    rotates the image
  --crop      crops the image
`
