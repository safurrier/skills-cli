package cmd

import (
	"fmt"
	"strings"

	"github.com/safurrier/skills-cli/internal/skill"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <skill-name>",
	Short: "Show detailed information about a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		skillName := args[0]
		skills, invalid := skill.GetSkills(dir)
		s := skill.FindSkill(skillName, skills)

		if s == nil {
			if errs, ok := invalid[skillName]; ok {
				fmt.Printf("Skill '%s' has validation errors:\n", skillName)
				for _, e := range errs {
					fmt.Printf("  - %s\n", e)
				}
				return fmt.Errorf("skill has validation errors")
			}
			fmt.Printf("Skill '%s' not found.\n", skillName)
			fmt.Println("Run `skills list` to see available skills.")
			return fmt.Errorf("skill not found")
		}

		p := s.Props
		fmt.Printf("Skill: %s\n", p.Name)
		fmt.Printf("Path:  %s\n\n", s.Dir)
		fmt.Printf("Description:\n  %s\n", p.Description)

		if p.License != "" {
			fmt.Printf("\nLicense: %s\n", p.License)
		}
		if p.Compatibility != "" {
			fmt.Printf("Compatibility: %s\n", p.Compatibility)
		}
		if p.AllowedTools != "" {
			fmt.Printf("Allowed Tools: %s\n", p.AllowedTools)
		}

		fmt.Printf("\nRender mode: %s\n", p.RenderMode)

		if len(p.RenderInputs) > 0 {
			fmt.Println("\nRender inputs:")
			for _, inp := range p.RenderInputs {
				req := ""
				if inp.Required {
					req = " (required)"
				}
				def := ""
				if inp.Default != "" {
					def = fmt.Sprintf(" [default: %s]", inp.Default)
				}
				fmt.Printf("  --%s  <%s>%s%s\n", inp.Name, inp.Type, req, def)
				if inp.Description != "" {
					fmt.Printf("      %s\n", inp.Description)
				}
			}
		}

		// Non-cli metadata
		if len(p.Metadata) > 0 {
			var nonCLI []string
			for k := range p.Metadata {
				if !strings.HasPrefix(k, "skills-cli.") {
					nonCLI = append(nonCLI, k)
				}
			}
			if len(nonCLI) > 0 {
				fmt.Println("\nMetadata:")
				for _, k := range nonCLI {
					fmt.Printf("  %s: %s\n", k, p.Metadata[k])
				}
			}
		}

		// Resources
		fmt.Println("\nResources:")
		if len(s.Scripts) > 0 {
			fmt.Printf("  scripts/ (%d files)\n", len(s.Scripts))
		}
		if len(s.References) > 0 {
			fmt.Printf("  references/ (%d files)\n", len(s.References))
		}
		if len(s.Assets) > 0 {
			fmt.Printf("  assets/ (%d files)\n", len(s.Assets))
		}
		if s.ResourceCount() == 0 {
			fmt.Println("  (none)")
		}

		// Body preview
		bodyLines := strings.Split(s.Body, "\n")
		preview := bodyLines
		if len(preview) > 8 {
			preview = preview[:8]
		}
		fmt.Println("\nInstructions preview:")
		for _, line := range preview {
			fmt.Printf("  %s\n", line)
		}
		if len(bodyLines) > 8 {
			fmt.Printf("  ... (%d more lines)\n", len(bodyLines)-8)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
