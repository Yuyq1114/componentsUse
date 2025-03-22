package stat

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type Stat struct {
	Ctx context.Context
	Eg  *errgroup.Group
}

func New(Ctx context.Context, eg *errgroup.Group) *Stat {
	return &Stat{
		Ctx: Ctx,
		Eg:  eg,
	}
}

func (s *Stat) Start() {
	s.Eg.Go(func() error {
		return s.GetStat()
	})
}

func (s *Stat) GetStat() error {
	return nil
}
