package skill

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Pointer generates a location-based prompt that points the agent to the skill
// rather than rendering the full content inline.
// Returns (pointer_text, errors).
func Pointer(skill *Skill, inputs map[string]string, freeform string) (string, []string) {
	// Validate inputs
	errors := ValidateInputs(skill.Props.RenderInputs, inputs)
	if len(errors) > 0 {
		return "", errors
	}

	fullInputs := resolveInputs(skill.Props.RenderInputs, inputs)

	var parts []string

	// Title
	parts = append(parts, fmt.Sprintf("# Skill: %s", skill.Props.Name))
	parts = append(parts, "")

	// Location pointers
	parts = append(parts, fmt.Sprintf("**Location**: %s", skill.SkillMD))
	parts = append(parts, fmt.Sprintf("**Directory**: %s/", skill.Dir))
	parts = append(parts, "")

	// Description as blockquote
	parts = append(parts, fmt.Sprintf("> %s", skill.Props.Description))
	parts = append(parts, "")

	// Instruction to load
	parts = append(parts, "Load and follow the instructions in the skill file above.")
	parts = append(parts, "")

	// Parameters
	if len(fullInputs) > 0 {
		parts = append(parts, "## Parameters")
		for _, inp := range skill.Props.RenderInputs {
			if val, ok := fullInputs[inp.Name]; ok {
				parts = append(parts, fmt.Sprintf("- **%s**: %s", inp.Name, val))
			}
		}
		// Ad-hoc flags
		declaredNames := make(map[string]bool)
		for _, inp := range skill.Props.RenderInputs {
			declaredNames[inp.Name] = true
		}
		for k, v := range fullInputs {
			if !declaredNames[k] {
				parts = append(parts, fmt.Sprintf("- **%s**: %s", k, v))
			}
		}
		parts = append(parts, "")
	}

	// Freeform context
	if strings.TrimSpace(freeform) != "" {
		parts = append(parts, "## Additional Context")
		parts = append(parts, strings.TrimSpace(freeform))
		parts = append(parts, "")
	}

	// Resources
	if skill.ResourceCount() > 0 {
		parts = append(parts, "## Resources")
		for _, ref := range skill.References {
			rel, _ := filepath.Rel(skill.Dir, ref)
			parts = append(parts, fmt.Sprintf("- %s", rel))
		}
		for _, script := range skill.Scripts {
			rel, _ := filepath.Rel(skill.Dir, script)
			parts = append(parts, fmt.Sprintf("- %s", rel))
		}
		for _, asset := range skill.Assets {
			rel, _ := filepath.Rel(skill.Dir, asset)
			parts = append(parts, fmt.Sprintf("- %s", rel))
		}
		parts = append(parts, "")
	}

	return strings.Join(parts, "\n"), nil
}
