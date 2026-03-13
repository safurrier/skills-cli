// Package skill provides types and functions for parsing, discovering,
// and rendering Agent Skills (SKILL.md files).
package skill

// RenderInput describes a single parameter a skill accepts.
type RenderInput struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Positional  bool   `yaml:"positional"`
}

// Properties holds the parsed YAML frontmatter of a SKILL.md.
type Properties struct {
	Name          string
	Description   string
	License       string
	Compatibility string
	AllowedTools  string
	Metadata      map[string]string

	// Parsed from metadata (skills-cli.* keys)
	RenderMode   string // "template" | "passthrough"
	RenderInputs []RenderInput
}

// Skill is a fully discovered skill with its directory context.
type Skill struct {
	Dir        string // absolute path to skill directory
	SkillMD    string // absolute path to SKILL.md
	Props      Properties
	Body       string   // markdown body after frontmatter
	Scripts    []string // paths in scripts/
	References []string // paths in references/
	Assets     []string // paths in assets/
}

// ResourceCount returns the total number of resource files.
func (s *Skill) ResourceCount() int {
	return len(s.Scripts) + len(s.References) + len(s.Assets)
}
