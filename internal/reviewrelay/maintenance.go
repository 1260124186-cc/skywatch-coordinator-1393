package reviewrelay

import "context"

func (r *Relay) Reset(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
