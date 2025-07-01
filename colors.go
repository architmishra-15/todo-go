package main

import (
	"strings"
)

// Color holds ANSI escape codes for terminal text styling and coloring.
type Color struct {
	// Colors
	Black  string
	Red    string
	Green  string
	Yellow string
	Blue   string
	Purple string
	Cyan   string
	White  string

	// Styles
	Bold      string
	Italic    string
	Underline string
	Strike    string

	// Reset code
	Reset string
}

// Colors is the global instance you can use anywhere in your project.
var Colors = Color{
	Black:  "\033[30m",
	Red:    "\033[31m",
	Green:  "\033[32m",
	Yellow: "\033[33m",
	Blue:   "\033[34m",
	Purple: "\033[35m",
	Cyan:   "\033[36m",
	White:  "\033[37m",

	Bold:      "\033[1m",
	Italic:    "\033[3m",
	Underline: "\033[4m",
	Strike:    "\033[9m",

	Reset: "\033[0m",
}

// Format wraps the given text with all provided ANSI codes (colors or styles) and appends a reset.
// Usage: Colors.Format("text", Colors.Bold, Colors.Underline, Colors.Red)
func (c Color) Format(text string, codes ...string) string {
	prefix := strings.Join(codes, "")
	return prefix + text + c.Reset
}

// Convenience methods for common combinations.
func (c Color) BoldText(text string) string      { return c.Format(text, c.Bold) }
func (c Color) ItalicText(text string) string    { return c.Format(text, c.Italic) }
func (c Color) UnderlineText(text string) string { return c.Format(text, c.Underline) }
func (c Color) StrikeText(text string) string    { return c.Format(text, c.Strike) }
func (c Color) RedText(text string) string       { return c.Format(text, c.Red) }
func (c Color) GreenText(text string) string     { return c.Format(text, c.Green) }
func (c Color) BlueText(text string) string      { return c.Format(text, c.Blue) }

// Example usage
// func main() {
// 	fmt.Println(Colors.RedText("This is red text"))
// 	fmt.Println(Colors.GreenText("This is green text"))
// 	fmt.Println(Colors.Format("Bold and underlined Purple Text", Colors.Bold, Colors.Underline, Colors.Purple))
// }
