package cmd

import (
	"fmt"
	"strings"

	"github.com/safurrier/skills-cli/internal/skill"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <skill-name> [--render] [--key value ...] [-- freeform...]",
	Short: "Invoke a skill (default: output pointer prompt)",
	Long: `Invoke a skill. By default outputs a pointer prompt that tells the agent
where to find the skill. Use --render to output the full rendered content.

Dynamic key-value pairs are passed as --key value after the skill name.
Freeform text can be passed after --.`,
	DisableFlagParsing: true,
	SilenceUsage:       true,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName, flags, freeform := parseRunArgs(args)

		// Extract known flags
		renderFull := flags["render"] == "true"
		if v, ok := flags["dir"]; ok {
			dir = v
		}
		delete(flags, "render")
		delete(flags, "dir")

		// Handle --help
		if _, ok := flags["help"]; ok {
			return cmd.Help()
		}
		delete(flags, "help")

		if skillName == "" {
			fmt.Println("Usage: skills run <skill-name> [--render] [--key value] [-- freeform]")
			return fmt.Errorf("skill name required")
		}

		skills, invalid := skill.GetSkills(dir)
		s := skill.FindSkill(skillName, skills)

		if s == nil {
			if errs, ok := invalid[skillName]; ok {
				fmt.Printf("Cannot run '%s': validation errors\n", skillName)
				for _, e := range errs {
					fmt.Printf("  - %s\n", e)
				}
				return fmt.Errorf("skill has validation errors")
			}
			fmt.Printf("Skill '%s' not found.\n", skillName)
			fmt.Println("Run `skills list` to see available skills.")
			return fmt.Errorf("skill not found")
		}

		var output string
		var errors []string

		if renderFull {
			output, errors = skill.Render(s, flags, freeform)
		} else {
			output, errors = skill.Pointer(s, flags, freeform)
		}

		if len(errors) > 0 {
			fmt.Printf("Cannot run '%s':\n", skillName)
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
			return fmt.Errorf("render failed")
		}

		fmt.Print(output)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// parseRunArgs parses the run command args into (skillName, flags, freeform).
// Format: <skill-name> [--key value ...] [-- freeform...]
func parseRunArgs(args []string) (string, map[string]string, string) {
	flags := make(map[string]string)
	var skillName string
	var freeformParts []string

	i := 0
	// First non-flag arg is the skill name
	for i < len(args) {
		if args[i] == "--" {
			freeformParts = args[i+1:]
			break
		}
		if strings.HasPrefix(args[i], "--") {
			break
		}
		if skillName == "" {
			skillName = args[i]
		}
		i++
	}

	// Parse flags
	for i < len(args) {
		arg := args[i]
		if arg == "--" {
			freeformParts = args[i+1:]
			break
		} else if strings.HasPrefix(arg, "--") && strings.Contains(arg, "=") {
			kv := strings.SplitN(arg[2:], "=", 2)
			flags[kv[0]] = kv[1]
		} else if strings.HasPrefix(arg, "--") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
			flags[arg[2:]] = args[i+1]
			i++
		} else if strings.HasPrefix(arg, "--") {
			flags[arg[2:]] = "true"
		}
		i++
	}

	freeform := strings.Join(freeformParts, " ")
	return skillName, flags, freeform
}
