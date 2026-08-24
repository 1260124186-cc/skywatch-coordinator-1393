package service

import (
	"fmt"
	"sync/atomic"
)

type IDGenerator struct{ sequence atomic.Uint64 }

func NewIDGenerator() *IDGenerator { return &IDGenerator{} }
func (g *IDGenerator) Next(prefix string) string {
	n := g.sequence.Add(1)
	return fmt.Sprintf("%s-%06d", prefix, n)
}
