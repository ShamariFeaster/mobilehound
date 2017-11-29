package edwards

import (
	"crypto/cipher"
	"io"
	"math/big"

	"mobilehound/v0-abstract"
	"mobilehound/group"
	"mobilehound/nist"
)

type projPoint struct {
	X, Y, Z nist.Int
	c       *ProjectiveCurve
}

func (P *projPoint) initXY(x, y *big.Int, c abstract.Group) {
	P.c = c.(*ProjectiveCurve)
	P.X.Init(x, &P.c.P)
	P.Y.Init(y, &P.c.P)
	P.Z.Init64(1, &P.c.P)
}

func (P *projPoint) getXY() (x, y *nist.Int) {
	P.normalize()
	return &P.X, &P.Y
}

func (P *projPoint) String() string {
	P.normalize()
	return P.c.pointString(&P.X, &P.Y)
}

func (P *projPoint) MarshalSize() int {
	return P.c.PointLen()
}

func (P *projPoint) MarshalBinary() ([]byte, error) {
	P.normalize()
	return P.c.encodePoint(&P.X, &P.Y), nil
}

func (P *projPoint) UnmarshalBinary(b []byte) error {
	P.Z.Init64(1, &P.c.P)
	return P.c.decodePoint(b, &P.X, &P.Y)
}

func (P *projPoint) MarshalTo(w io.Writer) (int, error) {
	return group.PointMarshalTo(P, w)
}

func (P *projPoint) UnmarshalFrom(r io.Reader) (int, error) {
	return group.PointUnmarshalFrom(P, r)
}

func (P *projPoint) HideLen() int {
	return P.c.hide.HideLen()
}

func (P *projPoint) HideEncode(rand cipher.Stream) []byte {
	return P.c.hide.HideEncode(P, rand)
}

func (P *projPoint) HideDecode(rep []byte) {
	P.c.hide.HideDecode(P, rep)
}

// Equality test for two Points on the same curve.
// We can avoid inversions here because:
//
//	(X1/Z1,Y1/Z1) == (X2/Z2,Y2/Z2)
//		iff
//	(X1*Z2,Y1*Z2) == (X2*Z1,Y2*Z1)
//
func (P1 *projPoint) Equal(CP2 abstract.Point) bool {
	P2 := CP2.(*projPoint)
	var t1, t2 nist.Int
	xeq := t1.Mul(&P1.X, &P2.Z).Equal(t2.Mul(&P2.X, &P1.Z))
	yeq := t1.Mul(&P1.Y, &P2.Z).Equal(t2.Mul(&P2.Y, &P1.Z))
	return xeq && yeq
}

func (P *projPoint) Set(CP2 abstract.Point) abstract.Point {
	P2 := CP2.(*projPoint)
	P.c = P2.c
	P.X.Set(&P2.X)
	P.Y.Set(&P2.Y)
	P.Z.Set(&P2.Z)
	return P
}

func (P *projPoint) Clone() abstract.Point {
	return &projPoint{
		c: P.c,
		X: P.X,
		Y: P.Y,
		Z: P.Z,
	}
}

func (P *projPoint) Null() abstract.Point {
	P.Set(&P.c.null)
	return P
}

func (P *projPoint) Base() abstract.Point {
	P.Set(&P.c.base)
	return P
}

func (P *projPoint) PickLen() int {
	return P.c.pickLen()
}

// Normalize the point's representation to Z=1.
func (P *projPoint) normalize() {
	P.Z.Inv(&P.Z)
	P.X.Mul(&P.X, &P.Z)
	P.Y.Mul(&P.Y, &P.Z)
	P.Z.V.SetInt64(1)
}

func (P *projPoint) Pick(data []byte, rand cipher.Stream) (abstract.Point, []byte) {
	return P, P.c.pickPoint(P, data, rand)
}

// Extract embedded data from a point group element
func (P *projPoint) Data() ([]byte, error) {
	P.normalize()
	return P.c.data(&P.X, &P.Y)
}

