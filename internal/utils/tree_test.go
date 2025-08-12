package utils

import (
	"path"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/cover"
)

func Test_GetProfilesTree(t *testing.T) {
	testCases := []struct {
		name     string
		profiles []*cover.Profile
		expected []Directory
	}{
		{
			name:     "Empty profiles",
			profiles: []*cover.Profile{},
			expected: nil,
		},
		{
			name: "Single profile",
			profiles: []*cover.Profile{makeProfile("internal/database/user_repository.go")},
			expected: []Directory{
				{
					Path: "/internal/database",
					Profiles: []*cover.Profile{
						makeProfile("internal/database/user_repository.go"),
					},
				},
			},
		},
		{
			name: "Multiple profiles in same directory",
			profiles: []*cover.Profile{
				makeProfile("internal/database/user_repository.go"),
				makeProfile("internal/database/order_repository.go"),
			},
			expected: []Directory{
				{
					Path: "/internal/database",
					Profiles: []*cover.Profile{
						makeProfile("internal/database/user_repository.go"),
						makeProfile("internal/database/order_repository.go"),
					},
				},
			},
		},
		{
			name: "Profiles in different directories",
			profiles: []*cover.Profile{
				makeProfile("internal/database/user_repository.go"),
				makeProfile("internal/database/article_repository.go"),
				makeProfile("internal/api/user_handler.go"),
			},
			expected: []Directory{
				{
					Path: "/internal/database",
					Profiles: []*cover.Profile{
						makeProfile("internal/database/user_repository.go"),
						makeProfile("internal/database/article_repository.go"),
					},
				},
				{
					Path: "/internal/api",
					Profiles: []*cover.Profile{
						makeProfile("internal/api/user_handler.go"),
					},
				},
			},
		},
		{
			name: "Profiles with nested directories",
			profiles: []*cover.Profile{
				makeProfile("internal/database/user_repository.go"),
				makeProfile("internal/database/subdir/article_repository.go"),
				makeProfile("internal/api/user_handler.go"),
			},
			expected: []Directory{
				{
					Path: "/internal/database",
					Profiles: []*cover.Profile{
						makeProfile("internal/database/user_repository.go"),
					},
				},
				{
					Path: "/internal/database/subdir",
					Profiles: []*cover.Profile{
						makeProfile("internal/database/subdir/article_repository.go"),
					},
				},
				{
					Path: "/internal/api",
					Profiles: []*cover.Profile{
						makeProfile("internal/api/user_handler.go"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetProfilesTree(tc.profiles)
			assert.Equal(t, sortDirectories(tc.expected), sortDirectories(result), "Expected and actual tree lengths should match")
		})
	}
}

func makeProfile(fileName string) *cover.Profile {
	return &cover.Profile{
		FileName: path.Join(GetModulePath(), fileName),
		Mode:     "set",
		Blocks:   []cover.ProfileBlock{},
	}
}

func sortDirectories(dirs []Directory) []Directory {
	// Sort directories by path for consistent order in tests
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Path < dirs[j].Path
	})

	// Sort profiles within each directory by file name
	for i := range dirs {
		sort.Slice(dirs[i].Profiles, func(a, b int) bool {
			return dirs[i].Profiles[a].FileName < dirs[i].Profiles[b].FileName
		})
	}

	return dirs
}