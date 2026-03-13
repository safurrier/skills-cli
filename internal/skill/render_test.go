package skill

import (
	"strings"
	"testing"
)

func TestValidateInputs(t *testing.T) {
	t.Parallel()

	declared := []RenderInput{
		{Name: "path", Required: true},
		{Name: "format", Default: "markdown"},
	}

	// Missing required input
	errs := ValidateInputs(declared, map[string]string{"format": "json"})
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}

	// All provided
	errs = ValidateInputs(declared, map[string]string{"path": "src/"})
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %v", len(errs), errs)
	}
}

func TestRender_TemplateMode(t *testing.T) {
	t.Parallel()

	s := &Skill{
		Dir:     "/tmp/test-skill",
		SkillMD: "/tmp/test-skill/SKILL.md",
		Props: Properties{
			Name:        "test-skill",
			Description: "A test skill",
			RenderMode:  "template",
			RenderInputs: []RenderInput{
				{Name: "path", Required: true},
			},
		},
		Body: "Review the code at {{path}}.",
	}

	output, errs := Render(s, map[string]string{"path": "src/"}, "")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !strings.Contains(output, "Review the code at src/.") {
		t.Errorf("template substitution failed, got:\n%s", output)
	}
}

func TestRender_PassthroughMode(t *testing.T) {
	t.Parallel()

	s := &Skill{
		Dir:     "/tmp/test-skill",
		SkillMD: "/tmp/test-skill/SKILL.md",
		Props: Properties{
			Name:        "test-skill",
			Description: "A test skill",
			RenderMode:  "passthrough",
			RenderInputs: []RenderInput{
				{Name: "path"},
			},
		},
		Body: "Do something.",
	}

	output, errs := Render(s, map[string]string{"path": "src/"}, "Extra context")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !strings.Contains(output, "## Invocation Parameters") {
		t.Error("expected Invocation Parameters block")
	}
	if !strings.Contains(output, "## Additional Context") {
		t.Error("expected Additional Context block")
	}
}
