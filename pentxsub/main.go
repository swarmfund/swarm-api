package pentxsub

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/db2/api"
	depkeypair "gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/go/network"
	"gitlab.com/swarmfund/go/xdr"
	horizon "gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/tokend/keypair"
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
	signer  keypair.Full
}

func New(q api.PenTXSubQI, horizon *horizon.Connector, signer keypair.Full) *System {
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
	transaction *xdr.TransactionEnvelope
	hashRaw     []byte
	hashHex     string
	err         error
}

func (s submission) Init() *submission {
	// crafting helper struct to abstract XDR
	if err := xdr.SafeUnmarshalBase64(s.envelope, &s.transaction); err != nil {
		s.err = errors.Wrap(err, "failed to unmarshal tx")
		return &s
	}

	// TODO pass proper passphrase
	rawhash, err := network.HashTransaction(&s.transaction.Tx, "passphrase")
	if err != nil {
		s.err = errors.Wrap(err, "failed to hash tx")
		return &s
	}
	s.hashHex = fmt.Sprintf("%x", rawhash)
	s.hashRaw = rawhash[:]

	s.pendingTX, err = s.system.q.TransactionByHash(s.hashHex)
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
	// TODO move submitter to proper keypair
	body, err := s.system.horizon.SubmitTXSignedVerbose(s.envelope, depkeypair.MustParse(s.system.signer.Seed()))
	if err == nil {
		// submit was successful, cleaning up transaction
		err := s.system.q.DeleteTransaction(s.hashHex)
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
		if len(s.transaction.Tx.Operations) < 1 {
			return nil, errors.New("expected at least one op")
		}
		opType := s.transaction.Tx.Operations[0].Body.Type
		s.pendingTX = api.NewPendingTransaction(int32(opType), s.envelope, s.hashHex, s.transaction.Tx.SourceAccount.Address())
		id, err := s.system.q.CreateTransaction(s.pendingTX)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create tx")
		}
		s.pendingTX.ID = id
	}
	return s.pendingTX, nil
}

func (s *submission) FindSigner() (*horizon.Signer, error) {
	signers, err := s.system.horizon.Signers(s.transaction.Tx.SourceAccount.Address())
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
		signerKP := keypair.MustParseAddress(signer.AccountID)
		for _, signature := range s.transaction.Signatures {
			err = signerKP.Verify(s.hashRaw, signature.Signature)
			if err == nil {
				return &signer, nil
			}
		}
	}
	return nil, nil
}
