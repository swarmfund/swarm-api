package pentxsub

import (
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"github.com/pkg/errors"
)

const (
	PENDING = 1
)

var (
	ErrTooManyOps = errors.New("too many operations")
)

type System struct {
	q       api.PenTXSubQI
	horizon *horizon.Connector
	signer  keypair.KP
}

func New(q api.PenTXSubQI, horizon *horizon.Connector, signer keypair.KP) *System {
	return &System{
		q:       q,
		horizon: horizon,
		signer:  signer,
	}
}

func (s *System) Submit(envelope string) ([]byte, error) {
	return submission{
		system:   s,
		envelope: envelope,
	}.Init().Process()
}

func (s *System) GetSigner(envelope string) (*horizon.Signer, error) {
	return submission{
		system:   s,
		envelope: envelope,
	}.Init().FindSigner()
}

// helper struct carrying state of submission
type submission struct {
	system *System

	envelope    string
	pendingTX   *api.PendingTransaction
	transaction *horizon.TransactionBuilder
	hash        *horizon.Hash
	err         error
}

func (s submission) Init() *submission {
	var err error
	// crafting helper struct to abstract XDR
	s.transaction = s.system.horizon.Transaction(&horizon.TransactionBuilder{
		Envelope: s.envelope,
	})

	s.hash, err = s.transaction.Hash()
	if err != nil {
		s.err = errors.Wrap(err, "failed to hash tx")
		return &s
	}

	s.pendingTX, err = s.system.q.TransactionByHash(s.hash.Hex())
	if err != nil {
		s.err = errors.Wrap(err, "failed to get pending tx")
		return &s
	}

	return &s
}

func (s *submission) Process() ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}

	// submit received tx as-is
	body, err := s.system.horizon.SubmitTXSignedVerbose(s.envelope, s.system.signer)
	if err == nil {
		// submit was successful, cleaning up transaction
		err := s.system.q.DeleteTransaction(s.hash.Hex())
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete transaction")
		}
		return body, nil
	}

	submitError, ok := errors.Cause(err).(horizon.SubmitError)
	if !ok {
		// submitter experienced technical issues
		return nil, err
	}

	// pending submit is restricted to a single operation
	if len(submitError.OperationCodes()) > 1 {
		return nil, ErrTooManyOps
	}

	// check if tx and op codes were not multi-sign related
	txCode := submitError.TransactionCode()
	opCodes := submitError.OperationCodes()
	txFailed := len(opCodes) == 0
	if (txFailed && txCode != "tx_bad_auth") || (!txFailed && opCodes[0] != "op_bad_auth") {
		// seems like transaction failed for non signature weight reason,
		// we are not going to update pending transaction/signers state.
		// by design it's admin responsibility to resolve transaction state by
		// either getting his shit together and producing valid signature or
		// deleting transaction altogether
		return nil, submitError
	}

	s.pendingTX, err = s.ensureTXSaved()
	if err != nil {
		return nil, errors.Wrap(err, "failed to save pending tx")
	}

	if err = s.saveSigner(); err != nil {
		return nil, errors.Wrap(err, "failed to save signer")
	}

	// returning original submission error, so whole pending flow should be implicit
	// for the client
	return nil, submitError
}

func (s *submission) saveSigner() error {
	signer, err := s.FindSigner()
	if signer == nil {
		return errors.Wrap(err, "failed to find signer")
	}

	return s.system.q.CreateTransactionSigner(&api.PendingTransactionSigner{
		PendingTransactionID: s.pendingTX.ID,
		SignerIdentity:       signer.Identity,
		SignerPublicKey:      signer.AccountID,
		SignerName:           signer.Name,
	})
}

func (s *submission) ensureTXSaved() (*api.PendingTransaction, error) {
	if s.pendingTX == nil {
		// transaction not yet in database, let's fix that
		// first we need to determine op type, checking ops length just in case
		if len(s.transaction.Operations) < 1 {
			return nil, errors.New("expected at least one op")
		}
		opType := s.transaction.Operations[0].Body.Type
		s.pendingTX = api.NewPendingTransaction(int32(opType), s.transaction.Envelope, s.hash.Hex(), s.transaction.Source.Address())
		id, err := s.system.q.CreateTransaction(s.pendingTX)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create tx")
		}
		s.pendingTX.ID = id
	}
	return s.pendingTX, nil
}

func (s *submission) FindSigner() (*horizon.Signer, error) {
	signers, err := s.system.horizon.Signers(s.transaction.Source.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get source signers")
	}

	var alreadySignedBy []api.PendingTransactionSigner
	if s.pendingTX != nil {
		alreadySignedBy, err = s.system.q.TransactionSigners(s.pendingTX.ID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get existing signers")
		}
	}

SIGNERS:
	for _, signer := range signers {
		for _, previousSigner := range alreadySignedBy {
			if previousSigner.SignerPublicKey == signer.AccountID || previousSigner.SignerIdentity == signer.Identity {
				continue SIGNERS
			}
		}
		signerKP := keypair.MustParse(signer.AccountID)
		for _, encoded := range s.transaction.Signatures {
			signature := xdr.DecoratedSignature{}
			err := xdr.SafeUnmarshalBase64(encoded, &signature)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal signature")
			}

			err = signerKP.Verify(s.hash.Slice(), signature.Signature)
			if err == nil {
				return &signer, nil
			}
		}
	}
	return nil, nil
}
