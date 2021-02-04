package zk

import (
	"encoding/binary"
	"errors"
	"filippo.io/edwards25519"

	"crypto/sha512"

	"github.com/taurusgroup/tg-tss/pkg/helpers/common"
)

type Schnorr struct {
	commitment edwards25519.Point  // commitment = v•G for random v
	response   edwards25519.Scalar // response = v - privateInput * challenge
}

var edwards25519GeneratorBytes = edwards25519.NewGeneratorPoint().Bytes()

func computeChallenge(commitmentPublic, public *edwards25519.Point, partyID uint32, params string) *edwards25519.Scalar {
	var partyIDBytes [4]byte
	binary.BigEndian.PutUint32(partyIDBytes[:], partyID)

	// Compute challenge
	// c = H(G || V || public || partyID || params)
	h := sha512.New()
	h.Write(edwards25519GeneratorBytes)
	h.Write(commitmentPublic.Bytes())
	h.Write(public.Bytes())
	h.Write(partyIDBytes[:])
	h.Write([]byte(params))

	challenge := new(edwards25519.Scalar).SetUniformBytes(h.Sum(nil))

	return challenge
}

// NewSchnorr is generates a ZK proof of knowledge of privateInput.
// Follows https://tools.ietf.org/html/rfc8235#section-3
func NewSchnorrProof(private *edwards25519.Scalar, partyID uint32, params string) (proof *Schnorr, public *edwards25519.Point) {
	var p edwards25519.Point
	proof = new(Schnorr)

	// public = x•G
	p.ScalarBaseMult(private)
	public = &p

	// Compute commitment for random nonce
	// V = v•G
	commitmentSecret := common.NewScalarRandom()      // = v
	proof.commitment.ScalarBaseMult(commitmentSecret) // V = v•G

	challenge := computeChallenge(&proof.commitment, public, partyID, params)

	proof.response.Multiply(challenge, private)                // = c•private
	proof.response.Subtract(commitmentSecret, &proof.response) // r = v - c•private

	return proof, public
}

// Schnorr.Verify verifies that the zero knowledge proof is valid.
// Follows https://tools.ietf.org/html/rfc8235#section-3
func (proof *Schnorr) Verify(public *edwards25519.Point, partyID uint32, params string) bool {
	var commitmentComputed edwards25519.Point

	// Check that the public point is not the identity
	if public.Equal(edwards25519.NewIdentityPoint()) == 1 {
		return false
	}
	// TODO: Check cofactor?

	challenge := computeChallenge(&proof.commitment, public, partyID, params)

	commitmentComputed.VarTimeDoubleScalarBaseMult(challenge, public, &proof.response) // = r•G + c•Public

	return commitmentComputed.Equal(&proof.commitment) == 1
}

func (proof *Schnorr) BytesAppend(existing []byte) (data []byte, err error) {
	existing = append(existing, proof.commitment.Bytes()...)
	existing = append(existing, proof.response.Bytes()...)
	return existing, nil
}

func (proof *Schnorr) MarshalBinary() (data []byte, err error) {
	var buf [64]byte
	return proof.BytesAppend(buf[:0])
}

func (proof *Schnorr) UnmarshalBinary(data []byte) error {
	if len(data) != 64 {
		return errors.New("length is wrong")
	}
	var err error
	_, err = proof.commitment.SetBytes(data[:32])
	if err != nil {
		return err
	}
	_, err = proof.response.SetCanonicalBytes(data[32:])
	if err != nil {
		return err
	}
	return nil
}

func (proof *Schnorr) Size() int {
	return 64
}
