package coverage

import "github.com/IcedElect/goverage/internal/utils"

type Coverage struct {
	// Coverage summary by statements.
	Statements CoverageItem

	// Coverage summary by lines.
	Lines CoverageItem

	// Coverage summary by functions.
	Functions CoverageItem

	// Coverage total summary
	TotalPercent float64
}

type CoverageItem struct {
	Total   int
	Covered int
	Percent float64
}

func NewCoverageItem(total, covered int) CoverageItem {
	return CoverageItem{
		Total:   total,
		Covered: covered,
		Percent: utils.Percent(int64(covered), int64(total)),
	}
}
