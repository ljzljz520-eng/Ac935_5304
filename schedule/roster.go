package schedule

import (
	"hospitaldesk/model"
	"sort"
	"strings"
)

func GroupByPeriod(files []model.ScheduleFile) map[string][]model.ScheduleFile {
	groups := make(map[string][]model.ScheduleFile)
	for _, f := range files {
		groups[f.Period] = append(groups[f.Period], f)
	}
	for period := range groups {
		sort.Slice(groups[period], func(i, j int) bool { return groups[period][i].Department < groups[period][j].Department })
	}
	return groups
}

func NormalizePeriod(period string) string { return strings.ToUpper(strings.TrimSpace(period)) }

func IsCurrentPeriod(file model.ScheduleFile, period string) bool {
	return NormalizePeriod(file.Period) == NormalizePeriod(period)
}
