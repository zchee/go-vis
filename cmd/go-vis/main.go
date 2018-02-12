// Copyright 2018 The go-vis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command go-vis is the visualizes the package methods dependencies.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/zchee/go-vis/internal/logger"
	"github.com/zchee/go-vis/parser"
	"go.uber.org/zap"
)

var (
	fpath = flag.String("path", "", "comma separated path to visualize package")
)

func main() {
	flag.Parse()
	if *fpath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zaplogger, undo := logger.NewRedirectZapLogger()
	defer undo()
	ctx = logger.NewContext(ctx, zaplogger)

	log := logger.FromContext(ctx).Named("main")

	paths := strings.Split(*fpath, ",")
	paths[len(paths)-1] = strings.TrimSuffix(paths[len(paths)-1], ",")
	pkgTypes, err := parser.InspectDir(paths...)
	if err != nil {
		log.Fatal("parser.ParseFile", zap.Error(err))
	}

	builder := new(strings.Builder)
	parser.Render(builder, pkgTypes)

	fmt.Printf("%s", builder.String())
}
