package machine

import (
	"koding/klientctl/commands/cli"
	"koding/klientctl/endpoint/machine"
	"koding/klientctl/endpoint/team"

	"github.com/spf13/cobra"
)

type listOptions struct {
	jsonOutput bool
}

// NewListCommand creates a command that displays remote machines which belong
// to the user or that can be accessed by their.
func NewListCommand(c *cli.CLI, aliasPath ...string) *cobra.Command {
	opts := &listOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available machines",
		Long:    `This command gets currently available machines and displays them to the caller.`,
		RunE:    listCommand(c, opts),
	}

	// Flags.
	flags := cmd.Flags()
	flags.BoolVar(&opts.jsonOutput, "json", false, "output in JSON format")

	// Middlewares.
	cli.MultiCobraCmdMiddleware(
		cli.WithMetrics(aliasPath...), // Gather statistics for this command.
		cli.NoArgs,                    // No custom arguments are accepted.
	)(c, cmd)

	return cmd
}

func listCommand(c *cli.CLI, opts *listOptions) cli.CobraFuncE {
	return func(cmd *cobra.Command, args []string) error {
		infos, err := machine.List(&machine.ListOptions{})
		if err != nil {
			return err
		}

		if t := team.Used(); t.Valid() == nil {
			all := infos
			infos = infos[:0]

			for _, i := range all {
				if i.Team == t.Name {
					infos = append(infos, i)
				}
			}
		}

		if opts.jsonOutput {
			cli.PrintJSON(c.Out(), infos)
			return nil
		}

		tabListFormatter(c.Out(), infos)
		return nil

	}
}
