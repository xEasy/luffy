package xiface

type IDataPack interface {
	GetHeadLen() uint32                // get package head len
	Pack(msg IMessage) ([]byte, error) // pack message method
	UnPack([]byte) (IMessage, error)   // unpack data method
	SetPackSize(size uint32)
}
