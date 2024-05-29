package blsSignatures

import (
    "crypto/rand"
    "fmt"
    "math/big"
	utils "github.com/peitalin/go-zktest"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)



func Ex1() {

	utils.PrintHeading("BLS Signature Tests")

    fmt.Println("Generating test BLS keys...")
    blsKey := GenerateRandomBlsKeys()
    blsKey2 := GenerateRandomBlsKeys()
    blsKey3 := GenerateRandomBlsKeys()

    fmt.Println("blsKey.PrivKey: ", blsKey.PrivKey.String())
    fmt.Println("blsKey.PubKey: ", blsKey.PubKey.String())

	var msgBytes [32]byte
    copy(msgBytes[:], "awww hell naw") // all three BLS keys sign the same message

    // signature is a G2Point
    fmt.Println("\nSigning message using BLS key...")
    signature := blsKey.SignMessage(msgBytes)
    signature2 := blsKey2.SignMessage(msgBytes)
    signature3 := blsKey3.SignMessage(msgBytes)

    fmt.Println("Message signature1:\t", signature.String())
    // fmt.Println(">>>>Message signature2: ", signature2.String())
    // fmt.Println(">>>>Message signature3: ", signature3.String())

    aggregated_sig := signature.Add(signature2).Add(signature3)
    fmt.Println("Aggregated signature:\t", aggregated_sig.String())

    aggregated_pubkey := blsKey.PubKey.Add(blsKey2.PubKey).Add(blsKey3.PubKey)
    fmt.Println("Aggregated pubkey:\t", aggregated_pubkey.String())

    agg_result, _ := aggregated_sig.Verify(&aggregated_pubkey, msgBytes)
    fmt.Println("\naggregated_signature.Verify(aggregated_pubkey, msgBytes) result: ", agg_result)
    // e(sig_agg, g2) == e(H(m), pk_agg)
    fmt.Println("e(sig_agg, g2) == e(H(m), pk_agg)")

    fmt.Println("\nVerifying signature using msg and BLS public key...")
    result, _ := signature.Verify(&blsKey.PubKey, msgBytes)
    fmt.Println("signature.Verify(pk, msgBytes) result: ", result)

}


type PrivateKey = fr.Element

func NewPrivateKey(sk string) (*PrivateKey, error) {
    ele, err := new(fr.Element).SetString(sk)
    if err != nil {
        return nil, err
    }
    return ele, nil
}

type G1Point struct {
    bn254.G1Affine
}

type G2Point struct {
	bn254.G2Affine
}

type Signature struct {
    G2Point `json:"g2_point"`
}

type KeyPair struct {
    PrivKey PrivateKey
    PubKey G1Point
}

func NewKeyPair(sk PrivateKey) KeyPair {
    pk := MulByGeneratorG1(sk)
    return KeyPair{ sk, G1Point{ pk }}
}

func NewKeyPairFromString(sk string) KeyPair {
    ele, _ := new(fr.Element).SetString(sk)
    return NewKeyPair(*ele)
}

func GenerateRandomBlsKeys() KeyPair {
    // Max value is order of the curve
    max := new(big.Int)
    max.SetString(fr.Modulus().String(), 10)
	// Generate cryptographically strong pseudo-random between 0 - max
    n, _ := rand.Int(rand.Reader, max)
    sk := new(PrivateKey).SetBigInt(n)
    return NewKeyPair(*sk)
}

func MulByGeneratorG1(a fr.Element) bn254.G1Affine {
    _, _, g1Gen, _ := bn254.Generators()
    return *new(bn254.G1Affine).ScalarMultiplication(&g1Gen, a.BigInt(new(big.Int)))
}

func MulByGeneratorG2(a fr.Element) bn254.G2Affine {
    _, _, _, g2Gen := bn254.Generators()
	return *new(bn254.G2Affine).ScalarMultiplication(&g2Gen, a.BigInt(new(big.Int)))
}


func (k *KeyPair) SignMessage(message [32]byte) Signature {
    H := HashToCurve(message) // G2Point
    sk := k.PrivKey.BigInt(new(big.Int)) // fr.Element
    sig := new(bn254.G2Affine).ScalarMultiplication(H, sk)
    return Signature{ G2Point{ *sig } }
}

// Verify a message against a public key
func (sig *Signature) Verify(pubkey *G1Point, message [32]byte) (bool, error) {
    ok, err := VerifySig(pubkey.G1Affine, sig.G2Affine, message)
    if err != nil {
        return false, err
    }
    return ok, nil
}

