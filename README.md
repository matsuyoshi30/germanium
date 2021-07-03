# Germanium

Germanium is an alternative to [Carbon](https://github.com/carbon-app/carbon) and [Silicon](https://github.com/Aloxaf/silicon) implemented in Go.

![screenshot](assets/img/sample.png)

Germanium can work without browser and internet like Silicon.


## Usage

```
germanium

USAGE:
    germanium [FLAGS] [FILE]

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
```

### Example

Generate image From file

```
germanium -o main.png main.go
```

Generate image From Stdin (need to add option `-l` or `--language`)

```
cat main.go | germanium -l go -o main.png -
```

Generate image with another style (you can get a list of them with `--list-styles`)

```
germanium -s solarized-dark -o main.png main.go
```

Generate image without line number

```
germanium --no-line-number -o main.png main.go
```

Generate image without window control bar

```
germanium --no-window-access-bar -o main.png main.go
```

Generate image and copy to clipboard

```
germanium --no-window-access-bar -o main.png main.go -c
```

## Install

### GitHub releases

You can download from [here](https://github.com/matsuyoshi30/germanium/releases).

### Arch Linux

```
yay germanium
```

### Build from source

```
git clone https://github.com/matsuyoshi30/germanium
cd cmd/germanium && go install
```

#### Requirements

- go1.16


## LICENSE

[MIT](./LICENSE)


## Author

[matsuyoshi30](https://twitter.com/matsuyoshi30)
