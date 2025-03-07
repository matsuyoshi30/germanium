package cli

const Usage = `USAGE:
    %s [FLAGS] [FILE]

FLAGS:
    -o, --output <PATH>       Write output image to specific filepath [default: ./output.png]
    -b, --background <COLOR>  Background color of the image [default: #aaaaff]
    -f, --font <FONT>         Specify font eg. 'Hack-Bold'
    -l, --language <LANG>     The language for syntax highlighting eg. 'go'
    -s, --style <STYLE>       The style for syntax highlighting eg. 'dracula'
    -c, --clip                Copy image to clipboard
    --list-styles             List all available styles for syntax highlighting
    --list-fonts              List all available fonts in your system
    --no-line-number          Hide the line number
    --no-window-access-bar    Hide the window access bar
    -v, --version             Show Version
    --square                  Adds padding to reach 1:1 aspect ratio 
    --padding                 Padding in px

AUTHOR:
    matsuyoshi30 <sfbgwm30@gmail.com>
`
