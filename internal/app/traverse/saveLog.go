package traverse

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type saveLog struct {
	SavedAt      string        `json:"savedAt"`
	Source       string        `json:"source"`
	Selector     string        `json:"selector"`
	Algorithm    string        `json:"algorithm"`
	Limit        int           `json:"limit"`
	Stats        StatsDTO      `json:"stats"`
	Log          []LogEntryDTO `json:"log"`
	Matches      []string      `json:"matches"`
	VisitedOrder []string      `json:"visitedOrder"`
}

func saveTraversalLog(req Request, stats StatsDTO, log []LogEntryDTO, matches []string, visited []string) error {
	saveTraversal := saveLog{
		SavedAt:      time.Now().Format(time.RFC3339),
		Source:       req.Source,
		Selector:     req.Selector,
		Algorithm:    req.Algorithm,
		Limit:        req.Limit,
		Stats:        stats,
		Log:          log,
		Matches:      matches,
		VisitedOrder: visited,
	}

	buf, err := json.MarshalIndent(saveTraversal, "", "  ")
	if err != nil {
		return err
	}

	baseDir := filepath.Join("docs", "traversalLogs")
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return err
	}

	ts := time.Now().Format("20060102-150405")
	historyPath := filepath.Join(baseDir, fmt.Sprintf("traversal-%s.json", ts))
	if err := os.WriteFile(historyPath, buf, 0o644); err != nil {
		return err
	}

	latestPath := filepath.Join(baseDir, "latest.json")
	if err := os.WriteFile(latestPath, buf, 0o644); err != nil {
		return err
	}

	return nil
}
