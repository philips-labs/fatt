package options

import "github.com/spf13/cobra"

// CommandFlagger allows to add flags to commands
type CommandFlagger interface {
	AddFlags(*cobra.Command)
}
