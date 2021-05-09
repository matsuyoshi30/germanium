#!/usr/bin/env bash
# Test a style that has a light background color, and missing colors for some tokens
# The line numbers should be black
../../../germanium -s autumn main.go -o light-style-gen.png
