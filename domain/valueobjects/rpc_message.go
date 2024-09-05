package valueobjects

type RPCMessage []byte

var (
	RPC_ACK  = NewRPCMessage("ack")
	RPC_NACK = NewRPCMessage("nack")
)

func NewRPCMessage(message string) RPCMessage {
	return []byte(message)
}

func (m RPCMessage) Export() []byte {
	return ([]byte)(m)
}
