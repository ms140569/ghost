package main

func ParseFrames(data []byte) []Frame {
	frames := []Frame{}

	Scanner(string(data))

	return frames
}
