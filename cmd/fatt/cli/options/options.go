package options

import "github.com/spf13/cobra"

type CommandFlagger interface {
	AddFlags(*cobra.Command)
}
