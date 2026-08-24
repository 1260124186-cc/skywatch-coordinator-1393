package reviewrelay

import "context"

func (r *Relay) Reset(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.active = ""
	return nil
}
func (r *Relay) Active() string { return r.active }
