package skill

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateInputs checks that all required inputs are provided.
func ValidateInputs(declared []RenderInput, provided map[string]string) []string {
	var errors []string
	for _, inp := range declared {
		if inp.Required && inp.Default == "" {
			if _, ok := provided[inp.Name]; !ok {
				errors = append(errors, fmt.Sprintf("required input '%s' not provided", inp.Name))
			}
		}
	}
	return errors
}

// resolveInputs merges provided inputs with defaults.
func resolveInputs(declared []RenderInput, provided map[string]string) map[string]string {
	full := make(map[string]string)
	for _, inp := range declared {
		if v, ok := provided[inp.Name]; ok {
			full[inp.Name] = v
		} else if inp.Default != "" {
			full[inp.Name] = inp.Default
		}
	}
	// Include any extra flags passed by caller
	for k, v := range provided {
		if _, ok := full[k]; !ok {
			full[k] = v
		}
	}
	return full
}

// fillTemplate does simple {{var}} placeholder substitution.
func fillTemplate(body string, inputs map[string]string) string {
	result := body
	for key, val := range inputs {
		result = strings.ReplaceAll(result, "{{"+key+"}}", val)
	}
	return result
}

// buildInputsBlock renders a normalized parameter block.
func buildInputsBlock(provided map[string]string, declared []RenderInput) string {
	if len(provided) == 0 {
		return ""
	}
	var lines []string
	lines = append(lines, "## Invocation Parameters", "")

	declaredNames := make(map[string]bool)
	for _, inp := range declared {
		declaredNames[inp.Name] = true
		val := provided[inp.Name]
		if val == "" {
			val = inp.Default
		}
		if val != "" {
			lines = append(lines, fmt.Sprintf("- **%s**: %s", inp.Name, val))
		}
	}
	// Ad-hoc flags not declared in inputs
	for k, v := range provided {
		if !declaredNames[k] {
			lines = append(lines, fmt.Sprintf("- **%s**: %s", k, v))
		}
	}
	return strings.Join(lines, "\n")
}

// Render renders a skill into a full prompt string (the "inline" mode).
// Returns (rendered_text, errors).
func Render(skill *Skill, inputs map[string]string, freeform string) (string, []string) {
	// Validate inputs
	errors := ValidateInputs(skill.Props.RenderInputs, inputs)
	if len(errors) > 0 {
		return "", errors
	}

	fullInputs := resolveInputs(skill.Props.RenderInputs, inputs)

	var parts []string

	// Header
	parts = append(parts, fmt.Sprintf("# Skill: %s\n", skill.Props.Name))
	parts = append(parts, fmt.Sprintf("_%s_\n", skill.Props.Description))
	if skill.Props.Compatibility != "" {
		parts = append(parts, fmt.Sprintf("> **Compatibility**: %s\n", skill.Props.Compatibility))
	}
	parts = append(parts, "---\n")

	// Instructions
	body := skill.Body
	if skill.Props.RenderMode == "template" && len(fullInputs) > 0 {
		body = fillTemplate(body, fullInputs)
	}
	parts = append(parts, body)

	// Parameters block
	if len(fullInputs) > 0 {
		block := buildInputsBlock(fullInputs, skill.Props.RenderInputs)
		if block != "" {
			parts = append(parts, "\n---\n")
			parts = append(parts, block)
		}
	}

	// Freeform context
	if strings.TrimSpace(freeform) != "" {
		parts = append(parts, "\n---\n")
		parts = append(parts, "## Additional Context\n")
		parts = append(parts, strings.TrimSpace(freeform))
	}

	// Resource references
	if len(skill.References) > 0 {
		parts = append(parts, "\n---\n")
		parts = append(parts, "## Referenced Files\n")
		for _, ref := range skill.References {
			rel, _ := filepath.Rel(skill.Dir, ref)
			parts = append(parts, fmt.Sprintf("- `%s`", rel))
		}
	}

	return strings.Join(parts, "\n"), nil
}