// Add another G1 point to this one
func (p G1Point) Add(p2 G1Point) G1Point {
	p.G1Affine.Add(&p.G1Affine, &p2.G1Affine)
	return p
}

// Add another G2 point to this one
func (p G2Point) Add(p2 G2Point) G2Point {
	p.G2Affine.Add(&p.G2Affine, &p2.G2Affine)
	return p
}

func (sig Signature) Add(otherS Signature) Signature {
	sig.G2Point = sig.G2Point.Add(otherS.G2Point)
	return sig
}


func VerifySig(
	pubkey bn254.G1Affine,
	sig bn254.G2Affine,
	msgBytes [32]byte,
) (bool, error) {

	_, _, g1Gen, _ := bn254.Generators()
	// g2GenB := getG2Generator()

	hashMsg := *HashToCurve(msgBytes) // message hashed to a G2Point

    var negSig bn254.G2Affine
    negSig.Neg(&sig)

    P := []bn254.G1Affine{ pubkey, g1Gen }
    Q := []bn254.G2Affine{ hashMsg, negSig }

    //// Hashing Message to G1 Version:
    // sig = [sk]*H(m)
    // e(H(m), pubkey) == e(H(m), [sk]*g2) == e(H(m), g2)^sk == e([sk]*H(m), g2) == e(sig, g2)
    // e(hashMsg, pubkey) == e(sig, g2Gen)

    //// Hashing Message to G2 Version:
    // e(pubkey, H(m)) == e([sk]*g1, H(m)) == e(g1, H(m))^sk == e(g1, [sk]*H(m)) == e(g1, sig)
    // e(pubkey, hashMsg) == e(g1Gen, sig)
    ok, err := bn254.PairingCheck(P, Q)
	return ok, err
}

func newFpElement(x *big.Int) fp.Element {
	var p fp.Element
	p.SetBigInt(x)
	return p
}

func NewG1Point(x, y *big.Int) *bn254.G1Affine {
	return &bn254.G1Affine{
        X: newFpElement(x),
        Y: newFpElement(y),
    }
}

func NewG2Point(X, Y [2]*big.Int) *bn254.G2Affine {
    return &bn254.G2Affine{
        X: struct{ A0, A1 fp.Element }{
            // why is [1] and [0] back to front?
            A0: newFpElement(X[1]),
            A1: newFpElement(X[0]),
        },
        Y: struct{ A0, A1 fp.Element }{
            A0: newFpElement(Y[0]),
            A1: newFpElement(Y[1]),
        },
    }
}

func HashToCurve(digest [32]byte) *bn254.G2Affine {
    // HashToCurve implements the simple hash-and-check (also sometimes try-and-increment) algorithm
    // see https://hackmd.io/@benjaminion/bls12-381#Hash-and-check

	one := new(big.Int).SetUint64(1)
	three := new(big.Int).SetUint64(3)
	x := new(big.Int)
	x.SetBytes(digest[:])

	for {
		// ECC curve: y = x^3 + 3
		xP3 := new(big.Int).Exp(x, big.NewInt(3), fp.Modulus())
        y := new(big.Int).Add(xP3, three)
        y.Mod(y, fp.Modulus())

        // Check if there is a point on the curve with this x
        // if nil -> add one and repeat
        // else -> there's your point.
        if y.ModSqrt(y, fp.Modulus()) == nil {
            x.Add(x, one).Mod(x, fp.Modulus())
        } else {
            // hash message H(m) to G1 point
            // return NewG1Point(x, y)

            // if we want to hash to G2 instead, we multiply by the G2 cofactor to convert it into a point in G2
            xx := new(fr.Element).SetBigInt(x)
            g2Msg := MulByGeneratorG2(*xx)
            fmt.Println("Message Hashed to G2: ", g2Msg.String())
            return &g2Msg
        }
	}
}

func getG1Generator() *bn254.G1Affine {
    g1Gen := new(bn254.G1Affine)
    _, err := g1Gen.X.SetString("1")
    if err != nil {
        return nil
    }
    _, err = g1Gen.Y.SetString("2")
    if err != nil {
        return nil
    }
    return g1Gen
}

func getG2Generator() *bn254.G2Affine {

	g2Gen := new(bn254.G2Affine)

	g2Gen.X.SetString(
		"10857046999023057135944570762232829481370756359578518086990519993285655852781",
		"11559732032986387107991004021392285783925812861821192530917403151452391805634",
	)
	g2Gen.Y.SetString(
		"8495653923123431417604973247489272438418190587263600148770280649306958101930",
		"4082367875863433681332203403145435568316851327593401208105741076214120093531",
	)

	return g2Gen
}