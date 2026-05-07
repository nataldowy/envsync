package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/differ"
	"github.com/yourorg/envsync/internal/parser"
	"github.com/yourorg/envsync/internal/resolver"
	"github.com/yourorg/envsync/internal/syncer"
)

func main() {
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)

	syncMode := syncCmd.String("mode", "add", "sync mode: add|overwrite")
	resolveVars := diffCmd.Bool("resolve", false, "resolve variable interpolation before diffing")
	allowMissing := diffCmd.Bool("allow-missing", false, "allow missing variables during resolution")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envsync <diff|sync> [options] <source> <target>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "diff":
		_ = diffCmd.Parse(os.Args[2:])
		args := diffCmd.Args()
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "diff requires <source> <target>")
			os.Exit(1)
		}
		runDiff(args[0], args[1], *resolveVars, *allowMissing)
	case "sync":
		_ = syncCmd.Parse(os.Args[2:])
		args := syncCmd.Args()
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "sync requires <source> <target>")
			os.Exit(1)
		}
		runSync(args[0], args[1], *syncMode)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runDiff(src, dst string, resolveInterp, allowMissing bool) {
	srcEnv, err := parser.Parse(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", src, err)
		os.Exit(1)
	}
	dstEnv, err := parser.Parse(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", dst, err)
		os.Exit(1)
	}

	if resolveInterp {
		opts := resolver.DefaultOptions()
		opts.AllowMissing = allowMissing
		if srcEnv, err = resolver.Resolve(srcEnv, opts); err != nil {
			fmt.Fprintf(os.Stderr, "resolve %s: %v\n", src, err)
			os.Exit(1)
		}
		if dstEnv, err = resolver.Resolve(dstEnv, opts); err != nil {
			fmt.Fprintf(os.Stderr, "resolve %s: %v\n", dst, err)
			os.Exit(1)
		}
	}

	entries := differ.Diff(srcEnv, dstEnv)
	fmt.Print(differ.Format(entries, true))
}

func runSync(src, dst, mode string) {
	srcEnv, err := parser.Parse(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse %s: %v\n", src, err)
		os.Exit(1)
	}
	if err := syncer.Sync(srcEnv, dst, mode); err != nil {
		fmt.Fprintf(os.Stderr, "sync: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("synced %s -> %s (mode=%s)\n", src, dst, mode)
}
