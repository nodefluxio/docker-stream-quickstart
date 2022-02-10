package zipper

import (
	"context"
	"os"
)

type Zipper interface {
	Create(ctx context.Context, source, target string) (*os.File, error)
}
