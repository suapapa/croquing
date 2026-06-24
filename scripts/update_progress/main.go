// update_progress reads docs/progress/workitems.json and regenerates docs/progress/PROGRESS.md.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type WorkItemsFile struct {
	SchemaVersion int     `json:"schema_version"`
	Phases        []Phase `json:"phases"`
}

type Phase struct {
	ID            string `json:"id"`
	Label         string `json:"label"`
	Range         string `json:"range"`
	Prerequisite  string `json:"prerequisite,omitempty"`
	Items         []Item `json:"items"`
}

type Item struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Deps        []string `json:"deps"`
	Status      string   `json:"status"`
	CompletedAt string   `json:"completed_at,omitempty"`
	Commit      string   `json:"commit,omitempty"`
	Artifacts   []string `json:"artifacts,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

func main() {
	root, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "repo root: %v\n", err)
		os.Exit(1)
	}

	jsonPath := filepath.Join(root, "docs/progress/workitems.json")
	mdPath := filepath.Join(root, "docs/progress/PROGRESS.md")

	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read workitems: %v\n", err)
		os.Exit(1)
	}

	var file WorkItemsFile
	if err := json.Unmarshal(data, &file); err != nil {
		fmt.Fprintf(os.Stderr, "parse workitems: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().Format("2006-01-02 15:04 MST")
	var b strings.Builder

	b.WriteString("# Croquing — 구현 진도\n\n")
	b.WriteString("> **자동 생성 문서** — 직접 수정하지 마세요.  \n")
	b.WriteString("> 소스: [`workitems.json`](workitems.json) · 갱신: `make progress`\n\n")
	b.WriteString(fmt.Sprintf("마지막 갱신: **%s**\n\n", now))

	var totalDone, totalAll int
	var nextCandidates []Item

	for _, phase := range file.Phases {
		done, all := phaseStats(phase.Items)
		totalDone += done
		totalAll += all

		b.WriteString(fmt.Sprintf("## %s (%s)\n\n", phase.Label, phase.Range))
		if phase.Prerequisite != "" {
			b.WriteString(fmt.Sprintf("선행 조건: %s\n\n", phase.Prerequisite))
		}
		b.WriteString(fmt.Sprintf("진행: **%d / %d** (%s)\n\n", done, all, pct(done, all)))
		b.WriteString("| Index | 상태 | 제목 | deps | 산출물 검증 |\n")
		b.WriteString("|-------|------|------|------|-------------|\n")

		for _, item := range phase.Items {
			check := verifyArtifacts(root, item)
			b.WriteString(fmt.Sprintf("| **%s** | %s | %s | %s | %s |\n",
				item.ID,
				statusBadge(item.Status),
				item.Title,
				formatDeps(item.Deps),
				check,
			))
			if item.Status == "pending" || item.Status == "in_progress" {
				nextCandidates = append(nextCandidates, item)
			}
		}
		b.WriteString("\n")

		writeDoneDetails(&b, phase.Items)
	}

	b.WriteString("## 전체 요약\n\n")
	b.WriteString(fmt.Sprintf("- **전체:** %d / %d 완료 (%s)\n", totalDone, totalAll, pct(totalDone, totalAll)))
	b.WriteString(fmt.Sprintf("- **백엔드:** %s\n", phaseSummary(file, "backend")))
	b.WriteString(fmt.Sprintf("- **프론트엔드:** %s\n", phaseSummary(file, "frontend")))
	b.WriteString("\n")

	b.WriteString("## 다음 작업 후보\n\n")
	writeNextItems(&b, nextCandidates, file)

	b.WriteString("## 진도 갱신 방법\n\n")
	b.WriteString("1. [`workitems.json`](workitems.json)에서 해당 WorkItem의 `status`를 수정합니다.\n")
	b.WriteString("   - `pending` · `in_progress` · `done` · `blocked`\n")
	b.WriteString("2. 완료 시 `completed_at`, `commit`, `artifacts`, `notes`를 채웁니다.\n")
	b.WriteString("3. 저장 후 `make progress`를 실행합니다.\n\n")
	b.WriteString("상세 스펙·API·아키텍처는 [`../PROJECT_PLAN.md`](../PROJECT_PLAN.md)를 참고하세요.\n")

	if err := os.WriteFile(mdPath, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write progress: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("updated %s\n", mdPath)
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func phaseStats(items []Item) (done, all int) {
	for _, item := range items {
		all++
		if item.Status == "done" {
			done++
		}
	}
	return done, all
}

func pct(done, all int) string {
	if all == 0 {
		return "0%"
	}
	return fmt.Sprintf("%d%%", done*100/all)
}

func phaseSummary(file WorkItemsFile, id string) string {
	for _, phase := range file.Phases {
		if phase.ID == id {
			done, all := phaseStats(phase.Items)
			return fmt.Sprintf("%d / %d (%s)", done, all, pct(done, all))
		}
	}
	return "—"
}

func statusBadge(status string) string {
	switch status {
	case "done":
		return "✅ done"
	case "in_progress":
		return "🚧 in_progress"
	case "blocked":
		return "⛔ blocked"
	default:
		return "⬜ pending"
	}
}

func formatDeps(deps []string) string {
	if len(deps) == 0 {
		return "—"
	}
	return strings.Join(deps, ", ")
}

func verifyArtifacts(root string, item Item) string {
	if len(item.Artifacts) == 0 {
		return "—"
	}
	if item.Status != "done" {
		return "—"
	}

	var missing []string
	for _, artifact := range item.Artifacts {
		path := filepath.Join(root, artifact)
		if _, err := os.Stat(path); err != nil {
			missing = append(missing, artifact)
		}
	}
	if len(missing) == 0 {
		return "✓"
	}
	return "⚠ missing: " + strings.Join(missing, ", ")
}

func writeDoneDetails(b *strings.Builder, items []Item) {
	var done []Item
	for _, item := range items {
		if item.Status == "done" {
			done = append(done, item)
		}
	}
	if len(done) == 0 {
		return
	}

	b.WriteString("### 완료 항목 상세\n\n")
	for _, item := range done {
		b.WriteString(fmt.Sprintf("#### %s — %s\n\n", item.ID, item.Title))
		if item.CompletedAt != "" {
			b.WriteString(fmt.Sprintf("- 완료일: %s\n", item.CompletedAt))
		}
		if item.Commit != "" {
			b.WriteString(fmt.Sprintf("- 커밋: `%s`\n", item.Commit))
		}
		if len(item.Artifacts) > 0 {
			b.WriteString(fmt.Sprintf("- 산출물: %s\n", strings.Join(item.Artifacts, ", ")))
		}
		if item.Notes != "" {
			b.WriteString(fmt.Sprintf("- 메모: %s\n", item.Notes))
		}
		b.WriteString("\n")
	}
}

func writeNextItems(b *strings.Builder, candidates []Item, file WorkItemsFile) {
	doneSet := doneItems(file)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].ID < candidates[j].ID
	})

	var ready []Item
	for _, item := range candidates {
		if depsSatisfied(item.Deps, doneSet) {
			ready = append(ready, item)
		}
	}

	if len(ready) == 0 {
		b.WriteString("현재 선행 작업이 모두 완료된 후보가 없습니다.\n\n")
		return
	}

	limit := 5
	if len(ready) < limit {
		limit = len(ready)
	}
	for i := 0; i < limit; i++ {
		item := ready[i]
		b.WriteString(fmt.Sprintf("- **%s** — %s (`%s`)\n", item.ID, item.Title, item.Status))
	}
	b.WriteString("\n")
}

func doneItems(file WorkItemsFile) map[string]bool {
	m := make(map[string]bool)
	for _, phase := range file.Phases {
		for _, item := range phase.Items {
			if item.Status == "done" {
				m[item.ID] = true
			}
		}
	}
	return m
}

func depsSatisfied(deps []string, done map[string]bool) bool {
	for _, dep := range deps {
		if !done[dep] {
			return false
		}
	}
	return true
}
