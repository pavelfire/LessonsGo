package parser_test

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"lessonsgo/testmicroservice/internal/domain/port"
	"lessonsgo/testmicroservice/internal/parser"
	"lessonsgo/testmicroservice/internal/storage/memory"
)

func TestStreamPorts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := []struct {
		name       string
		input      string
		wantCount  int
		wantStored int
		wantErr    bool
	}{
		{
			name:       "empty array",
			input:      `[]`,
			wantCount:  0,
			wantStored: 0,
		},
		{
			name:       "single port",
			input:      `[{"id":"NLRTM","name":"Rotterdam"}]`,
			wantCount:  1,
			wantStored: 1,
		},
		{
			name: "latest duplicate wins in storage",
			input: `[
				{"id":"AEDXB","name":"Dubai"},
				{"id":"AEDXB","name":"Dubai (updated)"}
			]`,
			wantCount:  2,
			wantStored: 1,
		},
		{
			name:    "invalid json",
			input:   `{"id":"NLRTM"}`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := memory.New()
			svc := port.NewService(store)

			count, err := parser.StreamPorts(ctx, strings.NewReader(tc.input), func(ctx context.Context, p port.Port) error {
				return svc.Upsert(ctx, p)
			})

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
			assert.Equal(t, tc.wantStored, svc.Count())
		})
	}
}

func TestStreamPorts_RealFile(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	svc := port.NewService(store)

	f, err := openPortsFixture()
	require.NoError(t, err)
	defer f.Close()

	count, err := parser.StreamPorts(ctx, f, func(ctx context.Context, p port.Port) error {
		return svc.Upsert(ctx, p)
	})
	require.NoError(t, err)

	assert.Equal(t, 6, count)
	assert.Equal(t, 5, svc.Count())

	got, ok := svc.Get(ctx, "AEDXB")
	require.True(t, ok)
	assert.Equal(t, "Dubai (updated)", got.Name)
}

func openPortsFixture() (*os.File, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, os.ErrInvalid
	}

	path := filepath.Join(filepath.Dir(filename), "..", "..", "ports.json")
	return os.Open(path)
}
