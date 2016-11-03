package jwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
)

type claims struct {
	Sub   string `json:"sub"`
	Exp   int64  `json:"exp"`
	Admin bool   `json:"admin"`
}

func fromJWT(t string, c context.Context) (*claims, error) {
	ps := strings.Split(t, ".")
	if len(ps) != 3 {
		return nil, errors.New("invalid JWT segment count")
	}

	be := ps[0] + "." + ps[1]
	s, err := decodeSegment(ps[2])
	if err != nil {
		return nil, err
	}

	err = verify(be, s, c)
	if err != nil {
		return nil, err
	}

	pj, err := decodeSegment(ps[1])
	if err != nil {
		return nil, err
	}
	cls := &claims{}
	err = json.Unmarshal(pj, cls)
	if err != nil {
		return nil, err
	}

	if time.Now().After(time.Unix(cls.Exp, 0)) {
		return nil, errors.New("token expired")
	}

	return cls, nil
}

func toJWT(cls *claims, c context.Context) (string, error) {
	h := "{\"alg\": \"AppEngine\",\"typ\": \"JWT\"}"
	p, err := json.Marshal(cls)
	if err != nil {
		return "", err
	}

	he := encodeSegment([]byte(h))
	pe := encodeSegment([]byte(p))

	be := he + "." + pe
	s, err := sign(be, c)
	if err != nil {
		return "", err
	}
	se := encodeSegment(s)

	return be + "." + se, nil
}

func encodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

func sign(s string, c context.Context) ([]byte, error) {
	_, signature, err := appengine.SignBytes(c, []byte(s))
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// http://stackoverflow.com/a/32491810/470509
func verify(b string, s []byte, c context.Context) error {
	certs, err := appengine.PublicCertificates(c)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	hasher.Write([]byte(b))

	for _, cert := range certs {
		rsaKey, err := parseRSAPublicKeyFromPEM(cert.Data)
		if err != nil {
			return err
		}

		err = rsa.VerifyPKCS1v15(rsaKey, crypto.SHA256, hasher.Sum(nil), s)
		if err == nil {
			return nil
		}
	}
	return errors.New("could not verify signature")
}

func parseRSAPublicKeyFromPEM(bs []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bs)
	if block == nil {
		return nil, errors.New("no PEM key found")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a RSA public key")
	}

	return key, nil
}
