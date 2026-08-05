package twofactor

import (
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestTOTPSetupAndVerification(t *testing.T) {
	setup, err := NewTOTPSetup("NextPT", "member@example.com")
	if err != nil {
		t.Fatalf("NewTOTPSetup() error = %v", err)
	}
	if setup.Secret == "" || !strings.HasPrefix(setup.QRCodeDataURL, "data:image/png;base64,") {
		t.Fatalf("NewTOTPSetup() returned incomplete setup: %#v", setup)
	}

	now := time.Now().UTC()
	code, err := totp.GenerateCodeCustom(setup.Secret, now, totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		t.Fatalf("GenerateCodeCustom() error = %v", err)
	}
	if !VerifyTOTP(setup.Secret, code, now) {
		t.Fatal("VerifyTOTP() = false, want true")
	}
}

func TestEncryptSecretRoundTrip(t *testing.T) {
	const secret = "JBSWY3DPEHPK3PXP"
	encrypted, err := EncryptSecret(secret, "test-encryption-key")
	if err != nil {
		t.Fatalf("EncryptSecret() error = %v", err)
	}
	if encrypted == secret {
		t.Fatal("EncryptSecret() returned plaintext")
	}
	decrypted, err := DecryptSecret(encrypted, "test-encryption-key")
	if err != nil {
		t.Fatalf("DecryptSecret() error = %v", err)
	}
	if decrypted != secret {
		t.Fatalf("DecryptSecret() = %q, want %q", decrypted, secret)
	}
}

func TestRecoveryCodes(t *testing.T) {
	codes, err := NewRecoveryCodes(3)
	if err != nil {
		t.Fatalf("NewRecoveryCodes() error = %v", err)
	}
	if len(codes) != 3 {
		t.Fatalf("NewRecoveryCodes() count = %d, want 3", len(codes))
	}
	if RecoveryCodeHash(codes[0]) != RecoveryCodeHash(strings.ToLower(strings.ReplaceAll(codes[0], "-", ""))) {
		t.Fatal("RecoveryCodeHash() does not normalize recovery codes")
	}
}
