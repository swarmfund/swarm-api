package pentxsub

type Envelope interface {
	TxOpCounter
	TxOpKeyGetter
	TxSignaturesGetter
	TxOperationTypeGetter
	Envelope() string
}

type TxOpCounter interface {
	TxOperationsCount() int
}

type TxOpKeyGetter interface {
	TxOperationKey() (string, error)
}

type TxSignaturesGetter interface {
	TxSignatures() [][]byte
}

type TxOperationTypeGetter interface {
	TxOperationType(index int) (int32, error)
}
