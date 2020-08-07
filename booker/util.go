package booker

//Int64Result represent result with int64 data
type Int64Result struct {
	data int64
	err  error
}

func goroutine(f func(chan error)) chan error {
	out := make(chan error, 0)
	go f(out)
	return out
}

func goroutineInts(f func(chan Int64Result)) chan Int64Result {
	out := make(chan Int64Result, 0)
	go f(out)
	return out
}
