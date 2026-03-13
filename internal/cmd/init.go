package cmd

import (
	"fmt"

	"github.com/safurrier/skills-cli/internal/skill"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scan and validate skills in a directory",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("Scanning %s for Agent Skills...\n", dir)
		skills, invalid := skill.GetSkills(dir)
		total := len(skills) + len(invalid)

		fmt.Printf("\nFound %d skill director", total)
		if total == 1 {
			fmt.Println("y:")
		} else {
			fmt.Println("ies:")
		}
		fmt.Printf("  %d valid\n", len(skills))
		if len(invalid) > 0 {
			fmt.Printf("  %d invalid\n", len(invalid))
			for name, errs := range invalid {
				fmt.Printf("\n  %s:\n", name)
				for _, e := range errs {
					fmt.Printf("    - %s\n", e)
				}
			}
		}

		if len(skills) > 0 {
			fmt.Println("\nValid skills:")
			for _, s := range skills {
				fmt.Printf("  ✓ %s  %s\n", s.Props.Name, s.Dir)
			}
		}

		fmt.Println("\nDiscovery paths checked:")
		for _, sub := range []string{".agents/skills", ".claude/skills", "skills"} {
			fmt.Printf("  %s/%s\n", dir, sub)
		}

		if len(invalid) > 0 {
			return fmt.Errorf("found %d invalid skill(s)", len(invalid))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
