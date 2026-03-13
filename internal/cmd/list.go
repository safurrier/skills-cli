package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/safurrier/skills-cli/internal/skill"
	"github.com/spf13/cobra"
)

var listJSON bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List discovered skills",
	RunE: func(_ *cobra.Command, _ []string) error {
		skills, invalid := skill.GetSkills(dir)

		if listJSON {
			return listAsJSON(skills, invalid)
		}
		return listAsText(skills, invalid)
	},
}

func init() {
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(listCmd)
}

type jsonSkill struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Resources   int    `json:"resources"`
}

type jsonInvalid struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}

func listAsJSON(skills []*skill.Skill, invalid map[string][]string) error {
	out := make(map[string]interface{})
	items := make([]jsonSkill, 0, len(skills))
	for _, s := range skills {
		items = append(items, jsonSkill{
			Name:        s.Props.Name,
			Description: s.Props.Description,
			Path:        s.Dir,
			Resources:   s.ResourceCount(),
		})
	}
	out["skills"] = items

	if len(invalid) > 0 {
		inv := make([]jsonInvalid, 0, len(invalid))
		for name, errs := range invalid {
			inv = append(inv, jsonInvalid{Name: name, Errors: errs})
		}
		out["invalid"] = inv
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func listAsText(skills []*skill.Skill, invalid map[string][]string) error {
	if len(skills) == 0 && len(invalid) == 0 {
		fmt.Printf("No skills found in %s\n", dir)
		fmt.Println("  Scanned: .agents/skills/, .claude/skills/, skills/")
		return fmt.Errorf("no skills found")
	}

	fmt.Printf("Available skills (%s)\n\n", dir)
	for _, s := range skills {
		desc := s.Props.Description
		if len(desc) > 72 {
			desc = desc[:69] + "..."
		}
		res := ""
		if s.ResourceCount() > 0 {
			res = fmt.Sprintf("  [%d resources]", s.ResourceCount())
		}
		fmt.Printf("  %s%s\n", s.Props.Name, res)
		fmt.Printf("    %s\n\n", desc)
	}

	if len(invalid) > 0 {
		fmt.Printf("\n%d invalid skill(s):\n", len(invalid))
		for name, errs := range invalid {
			fmt.Printf("  %s: %s\n", name, errs[0])
		}
	}
	return nil
}
