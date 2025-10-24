package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long:                  "To load completions",
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(CompletionCmd)
}

func SlugCompletion(
	_ *cobra.Command,
	_ []string,
	toComplete string,
) ([]cobra.Completion, cobra.ShellCompDirective) {
	files, err := os.ReadDir(BlogPath)
	if err != nil {
		return []cobra.Completion{}, cobra.ShellCompDirectiveNoFileComp
	}

	candidates := make([]cobra.Completion, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), toComplete) {
			candidates = append(candidates, file.Name())
		}
	}

	return candidates, cobra.ShellCompDirectiveNoFileComp
}
