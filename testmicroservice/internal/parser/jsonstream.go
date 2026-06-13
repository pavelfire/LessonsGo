package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"lessonsgo/testmicroservice/internal/domain/port"
)

type Handler func(ctx context.Context, p port.Port) error

func StreamPorts(ctx context.Context, r io.Reader, handle Handler) (int, error) {
	dec := json.NewDecoder(r)

	token, err := dec.Token()
	if err != nil {
		return 0, fmt.Errorf("read opening token: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		return 0, fmt.Errorf("expected JSON array, got %v", token)
	}

	processed := 0
	for dec.More() {
		if err := ctx.Err(); err != nil {
			return processed, err
		}

		var p port.Port
		if err := dec.Decode(&p); err != nil {
			return processed, fmt.Errorf("decode port at index %d: %w", processed, err)
		}
		if err := handle(ctx, p); err != nil {
			return processed, err
		}
		processed++
	}

	token, err = dec.Token()
	if err != nil {
		return processed, fmt.Errorf("read closing token: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != ']' {
		return processed, fmt.Errorf("expected closing ']', got %v", token)
	}

	return processed, nil
}
