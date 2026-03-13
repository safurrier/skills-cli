package skill

import (
	"os"
	"path/filepath"
	"sort"
)

// discoveryPaths are the spec-defined locations to scan for skills.
var discoveryPaths = []string{
	".agents/skills",
	".claude/skills",
	"skills",
}

// findSkillMD looks for SKILL.md or skill.md in a directory.
func findSkillMD(dir string) string {
	for _, name := range []string{"SKILL.md", "skill.md"} {
		p := filepath.Join(dir, name)
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			return p
		}
	}
	return ""
}

// scanSubdir returns sorted file paths under dir/subdir.
func scanSubdir(dir, subdir string) []string {
	d := filepath.Join(dir, subdir)
	info, err := os.Stat(d)
	if err != nil || !info.IsDir() {
		return nil
	}

	var files []string
	_ = filepath.Walk(d, func(path string, fi os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !fi.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files
}

// readFile reads a file's contents as a string.
func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DiscoverResult holds a discovered skill and any validation errors.
type DiscoverResult struct {
	Skill  *Skill
	Errors []string
}

// Discover scans root for skill directories and returns results.
func Discover(root string) []DiscoverResult {
	root, _ = filepath.Abs(root)
	var candidateDirs []string
	seen := make(map[string]bool)

	// Spec-defined discovery locations
	for _, sub := range discoveryPaths {
		p := filepath.Join(root, sub)
		info, err := os.Stat(p)
		if err != nil || !info.IsDir() {
			continue
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				d := filepath.Join(p, e.Name())
				if !seen[d] {
					candidateDirs = append(candidateDirs, d)
					seen[d] = true
				}
			}
		}
	}

	// Also scan direct children of root
	entries, err := os.ReadDir(root)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			d := filepath.Join(root, e.Name())
			if !seen[d] && findSkillMD(d) != "" {
				candidateDirs = append(candidateDirs, d)
				seen[d] = true
			}
		}
	}

	sort.Strings(candidateDirs)

	var results []DiscoverResult
	for _, d := range candidateDirs {
		if findSkillMD(d) == "" {
			continue
		}
		skill, errors := LoadSkill(d)
		results = append(results, DiscoverResult{Skill: skill, Errors: errors})
	}
	return results
}

// GetSkills separates discover results into valid skills and invalid ones.
func GetSkills(root string) ([]*Skill, map[string][]string) {
	results := Discover(root)
	var valid []*Skill
	invalid := make(map[string][]string)

	for _, r := range results {
		if len(r.Errors) > 0 {
			name := "unknown"
			if r.Skill != nil {
				name = r.Skill.Props.Name
			}
			invalid[name] = r.Errors
		} else if r.Skill != nil {
			valid = append(valid, r.Skill)
		}
	}
	return valid, invalid
}

// FindSkill finds a skill by name in a list.
func FindSkill(name string, skills []*Skill) *Skill {
	for _, s := range skills {
		if s.Props.Name == name {
			return s
		}
	}
	return nil
}
