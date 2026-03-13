package skill

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", "sample-skills")
}

func TestDiscover(t *testing.T) {
	t.Parallel()

	root := testdataDir(t)
	results := Discover(root)

	if len(results) != 3 {
		t.Fatalf("Discover() got %d results, want 3", len(results))
	}

	// All should be valid
	for _, r := range results {
		if len(r.Errors) > 0 {
			t.Errorf("skill %q has errors: %v", r.Skill.Props.Name, r.Errors)
		}
	}
}

func TestGetSkills(t *testing.T) {
	t.Parallel()

	root := testdataDir(t)
	skills, invalid := GetSkills(root)

	if len(skills) != 3 {
		t.Errorf("GetSkills() got %d valid skills, want 3", len(skills))
	}
	if len(invalid) != 0 {
		t.Errorf("GetSkills() got %d invalid, want 0", len(invalid))
	}
}

func TestFindSkill(t *testing.T) {
	t.Parallel()

	root := testdataDir(t)
	skills, _ := GetSkills(root)

	s := FindSkill("code-review", skills)
	if s == nil {
		t.Fatal("FindSkill(code-review) returned nil")
	}
	if s.Props.Name != "code-review" {
		t.Errorf("Name = %q, want %q", s.Props.Name, "code-review")
	}

	if FindSkill("nonexistent", skills) != nil {
		t.Error("FindSkill(nonexistent) should return nil")
	}
}
