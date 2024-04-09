package main

import "github.com/jialequ/linux-sdk/tools/goctl/compare/cmd"

// EXPERIMENTAL: compare goctl generated code results between old and new, it will be removed in the feature.
// : BEFORE RUNNING: export DSN=$datasource, the database must be gozero, and there has no limit for tables.
// : AFTER RUNNING: diff --recursive old_fs new_fs

func main() {
	cmd.Execute()
}
