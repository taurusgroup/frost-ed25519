package sign

import (
	"filippo.io/edwards25519"
	"github.com/stretchr/testify/require"
	"github.com/taurusgroup/tg-tss/pkg/frost"
	"github.com/taurusgroup/tg-tss/pkg/helpers/common"
	"testing"
)

func TestMsg1_Encode_Decode(t *testing.T) {
	var d, e *edwards25519.Scalar
	var buf []byte
	var err error
	d = common.NewScalarRandom()
	e = common.NewScalarRandom()
	msg := Msg1{
		CommitmentD: new(edwards25519.Point).ScalarBaseMult(d),
		CommitmentE: new(edwards25519.Point).ScalarBaseMult(e),
	}

	realFrom := uint32(42)
	buf, err = msg.Encode(realFrom)
	require.NoError(t, err, "failed to encode")

	from, msgType, c := frost.DecodeBytes(buf)
	require.Equal(t, realFrom, from, "from not decoded")
	require.Equal(t, frost.MessageTypeSign1, msgType, "msgType not decoded")
	msg2, err := new(Msg1).Decode(c)
	require.Equal(t, 1, msg.CommitmentD.Equal(msg2.CommitmentD), "D not same")
	require.Equal(t, 1, msg.CommitmentE.Equal(msg2.CommitmentE), "E not same")
}

func TestMsg2_Encode_Decode(t *testing.T) {
	var s *edwards25519.Scalar
	var buf []byte
	var err error
	s = common.NewScalarRandom()
	msg := Msg2{
		SignatureShare: s,
	}
	realFrom := uint32(42)
	buf, err = msg.Encode(realFrom)
	require.NoError(t, err, "failed to encode")

	from, msgType, c := frost.DecodeBytes(buf)
	require.Equal(t, realFrom, from, "from not decoded")
	require.Equal(t, frost.MessageTypeSign2, msgType, "msgType not decoded")
	msg2, err := new(Msg2).Decode(c)
	require.Equal(t, 1, msg.SignatureShare.Equal(msg2.SignatureShare), "S not same")
}