package stdout

import (
	"log"
)

type Recorder struct{}

func (*Recorder) Record(m string) {
	log.Println(m)
}

func NewRecorder() *Recorder {
	return &Recorder{}
}
