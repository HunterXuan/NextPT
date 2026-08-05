package twofactor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image/png"
	"io"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const secretVersion = "v1"

type Setup struct {
	Secret        string
	QRCodeDataURL string
}

func NewTOTPSetup(issuer string, accountName string) (*Setup, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      strings.TrimSpace(issuer),
		AccountName: strings.TrimSpace(accountName),
		Period:      30,
		SecretSize:  20,
		Secret:      []byte{},
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, err
	}

	qrCode, err := key.Image(240, 240)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, qrCode); err != nil {
		return nil, err
	}
	return &Setup{
		Secret:        key.Secret(),
		QRCodeDataURL: "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes()),
	}, nil
}

func VerifyTOTP(secret string, code string, now time.Time) bool {
	valid, err := totp.ValidateCustom(strings.TrimSpace(code), strings.TrimSpace(secret), now, totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && valid
}

func EncryptSecret(secret string, encryptionKey string) (string, error) {
	block, err := aes.NewCipher(cipherKey(encryptionKey))
	if err != nil {
		return "", err
	}
	sealed, err := seal(block, []byte(strings.TrimSpace(secret)))
	if err != nil {
		return "", err
	}
	return secretVersion + "." + base64.RawURLEncoding.EncodeToString(sealed), nil
}

func DecryptSecret(value string, encryptionKey string) (string, error) {
	parts := strings.SplitN(strings.TrimSpace(value), ".", 2)
	if len(parts) != 2 || parts[0] != secretVersion {
		return "", fmt.Errorf("invalid two-step secret format")
	}
	sealed, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(cipherKey(encryptionKey))
	if err != nil {
		return "", err
	}
	plain, err := open(block, sealed)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func NewRecoveryCodes(count int) ([]string, error) {
	if count < 1 {
		return nil, nil
	}

	codes := make([]string, 0, count)
	for range count {
		value := make([]byte, 10)
		if _, err := io.ReadFull(rand.Reader, value); err != nil {
			return nil, err
		}
		code := strings.TrimRight(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(value), "=")
		codes = append(codes, code[:5]+"-"+code[5:10]+"-"+code[10:])
	}
	return codes, nil
}

func RecoveryCodeHash(code string) string {
	digest := sha256.Sum256([]byte(NormalizeRecoveryCode(code)))
	return hex.EncodeToString(digest[:])
}

func NormalizeRecoveryCode(code string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(code), "-", ""))
}

func cipherKey(encryptionKey string) []byte {
	digest := sha256.Sum256([]byte(encryptionKey))
	return digest[:]
}

func seal(block cipher.Block, plain []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plain, nil), nil
}

func open(block cipher.Block, sealed []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(sealed) < gcm.NonceSize() {
		return nil, fmt.Errorf("invalid two-step secret payload")
	}
	return gcm.Open(nil, sealed[:gcm.NonceSize()], sealed[gcm.NonceSize():], nil)
}
