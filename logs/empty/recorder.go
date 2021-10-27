package empty

type Recorder struct{}

func (Recorder) Record(m string) {
}

func NewRecorder() *Recorder {
	return &Recorder{}
}
