package workerpool

type Job struct {
	ID   string
	Func func(args ...interface{})
	Args []any
}
