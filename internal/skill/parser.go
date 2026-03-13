package skill

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v3"
)

// Spec validation limits.
const (
	MaxName          = 64
	MaxDescription   = 1024
	MaxCompatibility = 500
)

// allowedFields is the set of valid top-level frontmatter fields per the spec.
var allowedFields = map[string]bool{
	"name":          true,
	"description":   true,
	"license":       true,
	"allowed-tools": true,
	"metadata":      true,
	"compatibility": true,
}

// ParseError represents a SKILL.md parsing error.
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string { return e.Message }

// parseFrontmatter splits a SKILL.md into YAML frontmatter and body.
func parseFrontmatter(content string) (map[string]interface{}, string, error) {
	if !strings.HasPrefix(content, "---") {
		return nil, "", &ParseError{"SKILL.md must start with YAML frontmatter (---)"}
	}
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, "", &ParseError{"SKILL.md frontmatter not properly closed with ---"}
	}

	var meta map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &meta); err != nil {
		return nil, "", &ParseError{fmt.Sprintf("Invalid YAML in frontmatter: %v", err)}
	}
	if meta == nil {
		meta = make(map[string]interface{})
	}
	return meta, strings.TrimSpace(parts[2]), nil
}

// validateName checks that a skill name conforms to the spec.
func validateName(name, dirName string) []string {
	var errors []string
	normalized := norm.NFKC.String(strings.TrimSpace(name))

	if len(normalized) > MaxName {
		errors = append(errors, fmt.Sprintf("name exceeds %d chars (%d)", MaxName, len(normalized)))
	}
	if normalized != strings.ToLower(normalized) {
		errors = append(errors, "name must be lowercase")
	}
	if strings.HasPrefix(normalized, "-") || strings.HasSuffix(normalized, "-") {
		errors = append(errors, "name cannot start or end with hyphen")
	}
	if strings.Contains(normalized, "--") {
		errors = append(errors, "name cannot contain consecutive hyphens")
	}
	for _, r := range normalized {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			errors = append(errors, "name contains invalid characters (only a-z, 0-9, hyphens)")
			break
		}
	}
	if dirName != normalized {
		errors = append(errors, fmt.Sprintf("directory name '%s' must match skill name '%s'", dirName, normalized))
	}
	return errors
}

// parseRenderInputs parses the skills-cli.inputs YAML string into RenderInput list.
func parseRenderInputs(raw string) []RenderInput {
	var items []map[string]interface{}
	if err := yaml.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}

	var result []RenderInput
	for _, item := range items {
		nameVal, ok := item["name"]
		if !ok {
			continue
		}
		ri := RenderInput{
			Name: fmt.Sprintf("%v", nameVal),
			Type: "string",
		}
		if v, ok := item["type"]; ok {
			ri.Type = fmt.Sprintf("%v", v)
		}
		if v, ok := item["required"]; ok {
			ri.Required = fmt.Sprintf("%v", v) == "true"
		}
		if v, ok := item["description"]; ok {
			ri.Description = fmt.Sprintf("%v", v)
		}
		if v, ok := item["default"]; ok {
			ri.Default = fmt.Sprintf("%v", v)
		}
		if v, ok := item["positional"]; ok {
			ri.Positional = fmt.Sprintf("%v", v) == "true"
		}
		result = append(result, ri)
	}
	return result
}

// parseProps extracts Properties from the frontmatter map.
func parseProps(meta map[string]interface{}, dirName string) (Properties, []string) {
	var errors []string

	// Check for unexpected fields
	for k := range meta {
		if !allowedFields[k] {
			errors = append(errors, fmt.Sprintf("unexpected frontmatter field: %s", k))
		}
	}

	// Required fields
	if _, ok := meta["name"]; !ok {
		errors = append(errors, "missing required field: name")
	} else {
		errors = append(errors, validateName(fmt.Sprintf("%v", meta["name"]), dirName)...)
	}

	if _, ok := meta["description"]; !ok {
		errors = append(errors, "missing required field: description")
	} else if len(fmt.Sprintf("%v", meta["description"])) > MaxDescription {
		errors = append(errors, fmt.Sprintf("description exceeds %d chars", MaxDescription))
	}

	if v, ok := meta["compatibility"]; ok {
		if len(fmt.Sprintf("%v", v)) > MaxCompatibility {
			errors = append(errors, fmt.Sprintf("compatibility exceeds %d chars", MaxCompatibility))
		}
	}

	// Parse metadata for skills-cli extensions
	rawMetadata := make(map[string]string)
	renderMode := "passthrough"
	var renderInputs []RenderInput

	if m, ok := meta["metadata"]; ok {
		if mMap, ok := m.(map[string]interface{}); ok {
			for k, v := range mMap {
				rawMetadata[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	if v, ok := rawMetadata["skills-cli.render"]; ok {
		renderMode = v
	}
	if v, ok := rawMetadata["skills-cli.inputs"]; ok {
		renderInputs = parseRenderInputs(v)
	}

	props := Properties{
		Name:         strings.TrimSpace(fmt.Sprintf("%v", meta["name"])),
		Description:  strings.TrimSpace(fmt.Sprintf("%v", meta["description"])),
		Metadata:     rawMetadata,
		RenderMode:   renderMode,
		RenderInputs: renderInputs,
	}

	if v, ok := meta["license"]; ok {
		props.License = fmt.Sprintf("%v", v)
	}
	if v, ok := meta["compatibility"]; ok {
		props.Compatibility = fmt.Sprintf("%v", v)
	}
	if v, ok := meta["allowed-tools"]; ok {
		props.AllowedTools = fmt.Sprintf("%v", v)
	}

	return props, errors
}

// ParseSkillMD parses a SKILL.md file content and returns Properties, body, and errors.
func ParseSkillMD(content string, dirName string) (Properties, string, []string) {
	meta, body, err := parseFrontmatter(content)
	if err != nil {
		return Properties{}, "", []string{err.Error()}
	}

	props, errors := parseProps(meta, dirName)
	return props, body, errors
}

// LoadSkill loads and validates a skill from a directory path.
func LoadSkill(skillDir string) (*Skill, []string) {
	skillDir, _ = filepath.Abs(skillDir)
	skillMD := findSkillMD(skillDir)
	if skillMD == "" {
		return nil, []string{fmt.Sprintf("missing SKILL.md in %s", skillDir)}
	}

	content, err := readFile(skillMD)
	if err != nil {
		return nil, []string{fmt.Sprintf("cannot read %s: %v", skillMD, err)}
	}

	dirName := filepath.Base(skillDir)
	props, body, errors := ParseSkillMD(content, dirName)

	skill := &Skill{
		Dir:        skillDir,
		SkillMD:    skillMD,
		Props:      props,
		Body:       body,
		Scripts:    scanSubdir(skillDir, "scripts"),
		References: scanSubdir(skillDir, "references"),
		Assets:     scanSubdir(skillDir, "assets"),
	}
	return skill, errors
}
