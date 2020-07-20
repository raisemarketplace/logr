package logr

func init() {
	pool.New = func() interface{} {
		return &Message{}
	}
	SetBufferSize(10000)
}