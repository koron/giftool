package analyse

import "github.com/koron-go/subcmd"

var Analyse = subcmd.DefineSet("analyse", "analyse tools",
	Entropy,
)
