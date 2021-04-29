package cli

type Options struct {
	Output            string `short:"o" long:"output" default:"output.png" description:"Write output image to specific filepath"`
	BackgroundColor   string `short:"b" long:"background" default:"#aaaaff" description:"Background color of the image"`
	Font              string `short:"f" long:"font" default:"Hack-Regular" description:"Specify font eg. 'Hack-Bold'"`
	Language          string `short:"l" long:"language" description:"The language for syntax highlighting"`
	ListFonts         bool   `long:"list-fonts" description:"List all available fonts in your system"`
	NoLineNum         bool   `long:"no-line-number" description:"Hide the line number"`
	NoWindowAccessBar bool   `long:"no-window-access-bar" description:"Hide the window access bar"`
}

var opts Options
