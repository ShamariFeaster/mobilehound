package crypto

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/require"
	"mobilehound/v0-edwards"
	"mobilehound/v0-log"
)

var s = edwards.NewAES128SHA256Ed25519(false)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func TestPub64(t *testing.T) {
	b := &bytes.Buffer{}
	rand := s.Cipher([]byte("example"))
	p, _ := s.Point().Pick(nil, rand)
	log.ErrFatal(Write64Point(s, b, p))
	log.ErrFatal(Write64Point(s, b, p))
	p2, err := Read64Point(s, b)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
	p2, err = Read64Point(s, b)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
}

func TestScalar64(t *testing.T) {
	b := &bytes.Buffer{}
	rand := s.Cipher([]byte("example"))
	sc := s.Scalar().Pick(rand)
	log.ErrFatal(Write64Scalar(s, b, sc))
	log.ErrFatal(Write64Scalar(s, b, sc))
	s2, err := Read64Scalar(s, b)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
	s2, err = Read64Scalar(s, b)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
}

func TestPubHexStream(t *testing.T) {
	b := &bytes.Buffer{}
	rand := s.Cipher([]byte("example"))
	p, _ := s.Point().Pick(nil, rand)
	log.ErrFatal(WriteHexPoint(s, b, p))
	log.ErrFatal(WriteHexPoint(s, b, p))
	p2, err := ReadHexPoint(s, b)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
	p2, err = ReadHexPoint(s, b)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
}

func TestScalarHexStream(t *testing.T) {
	b := &bytes.Buffer{}
	rand := s.Cipher([]byte("example"))
	sc := s.Scalar().Pick(rand)
	log.ErrFatal(WriteHexScalar(s, b, sc))
	log.ErrFatal(WriteHexScalar(s, b, sc))
	s2, err := ReadHexScalar(s, b)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
	s2, err = ReadHexScalar(s, b)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
}

func TestPubHexString(t *testing.T) {
	rand := s.Cipher([]byte("example"))
	p, _ := s.Point().Pick(nil, rand)
	pstr, err := PointToStringHex(s, p)
	log.ErrFatal(err)
	p2, err := StringHexToPoint(s, pstr)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
}

func TestPub64String(t *testing.T) {
	rand := s.Cipher([]byte("example"))
	p, _ := s.Point().Pick(nil, rand)
	pstr, err := PointToString64(s, p)
	log.ErrFatal(err)
	p2, err := String64ToPoint(s, pstr)
	log.ErrFatal(err)
	require.Equal(t, p, p2)
}

func TestScalarHexString(t *testing.T) {
	rand := s.Cipher([]byte("example"))
	sc := s.Scalar().Pick(rand)
	scstr, err := ScalarToStringHex(s, sc)
	log.ErrFatal(err)
	s2, err := StringHexToScalar(s, scstr)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
}

func TestScalar64String(t *testing.T) {
	rand := s.Cipher([]byte("example"))
	sc := s.Scalar().Pick(rand)
	scstr, err := ScalarToString64(s, sc)
	log.ErrFatal(err)
	s2, err := String64ToScalar(s, scstr)
	log.ErrFatal(err)
	require.True(t, sc.Equal(s2))
}
