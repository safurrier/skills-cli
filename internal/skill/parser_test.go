package skill

import (
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid frontmatter",
			content: "---\nname: test-skill\ndescription: A test skill\n---\nBody content",
			wantErr: false,
		},
		{
			name:    "missing opening delimiter",
			content: "name: test-skill\n---\nBody content",
			wantErr: true,
		},
		{
			name:    "missing closing delimiter",
			content: "---\nname: test-skill\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, err := parseFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		skillDir string
		wantErrs int
	}{
		{"valid-skill", "valid-skill", 0},
		{"UPPER", "UPPER", 1},         // must be lowercase
		{"-leading", "-leading", 1},   // no leading hyphen
		{"trailing-", "trailing-", 1}, // no trailing hyphen
		{"double--hyphen", "double--hyphen", 1},
		{"mismatch", "other-dir", 1},
		{"has space", "has space", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			errs := validateName(tt.name, tt.skillDir)
			if len(errs) < tt.wantErrs {
				t.Errorf("validateName(%q, %q) got %d errors, want at least %d", tt.name, tt.skillDir, len(errs), tt.wantErrs)
			}
		})
	}
}

func TestParseSkillMD(t *testing.T) {
	t.Parallel()

	content := `---
name: my-skill
description: A helpful skill
compatibility: Claude Code
metadata:
  skills-cli.render: template
  skills-cli.inputs: |
    - name: path
      type: string
      required: true
    - name: format
      type: string
      default: markdown
---
# Instructions

Do the thing with {{path}}.`

	props, body, errors := ParseSkillMD(content, "my-skill")
	if len(errors) > 0 {
		t.Fatalf("unexpected errors: %v", errors)
	}

	if props.Name != "my-skill" {
		t.Errorf("Name = %q, want %q", props.Name, "my-skill")
	}
	if props.RenderMode != "template" {
		t.Errorf("RenderMode = %q, want %q", props.RenderMode, "template")
	}
	if len(props.RenderInputs) != 2 {
		t.Fatalf("RenderInputs len = %d, want 2", len(props.RenderInputs))
	}
	if !props.RenderInputs[0].Required {
		t.Error("first input should be required")
	}
	if props.RenderInputs[1].Default != "markdown" {
		t.Errorf("second input default = %q, want %q", props.RenderInputs[1].Default, "markdown")
	}
	if body == "" {
		t.Error("body should not be empty")
	}
}

func TestParseSkillMD_MissingFields(t *testing.T) {
	t.Parallel()

	content := "---\n---\nBody"
	_, _, errors := ParseSkillMD(content, "test")
	if len(errors) == 0 {
		t.Error("expected errors for missing name and description")
	}
}
