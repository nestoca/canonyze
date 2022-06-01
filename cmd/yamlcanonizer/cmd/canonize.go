package main

import (
	"github.com/nestoca/yamlcanonizer/cmd/yamlcanonizer/internal"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func NewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "yamlcanonizer",
		Short: "Canonizes multiple yaml docs for easier comparing",
		Long:  "Canonizes a set of yaml documents into a standardized and sorted form that is suitable for comparison against another set",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Canonize()
		},
	}

	return cmd
}
