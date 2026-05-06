package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envsync/internal/differ"
	"github.com/user/envsync/internal/parser"
	"github.com/user/envsync/internal/syncer"
)

func main() {
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	diffMask := diffCmd.Bool("mask", false, "mask secret values in output")

	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	syncOverwrite := syncCmd.Bool("overwrite", false, "overwrite changed keys in target")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsync <diff|sync> [options] <source> <target>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "diff":
		diffCmd.Parse(os.Args[2:])
		args := diffCmd.Args()
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: envsync diff [--mask] <source> <target>")
			os.Exit(1)
		}
		runDiff(args[0], args[1], *diffMask)

	case "sync":
		syncCmd.Parse(os.Args[2:])
		args := syncCmd.Args()
		if len(args) != 2 {
			fmt.Fprintln(os.Stderr, "usage: envsync sync [--overwrite] <source> <target>")
			os.Exit(1)
		}
		runSync(args[0], args[1], *syncOverwrite)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runDiff(src, dst string, maskSecrets bool) {
	source, err := parser.Parse(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading source: %v\n", err)
		os.Exit(1)
	}
	target, err := parser.Parse(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target: %v\n", err)
		os.Exit(1)
	}

	entries := differ.Diff(source, target)
	fmt.Print(differ.Format(entries, maskSecrets))
}

func runSync(src, dst string, overwrite bool) {
	mode := syncer.ModeAddMissing
	if overwrite {
		mode = syncer.ModeOverwrite
	}

	res, err := syncer.Sync(src, dst, mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sync error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("added: %d, updated: %d, skipped: %d\n",
		len(res.Added), len(res.Updated), len(res.Skipped))
	for _, k := range res.Added {
		fmt.Printf("  + %s\n", k)
	}
	for _, k := range res.Updated {
		fmt.Printf("  ~ %s\n", k)
	}
	for _, k := range res.Skipped {
		fmt.Printf("  . %s (skipped)\n", k)
	}
}
