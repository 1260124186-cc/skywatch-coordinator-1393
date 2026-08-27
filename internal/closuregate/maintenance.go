package closuregate

import "context"

func (g *Gate) Reset(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.active = ""
	return nil
}

func (g *Gate) Active() string { return g.active }
