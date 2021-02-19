package messages

import (
	"fmt"

	"github.com/taurusgroup/frost-ed25519/pkg/helpers/polynomial"
	"github.com/taurusgroup/frost-ed25519/pkg/helpers/zk"
)

type KeyGen1 struct {
	Proof       *zk.Schnorr
	Commitments *polynomial.Exponent
}

func NewKeyGen1(from uint32, Proof *zk.Schnorr, Commitments *polynomial.Exponent) *Message {
	return &Message{
		Type: MessageTypeKeyGen1,
		From: from,
		KeyGen1: &KeyGen1{
			Proof:       Proof,
			Commitments: Commitments,
		},
	}
}

func (m *KeyGen1) BytesAppend(existing []byte) ([]byte, error) {
	var err error

	if m.Proof == nil || m.Commitments == nil {
		return nil, fmt.Errorf("msg1: %w", ErrInvalidMessage)
	}

	existing, err = m.Proof.BytesAppend(existing)
	if err != nil {
		return nil, err
	}
	existing, err = m.Commitments.BytesAppend(existing)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func (m *KeyGen1) MarshalBinary() (data []byte, err error) {

	buf := make([]byte, 0, m.Size())

	return m.BytesAppend(buf[:0])
}

func (m *KeyGen1) UnmarshalBinary(data []byte) error {
	var err error
	var proof zk.Schnorr
	var commitments polynomial.Exponent

	//proof.
	err = proof.UnmarshalBinary(data[:64])
	if err != nil {
		return err
	}
	m.Proof = &proof
	data = data[64:]

	err = commitments.UnmarshalBinary(data[:])
	if err != nil {
		return err
	}
	m.Commitments = &commitments
	return nil
}

func (m *KeyGen1) Size() int {
	return m.Proof.Size() + m.Commitments.Size()
}
