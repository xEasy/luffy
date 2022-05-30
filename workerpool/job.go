package workerpool

type JobFunc func(args ...any)

type Job struct {
	ID   string
	Func func(args ...any)
	Args []any
}