// Add two points using optimized projective coordinate addition formulas.
// Formulas taken from:
//
//	http://eprint.iacr.org/2008/013.pdf
//	https://hyperelliptic.org/EFD/g1p/auto-twisted-projective.html
//
func (P *projPoint) Add(CP1, CP2 abstract.Point) abstract.Point {
	P1 := CP1.(*projPoint)
	P2 := CP2.(*projPoint)
	X1, Y1, Z1 := &P1.X, &P1.Y, &P1.Z
	X2, Y2, Z2 := &P2.X, &P2.Y, &P2.Z
	X3, Y3, Z3 := &P.X, &P.Y, &P.Z
	var A, B, C, D, E, F, G nist.Int

	A.Mul(Z1, Z2)
	B.Mul(&A, &A)
	C.Mul(X1, X2)
	D.Mul(Y1, Y2)
	E.Mul(&C, &D).Mul(&P.c.d, &E)
	F.Sub(&B, &E)
	G.Add(&B, &E)
	X3.Add(X1, Y1).Mul(X3, Z3.Add(X2, Y2)).Sub(X3, &C).Sub(X3, &D).
		Mul(&F, X3).Mul(&A, X3)
	Y3.Mul(&P.c.a, &C).Sub(&D, Y3).Mul(&G, Y3).Mul(&A, Y3)
	Z3.Mul(&F, &G)
	return P
}

// Subtract points so that their scalars subtract homomorphically
func (P *projPoint) Sub(CP1, CP2 abstract.Point) abstract.Point {
	P1 := CP1.(*projPoint)
	P2 := CP2.(*projPoint)
	X1, Y1, Z1 := &P1.X, &P1.Y, &P1.Z
	X2, Y2, Z2 := &P2.X, &P2.Y, &P2.Z
	X3, Y3, Z3 := &P.X, &P.Y, &P.Z
	var A, B, C, D, E, F, G nist.Int

	A.Mul(Z1, Z2)
	B.Mul(&A, &A)
	C.Mul(X1, X2)
	D.Mul(Y1, Y2)
	E.Mul(&C, &D).Mul(&P.c.d, &E)
	F.Add(&B, &E)
	G.Sub(&B, &E)
	X3.Add(X1, Y1).Mul(X3, Z3.Sub(Y2, X2)).Add(X3, &C).Sub(X3, &D).
		Mul(&F, X3).Mul(&A, X3)
	Y3.Mul(&P.c.a, &C).Add(&D, Y3).Mul(&G, Y3).Mul(&A, Y3)
	Z3.Mul(&F, &G)
	return P
}

// Find the negative of point A.
// For Edwards curves, the negative of (x,y) is (-x,y).
func (P *projPoint) Neg(CA abstract.Point) abstract.Point {
	A := CA.(*projPoint)
	P.c = A.c
	P.X.Neg(&A.X)
	P.Y.Set(&A.Y)
	P.Z.Set(&A.Z)
	return P
}

// Optimized point doubling for use in scalar multiplication.
func (P *projPoint) double() {
	var B, C, D, E, F, H, J nist.Int

	B.Add(&P.X, &P.Y).Mul(&B, &B)
	C.Mul(&P.X, &P.X)
	D.Mul(&P.Y, &P.Y)
	E.Mul(&P.c.a, &C)
	F.Add(&E, &D)
	H.Mul(&P.Z, &P.Z)
	J.Add(&H, &H).Sub(&F, &J)
	P.X.Sub(&B, &C).Sub(&P.X, &D).Mul(&P.X, &J)
	P.Y.Sub(&E, &D).Mul(&F, &P.Y)
	P.Z.Mul(&F, &J)
}

// Multiply point p by scalar s using the repeated doubling method.
func (P *projPoint) Mul(G abstract.Point, s abstract.Scalar) abstract.Point {
	v := s.(*nist.Int).V
	if G == nil {
		return P.Base().Mul(P, s)
	}
	T := P
	if G == P { // Must use temporary for in-place multiply
		T = &projPoint{}
	}
	T.Set(&P.c.null) // Initialize to identity element (0,1)
	for i := v.BitLen() - 1; i >= 0; i-- {
		T.double()
		if v.Bit(i) != 0 {
			T.Add(T, G)
		}
	}
	if T != P {
		P.Set(T)
	}
	return P
}

// ProjectiveCurve implements Twisted Edwards curves
// using projective coordinate representation (X:Y:Z),
// satisfying the identities x = X/Z, y = Y/Z.
// This representation still supports all Twisted Edwards curves
// and avoids expensive modular inversions on the critical paths.
// Uses the projective arithmetic formulas in:
// http://cr.yp.to/newelliptic/newelliptic-20070906.pdf
//
type ProjectiveCurve struct {
	curve           // generic Edwards curve functionality
	null  projPoint // Constant identity/null point (0,1)
	base  projPoint // Standard base point
}

// Create a new Point on this curve.
func (c *ProjectiveCurve) Point() abstract.Point {
	P := new(projPoint)
	P.c = c
	//P.Set(&c.null)
	return P
}

// Initialize the curve with given parameters.
func (c *ProjectiveCurve) Init(p *Param, fullGroup bool) *ProjectiveCurve {
	c.curve.init(c, p, fullGroup, &c.null, &c.base)
	return c
}
