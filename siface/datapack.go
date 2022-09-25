package siface

type DataPack interface {
	GetHeadLen() uint32
	Pack(msg Message) ([]byte, error)
	Unpack([]byte) (Message, error)
}
