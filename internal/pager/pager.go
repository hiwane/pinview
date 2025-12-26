package pager

import (
	"context"
	"os"

	"github.com/hiwane/pinview/internal/draw"
	"github.com/hiwane/pinview/internal/input"
	"github.com/hiwane/pinview/internal/term"
)

type Pager struct {
	model *Model
	input *input.Input
	draw  draw.Drawer
}

func New(model *Model, input *input.Input, draw draw.Drawer) *Pager {
	return &Pager{
		model: model,
		input: input,
		draw:  draw,
	}
}

func (p *Pager) Run(ctx context.Context, tty *os.File) error {
	resizeCh := term.WatchResize(ctx, tty)
	sigintCh := term.WatchInterrupt()
	inputCh := p.input.Runes()

	if true {
		p.draw.Render(p.model.View())
	}

	need_draw := true
	for {
		select {
		case size := <-resizeCh:
			p.model.SetHeight(size.Height)
			p.model.SetWidth(size.Width)
			p.model.SetHeight(size.Height)
			p.model.SetSizeUpdate(true)
			need_draw = true
		case <-sigintCh:
			return nil
		case key := <-inputCh:
			if p.model.Update(key) {
				return nil
			}
			need_draw = true
		}

		if need_draw {
			lines := p.model.View()
			p.draw.Render(lines)
			need_draw = false
			p.model.SetSizeUpdate(false)
			p.model.SetKey(0)
		}
	}
}
