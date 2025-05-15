package main

import (
	"fmt"
	"os"

	"bitmap/config"
	"bitmap/internal/core"
	"bitmap/internal/crop"
	"bitmap/internal/filter"
	"bitmap/internal/header"
	"bitmap/internal/mirror"
	"bitmap/internal/rotate"
)

var applyFeatures = map[string]func(*core.BitMap){
	"filter": filter.HandleFilter,
	"rotate": rotate.HandleRotate,
	"mirror": mirror.HandleMirror,
	"crop":   crop.HandleCrop,
}

func main() {
	config.InitFlags()

	file, err := os.Open(config.SourceFileName)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	b := core.NewBitMap()
	b.Read(file)

	if config.HeaderCmd != nil {
		header.PrintHeaderInfo(b)
	}

	if config.ApplyCmd != nil {
		for _, feature := range config.OrderedFlags {
			applyFeatures[feature](b)
		}
		file, err = os.Create(config.OutputFileName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		b.Save(file)
	}
}
