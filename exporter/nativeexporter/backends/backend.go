package backends

type Writer interface {
	Write(bIndex []byte, bTraces []byte, bBloom []byte) error
}

type Reader interface {
	Find(traceId []byte) ([]byte, error)
}
