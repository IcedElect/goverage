package coverage

import (
	"testing"

	"github.com/IcedElect/goverage/internal/structure/files"
	"github.com/IcedElect/goverage/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/tools/cover"
)

func TestCalculator_CoverageByProfile(t *testing.T) {
	testCases := []struct {
		name     string
		profiles []*cover.Profile
		setup    func(filesRegistry *MockFilesRegistry)
		expected map[string]Coverage
	}{
		{
			name: "single profile",
			profiles: []*cover.Profile{
				{
					FileName: "internal/database/user_repository.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 1, EndLine: 1, NumStmt: 0, Count: 0},
						{StartLine: 1, EndLine: 10, NumStmt: 10, Count: 8},
						{StartLine: 20, EndLine: 25, NumStmt: 5, Count: 0},
					},
				},
			},
			expected: map[string]Coverage{
				"internal/database/user_repository.go": {
					Statements:   NewCoverageItem(15, 10),
					Lines:        NewCoverageItem(17, 10),
					Functions:    NewCoverageItem(1, 1),
					TotalPercent: utils.Percent(21, 33),
				},
			},
			setup: func(filesRegistry *MockFilesRegistry) {
				filesRegistry.EXPECT().GetFile("internal/database/user_repository.go").
					Return(&files.File{
						Funcs: []*utils.FuncExtent{
							{StartLine: 1, EndLine: 5},
						},
					}, true)
			},
		},
		{
			name: "multiple profiles",
			profiles: []*cover.Profile{
				{
					FileName: "internal/database/user_repository.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 1, EndLine: 1, NumStmt: 0, Count: 0},
						{StartLine: 1, EndLine: 10, NumStmt: 10, Count: 8},
						{StartLine: 20, EndLine: 25, NumStmt: 5, Count: 0},
					},
				},
				{
					FileName: "internal/service/user_service.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 1, EndLine: 5, NumStmt: 5, Count: 5},
						{StartLine: 10, EndLine: 15, NumStmt: 6, Count: 0},
					},
				},
				{
					FileName: "internal/database/user_repository.go",
					Blocks: []cover.ProfileBlock{
						{StartLine: 1, EndLine: 1, NumStmt: 0, Count: 0},
						{StartLine: 1, EndLine: 10, NumStmt: 10, Count: 8},
						{StartLine: 20, EndLine: 25, NumStmt: 5, Count: 0},
					},
				},
			},
			expected: map[string]Coverage{
				"internal/database/user_repository.go": {
					Statements:   NewCoverageItem(15, 10),
					Lines:        NewCoverageItem(17, 10),
					Functions:    NewCoverageItem(1, 1),
					TotalPercent: utils.Percent(21, 33),
				},
				"internal/service/user_service.go": {
					Statements:   NewCoverageItem(11, 5),
					Lines:        NewCoverageItem(11, 5),
					Functions:    NewCoverageItem(0, 0),
					TotalPercent: utils.Percent(10, 22),
				},
			},
			setup: func(filesRegistry *MockFilesRegistry) {
				filesRegistry.EXPECT().GetFile("internal/database/user_repository.go").
					Return(&files.File{
						Funcs: []*utils.FuncExtent{
							{StartLine: 1, EndLine: 5},
						},
					}, true)

				filesRegistry.EXPECT().GetFile("internal/service/user_service.go").
					Return(&files.File{
						Funcs: []*utils.FuncExtent{},
					}, true)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			filesRegistry := NewMockFilesRegistry(ctrl)
			calculator := NewCalculator(filesRegistry)

			tc.setup(filesRegistry)

			for _, profile := range tc.profiles {
				coverage := calculator.CoverageByProfile(profile)
				expectedCoverage, ok := tc.expected[profile.FileName]
				assert.True(t, ok, "expected coverage not found for file: %s", profile.FileName)
				assert.Equal(t, expectedCoverage, coverage, "coverage mismatch for file: %s", profile.FileName)
			}
		})
	}
}
