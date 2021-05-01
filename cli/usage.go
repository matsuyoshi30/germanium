package cli

const Usage = `USAGE:
    %s [FLAGS] [FILE]

FLAGS:
    -o, --output <PATH>       Write output image to specific filepath [default: ./output.png]
    -b, --background <COLOR>  Background color of the image [default: #aaaaff]
    -f, --font <FONT>         Specify font eg. 'Hack-Bold'
    -l, --language <LANG>     The language for syntax highlighting eg. 'go'
    --list-fonts              List all available fonts in your system
    --no-line-number          Hide the line number
    --no-window-access-bar    Hide the window access bar

AUTHOR:
    matsuyoshi30 <sfbgwm30@gmail.com>
`
