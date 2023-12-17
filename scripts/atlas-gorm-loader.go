package main

import (
	"fmt"
	"io"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	"ariga.io/atlas-provider-gorm/gormschema"

	"github.com/julianstephens/license-server/pkg/model"
)

// Define the models to generate migrations for.
var models = []any{
	&model.User{},
	&model.Product{},
	&model.ProductFeature{},
	&model.License{},
	&model.APIKey{},
}

func main() {
	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	_, err = io.WriteString(os.Stdout, stmts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to write planned schema: %v\n", err)
	}
}
