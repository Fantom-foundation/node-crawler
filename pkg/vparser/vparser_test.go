package vparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseVersionString(t *testing.T) {
	type ParseTestCase struct {
		name string
		args string
		want *ParsedInfo
	}

	var test_data = []ParseTestCase{
		{
			name: "single",
			args: "opera",
			want: &ParsedInfo{
				Name: "opera",
			},
		},
		{
			name: "perfect-case",
			args: "go-opera/v1.1.2-rc.6-825a85c9-1689192286/linux-amd64/go1.20.4",
			want: &ParsedInfo{
				Name: "go-opera",
				Version: Version{
					Major: 1,
					Minor: 1,
					Patch: 2,
					Tag:   "rc.6",
					Build: "825a85c9",
					Date:  "1689192286",
				},
				Os: OSInfo{
					Os:           "linux",
					Architecture: "amd64",
				},
				Language: LanguageInfo{
					Name:    "go",
					Version: "1.20.4",
				},
			},
		},
		{
			name: "without-build",
			args: "go-opera/v1.1.2-rc.6/linux-amd64/go1.19.3",
			want: &ParsedInfo{
				Name: "go-opera",
				Version: Version{
					Major: 1,
					Minor: 1,
					Patch: 2,
					Tag:   "rc.6",
					Build: "",
				},
				Os: OSInfo{
					Os:           "linux",
					Architecture: "amd64",
				},
				Language: LanguageInfo{
					Name:    "go",
					Version: "1.19.3",
				},
			},
		},
		{
			name: "extended-tag",
			args: "go-opera/v1.1.2-txtracing-rc.6.1-fef7ab09-1693339901/linux-amd64/go1.19.12",
			want: &ParsedInfo{
				Name: "go-opera",
				Version: Version{
					Major: 1,
					Minor: 1,
					Patch: 2,
					Tag:   "txtracing-rc.6.1",
					Build: "fef7ab09",
					Date:  "1693339901",
				},
				Os: OSInfo{
					Os:           "linux",
					Architecture: "amd64",
				},
				Language: LanguageInfo{
					Name:    "go",
					Version: "1.19.12",
				},
			},
		},
		{
			name: "windows",
			args: "erigon/v2021.06.5-alpha-a0694dd3/windows-x86_64/go1.16.5",
			want: &ParsedInfo{
				Name: "erigon",
				Version: Version{
					Major: 2021,
					Minor: 06,
					Patch: 5,
					Tag:   "alpha",
					Build: "a0694dd3",
				},
				Os: OSInfo{
					Os:           "windows",
					Architecture: "x86_64",
				},
				Language: LanguageInfo{
					Name:    "go",
					Version: "1.16.5",
				},
			},
		},
		{
			name: "rust",
			args: "OpenEthereum/v3.2.6-stable-f9f4926-20210514/x86_64-linux-gnu/rustc1.52.1",
			want: &ParsedInfo{
				Name: "openethereum",
				Version: Version{
					Major: 3,
					Minor: 2,
					Patch: 6,
					Tag:   "stable",
					Build: "f9f4926",
					Date:  "20210514",
				},
				// This doesn't work
				// Os: OSInfo{
				// 	Os: "linux",
				// 	Architecture: "x86_64",
				// },
				Language: LanguageInfo{
					Name:    "rustc",
					Version: "1.52.1",
				},
			},
		},
		{
			name: "with-label",
			args: "Q-Client/v1.0.8-stable/Geth/v1.10.8-stable-825470ee/linux-amd64/go1.16.15",
			want: nil,
		},
		{
			name: "with-enode",
			args: "Geth/enode://91a3c3d5e76b0acf05d9abddee959f1bcbc7c91537d2629288a9edd7a3df90acaa46ffba0e0e5d49a20598e0960ac458d76eb8fa92a1d64938c0a3a3d60f8be4@127.0.0.1:21000/v1.10.0-stable(quorum-v22.1.0)/linux-amd64/go1.17.2",
			want: nil,
		},
	}

	for i, tt := range test_data {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseVersionString(tt.args)
			require.Equal(t, tt.want, got, i)
		})
	}
}
