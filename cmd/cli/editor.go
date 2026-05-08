package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var knownIDFields = map[string]string{
	"auto_groups":              "group",
	"groups":                   "group",
	"distribution_groups":      "group",
	"peer_groups":              "group",
	"access_control_groups":    "group",
	"disabled_management_groups": "group",
	"source_posture_checks":    "posturecheck",
	"peers":                    "peer",
	"peer":                     "peer",
	"user_id":                  "user",
	"created_by":               "user",
	"routers":                  "peer",
	"policies":                 "policy",
	"sources":                  "group",
	"destinations":             "group",
	"zone_id":                  "dnszone",
	"network_id":               "network",
}

func annotateYAML(yamlBytes []byte, m map[string]interface{}) []byte {
	ids := collectIDs(m)
	if len(ids) == 0 {
		return yamlBytes
	}

	lines := strings.Split(string(yamlBytes), "\n")
	var out []string

	for _, line := range lines {
		annotated := line
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		for id, name := range ids {
			if name == "" {
				continue
			}
			if strings.Contains(line, id) {
				suffix := "  # " + name
				if !strings.Contains(line, suffix) {
					annotated = line + suffix
				}
				break
			}
		}
		out = append(out, annotated)
	}

	return []byte(strings.Join(out, "\n"))
}

func collectIDs(m map[string]interface{}) map[string]string {
	ids := make(map[string]string)

	for key, val := range m {
		resourceType, ok := knownIDFields[key]
		if !ok {
			continue
		}

		switch v := val.(type) {
		case string:
			if v != "" {
				ids[v] = c.IDToName(resourceType, v)
			}
		case []interface{}:
			for _, item := range v {
				switch s := item.(type) {
				case string:
					if s != "" {
						ids[s] = c.IDToName(resourceType, s)
					}
				case map[string]interface{}:
					if id, ok := s["id"].(string); ok && id != "" {
						ids[id] = c.IDToName(resourceType, id)
					}
				}
			}
		}
	}

	return ids
}

func stripComments(yamlBytes []byte) []byte {
	re := regexp.MustCompile(`  # .*$`)
	var lines []string
	for _, line := range strings.Split(string(yamlBytes), "\n") {
		lines = append(lines, re.ReplaceAllString(line, ""))
	}
	return []byte(strings.Join(lines, "\n"))
}

func createWithEditor(path string, template map[string]interface{}) (map[string]interface{}, error) {
	yamlData, err := yaml.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("YAML conversion: %w", err)
	}

	annotated := annotateYAML(yamlData, template)

	tmpFile, err := os.CreateTemp("", "netbird-create-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("# Creation: %s\n# Fill in the fields, save and quit.\n", path))
	buf.Write(annotated)

	original := buf.Bytes()

	if _, err := tmpFile.Write(original); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if _, err := exec.LookPath(editor); err != nil {
		return nil, fmt.Errorf("editor %s not found", editor)
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor: %w", err)
	}

	edited, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	if bytes.Equal(original, edited) {
		fmt.Println("Creation cancelled, no changes.")
		return nil, nil
	}

	cleaned := stripComments(edited)

	var result map[string]interface{}
	if err := yaml.Unmarshal(cleaned, &result); err != nil {
		return nil, fmt.Errorf("YAML parsing: %w", err)
	}

	resolveNamesToIDs(result)

	return result, nil
}

func resolveNamesToIDs(m map[string]interface{}) {
	for key, resourceType := range knownIDFields {
		switch v := m[key].(type) {
		case string:
			if v != "" {
				if id := resolveName(v, resourceType); id != "" {
					m[key] = id
				}
			}
		case []interface{}:
			var resolved []interface{}
			for _, item := range v {
				if s, ok := item.(string); ok {
					if id := resolveName(s, resourceType); id != "" {
						resolved = append(resolved, id)
					} else {
						resolved = append(resolved, s)
					}
				} else {
					resolved = append(resolved, item)
				}
			}
			m[key] = resolved
		}
	}
}

func resolveName(val, resourceType string) string {
	switch resourceType {
	case "group":
		if id, err := c.ResolveGroupID(val); err == nil {
			return id
		}
	case "user":
		if id, err := c.ResolveUserID(val); err == nil {
			return id
		}
	case "peer":
		if id, err := c.ResolvePeerID(val); err == nil {
			return id
		}
	case "network":
		if id, err := c.ResolveNetworkID(val); err == nil {
			return id
		}
	case "policy":
		if id, err := c.ResolvePolicyID(val); err == nil {
			return id
		}
	case "dnszone":
		if id, err := c.ResolveDNSZoneID(val); err == nil {
			return id
		}
	}
	return ""
}

func editWithEditor(path string, currentData []byte, overrides map[string]interface{}) (map[string]interface{}, error) {
	var current map[string]interface{}
	if err := json.Unmarshal(currentData, &current); err != nil {
		return nil, fmt.Errorf("JSON parsing: %w", err)
	}

	delete(current, "id")
	delete(current, "password")
	delete(current, "is_current")
	delete(current, "status")
	delete(current, "pending_approval")

	for k, v := range overrides {
		current[k] = v
	}

	yamlData, err := yaml.Marshal(current)
	if err != nil {
		return nil, fmt.Errorf("YAML conversion: %w", err)
	}

	annotated := annotateYAML(yamlData, current)

	tmpFile, err := os.CreateTemp("", "netbird-edit-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	var buf bytes.Buffer
	header := fmt.Sprintf("# Edit: %s\n# Deleted fields will be ignored.\n# Comments (# ...) after IDs are ignored.\n", path)
	buf.WriteString(header)
	buf.Write(annotated)

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if _, err := exec.LookPath(editor); err != nil {
		return nil, fmt.Errorf("editor %s not found (set EDITOR or install vim)", editor)
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("editor: %w", err)
	}

	edited, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	if bytes.Equal(buf.Bytes(), edited) {
		fmt.Println("Edit cancelled, no changes.")
		return nil, nil
	}

	cleaned := stripComments(edited)

	var updated map[string]interface{}
	if err := yaml.Unmarshal(cleaned, &updated); err != nil {
		return nil, fmt.Errorf("YAML parsing after edit: %w\n--- Cleaned YAML ---\n%s", err, string(cleaned))
	}

	resolveNamesToIDs(updated)

	for k := range updated {
		if k == "id" {
			delete(updated, k)
		}
	}

	if dryRun {
		fmt.Println("[dry-run] diff between original and edited version:")
		newYAML, _ := yaml.Marshal(updated)
		newAnnotated := annotateYAML(newYAML, updated)
		fmt.Println("--- original")
		fmt.Print(string(annotated))
		fmt.Println("--- modified")
		fmt.Print(string(newAnnotated))
		return nil, nil
	}

	result, err := c.PutRaw(path, updated)
	if err != nil {
		return nil, fmt.Errorf("PUT %s: %w", path, err)
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(result, &resultMap); err != nil {
		return nil, fmt.Errorf("response parsing: %w", err)
	}

	return resultMap, nil
}
