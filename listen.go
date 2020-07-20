package logr

// listen concurrently works through the buffered messages channel
func listen(ms <-chan *Message) {
	for {
		m := <-ms

		mutex.RLock()
		for c, w := range writers {
			if m.Type&c.filter != m.Type {
				continue
			}
			_, err := w.Write(c.format(m))
			if err != nil {
				Errorf("failed to write message to Writer: %v", err)
			}
		}
		mutex.RUnlock()

		close(m.done)

		m.Reset()
		pool.Put(m)
	}
}
