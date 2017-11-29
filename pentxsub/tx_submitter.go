package pentxsub

type TransactionSubmitter interface {
	Submit(envelope string) (*Result, error)
}

type SubmissionResult interface{}

type Result struct {
	Err            error
	Hash           string
	LedgerSequence int32
	EnvelopeXDR    string
	ResultXDR      string
	ResultMetaXDR  string
}

func NewResult(err error, hash string, seq int32, env string, result string, resultMeta string) *Result {
	return &Result{err, hash, seq, env, result, resultMeta}
}
