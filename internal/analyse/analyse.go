/*
Package analyse provides sub-commands to analyze GIF file.
*/
package analyse

import "github.com/koron-go/subcmd"

var Analyse = subcmd.DefineSet("analyse", "analyse tools",
	Entropy,
)
