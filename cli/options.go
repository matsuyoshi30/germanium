package cli

type Options struct {
	Output            string `short:"o" long:"output" default:"output.png" description:"Write output image to specific filepath"`
	BackgroundColor   string `short:"b" long:"background" default:"#aaaaff" description:"Background color of the image"`
	Font              string `short:"f" long:"font" default:"Hack-Regular" description:"Specify font eg. 'Hack-Bold'"`
	Language          string `short:"l" long:"language" description:"The language for syntax highlighting"`
	Style             string `short:"s" long:"style" description:"The style for syntax highlighting"`
	Clipboard         bool   `short:"c" long:"clip" description:"Copy image to clipboard"`
	ListStyles        bool   `long:"list-styles" description:"List all available styles for syntax highlighting"`
	ListFonts         bool   `long:"list-fonts" description:"List all available fonts in your system"`
	NoLineNum         bool   `long:"no-line-number" description:"Hide the line number"`
	NoWindowAccessBar bool   `long:"no-window-access-bar" description:"Hide the window access bar"`
	ShowVersion       bool   `short:"v" long:"version" description:"Show version"`
	FontSize          string `long:"font-size" default:"24" description:"Specify size of font"`
	Square            bool   `long:"square" description:"Image padded to 1:1 aspect ratio"`
    Padding           string `long:"padding" default:"60" description:"Padding around the code"`
}
