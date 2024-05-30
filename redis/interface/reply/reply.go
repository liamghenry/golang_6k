package reply

type Reply interface {
	Marshal() []byte
}