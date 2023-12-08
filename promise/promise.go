package promise

import "sync"

type Promise[T any] struct {
	wg  sync.WaitGroup
	res T
	err error
}

func NewPromise[T any](f func() (T, error)) *Promise[T] {
	p := &Promise[T]{}
	p.wg.Add(1)
	go func() {
		p.res, p.err = f()
		p.wg.Done()
	}()
	return p
}

func (p *Promise[T]) Then(then func(result T), err func(error)) {
	go func() {
		p.wg.Wait()
		if p.err != nil {
			err(p.err)
			return
		}
		then(p.res)
	}()
}

func (p *Promise[T]) Await() (T, error) {
	p.wg.Wait()
	return p.res, p.err
}
