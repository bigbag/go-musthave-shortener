package url

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Task struct {
	ShortIDs []string
	UserID   string
}
type Queue struct {
	arr  []*Task
	mu   sync.Mutex
	cond *sync.Cond
	stop bool
}

func (q *Queue) close() {
	q.cond.L.Lock()
	q.stop = true
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *Queue) PopWait() (*Task, bool) {
	q.cond.L.Lock()

	for len(q.arr) == 0 && !q.stop {
		q.cond.Wait()
	}

	if q.stop {
		q.cond.L.Unlock()
		return nil, false
	}

	t := q.arr[0]
	q.arr = q.arr[1:]

	q.cond.L.Unlock()

	return t, true
}

type TaskPool struct {
	l          logrus.FieldLogger
	r          URLRepository
	workerPool []*TaskWorker
	wg         *sync.WaitGroup
	queue      *Queue
	total      chan int
}

func NewTaskPool(ctx context.Context, l logrus.FieldLogger, r URLRepository) *TaskPool {
	p := &TaskPool{l: l, r: r}
	p.workerPool = make([]*TaskWorker, 0, runtime.NumCPU())
	p.queue = p.newQueue()

	for i := 0; i < runtime.NumCPU(); i++ {
		p.workerPool = append(p.workerPool, p.newWorker(i))
	}

	ctx, cancel := context.WithCancel(ctx)
	g, _ := errgroup.WithContext(ctx)
	p.wg = &sync.WaitGroup{}

	for _, w := range p.workerPool {
		p.wg.Add(1)
		worker := w
		f := func() error {
			return worker.loop(ctx)
		}
		g.Go(f)
	}

	go func() {
		if err := g.Wait(); err != nil {
			p.l.Info("worker: pool error ", err)
		}
	}()
	go func() {
		p.wg.Wait()
		close(p.total)
		cancel()
	}()

	p.total = make(chan int)
	go func() {
		total := 0
		for c := range p.total {
			total = total + c
		}
	}()

	return p
}

func (p *TaskPool) newQueue() *Queue {
	q := Queue{}
	q.cond = sync.NewCond(&q.mu)
	q.stop = false
	return &q
}

func (p *TaskPool) newWorker(id int) *TaskWorker {
	p.l.Info("worker: init ", id)
	return &TaskWorker{id, p}
}

func (p *TaskPool) Push(userID string, shortIDs []string) error {
	if p.queue.stop {
		return errors.New("worker: queue was stopped")
	}

	p.queue.cond.L.Lock()
	defer p.queue.cond.L.Unlock()

	t := Task{UserID: userID, ShortIDs: shortIDs}
	p.queue.arr = append(p.queue.arr, &t)
	p.queue.cond.Signal()
	return nil
}

func (p *TaskPool) Close() {
	p.queue.close()
}

type TaskWorker struct {
	id   int
	pool *TaskPool
}

func (w *TaskWorker) loop(ctx context.Context) error {
	defer func() {
		w.pool.wg.Done()
		w.pool.queue.close()

		<-ctx.Done()
	}()

	for {
		t, ok := w.pool.queue.PopWait()
		if !ok {
			return nil
		}

		w.pool.l.Info("worker: new task ", w.id)
		if err := w.pool.r.DeleteUserURLs(t.UserID, t.ShortIDs); err != nil {
			w.pool.l.Info("worker: run to out from loop ")
			return err
		}

		w.pool.total <- len(t.ShortIDs)
	}
}
