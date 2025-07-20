package application

func WrapMessageStreamer(streamFn func(eventType, data string) error) func(data string) error {
	return func(data string) error {
		return streamFn("message", data)
	}
}

