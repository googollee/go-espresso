package espresso

type Brewing interface {
	Next()
}

type brewBinder struct {
	Name     string
	BindFunc bindFunc
}

type brew[ContextData any] struct {
	handlers     []Handler[ContextData]
	handlerIndex int
	ctx          *brewContext[ContextData]
}

func (b *brew[ContextData]) Next() {
	for b.handlerIndex < len(b.handlers) && !b.ctx.isAborted {
		handler := b.handlers[b.handlerIndex]
		b.handlerIndex++
		if err := handler(b.ctx); err != nil {
			if ig, ok := err.(HTTPIgnore); ok && ig.Ignore() {
				continue
			}
			b.ctx.isAborted = true
			b.ctx.error = err
			return
		}
	}
}
