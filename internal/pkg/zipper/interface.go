package zipper

import "context"

type Zipper interface {
	Create(ctx context.Context, source, target string) error
}
