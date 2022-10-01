package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "kitus-protoc",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.Flags().BoolVar(&opts.GoPb, "go-pb", false, "")
	rootCmd.Flags().BoolVar(&opts.KitusPb, "kitus-pb", false, "")
	rootCmd.Flags().StringVar(&opts.generatorOpts.FilePath, "file-path", "", "")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
