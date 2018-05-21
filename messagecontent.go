package kait

type messageContent interface {
	MarshalKV() ([]byte, error)
	MarshalJSON() ([]byte, error)
}
