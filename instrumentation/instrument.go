package instrumentation

type Instrument interface {
	AddFile(filepath string) (err error)
	Instrument() (err error)
}
