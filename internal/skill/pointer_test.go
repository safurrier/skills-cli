package skill

import (
	"strings"
	"testing"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	s := &Skill{
		Dir:     "/tmp/test-skill",
		SkillMD: "/tmp/test-skill/SKILL.md",
		Props: Properties{
			Name:        "test-skill",
			Description: "A test skill",
			RenderInputs: []RenderInput{
				{Name: "path", Required: true},
			},
		},
		Body: "Instructions here.",
		References: []string{
			"/tmp/test-skill/references/GUIDE.md",
		},
	}

	output, errs := Pointer(s, map[string]string{"path": "src/"}, "Focus on tests")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}

	checks := []string{
		"# Skill: test-skill",
		"**Location**: /tmp/test-skill/SKILL.md",
		"**Directory**: /tmp/test-skill/",
		"> A test skill",
		"Load and follow the instructions",
		"- **path**: src/",
		"## Additional Context",
		"Focus on tests",
		"## Resources",
		"references/GUIDE.md",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("output missing %q\nGot:\n%s", check, output)
		}
	}
}

func TestPointer_RequiredInputMissing(t *testing.T) {
	t.Parallel()

	s := &Skill{
		Props: Properties{
			RenderInputs: []RenderInput{
				{Name: "topic", Required: true},
			},
		},
	}

	_, errs := Pointer(s, map[string]string{}, "")
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}
