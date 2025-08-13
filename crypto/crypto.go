package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/scrypt"
)

// GenerateRandomBytes generates cryptographically secure random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

// GenerateRandomString generates a cryptographically secure random string
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random string: %w", err)
		}
		bytes[i] = charset[num.Int64()]
	}

	return string(bytes), nil
}

// GenerateRandomHex generates a random hex string of specified length
func GenerateRandomHex(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length / 2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomBase64 generates a random base64 string
func GenerateRandomBase64(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateToken generates a secure random token for authentication
func GenerateToken(length int) (string, error) {
	if length < 16 {
		length = 32 // Minimum secure length
	}
	return GenerateRandomBase64(length)
}

// GenerateAPIKey generates a secure API key
func GenerateAPIKey() (string, error) {
	return GenerateRandomBase64(32)
}

// GenerateSecretKey generates a secure secret key for encryption
func GenerateSecretKey() ([]byte, error) {
	return GenerateRandomBytes(32) // 256-bit key
}

// DeriveKey derives a key from a password using scrypt
func DeriveKey(password, salt []byte, keyLen int) ([]byte, error) {
	// scrypt parameters: N=32768, r=8, p=1
	return scrypt.Key(password, salt, 32768, 8, 1, keyLen)
}

// GenerateSalt generates a random salt for password hashing
func GenerateSalt() ([]byte, error) {
	return GenerateRandomBytes(16)
}

// HashPasswordWithSalt hashes a password with a given salt using scrypt
func HashPasswordWithSalt(password, salt []byte) ([]byte, error) {
	return DeriveKey(password, salt, 32)
}

// HashPasswordSecure hashes a password with a random salt
func HashPasswordSecure(password string) (hash, salt string, err error) {
	saltBytes, err := GenerateSalt()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hashBytes, err := HashPasswordWithSalt([]byte(password), saltBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	return base64.StdEncoding.EncodeToString(hashBytes),
		base64.StdEncoding.EncodeToString(saltBytes),
		nil
}

// VerifyPasswordSecure verifies a password against a hash and salt
func VerifyPasswordSecure(password, hash, salt string) (bool, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	saltBytes, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	computedHash, err := HashPasswordWithSalt([]byte(password), saltBytes)
	if err != nil {
		return false, fmt.Errorf("failed to compute hash: %w", err)
	}

	// Compare hashes
	if len(hashBytes) != len(computedHash) {
		return false, nil
	}

	for i := range hashBytes {
		if hashBytes[i] != computedHash[i] {
			return false, nil
		}
	}

	return true, nil
}

// AESEncrypt encrypts data using AES-GCM
func AESEncrypt(data, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes (256 bits)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// AESDecrypt decrypts data using AES-GCM
func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("key must be 32 bytes (256 bits)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string using AES-GCM and returns base64 encoded result
func EncryptString(plaintext string, key []byte) (string, error) {
	encrypted, err := AESEncrypt([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptString decrypts a base64 encoded string using AES-GCM
func DecryptString(ciphertext string, key []byte) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	decrypted, err := AESDecrypt(encrypted, key)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// GenerateRSAKeyPair generates an RSA key pair
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bits < 2048 {
		bits = 2048 // Minimum secure key size
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// RSAPrivateKeyToPEM converts an RSA private key to PEM format
func RSAPrivateKeyToPEM(key *rsa.PrivateKey) ([]byte, error) {
	privateKeyDER := x509.MarshalPKCS1PrivateKey(key)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}

	return pem.EncodeToMemory(privateKeyBlock), nil
}

// RSAPublicKeyToPEM converts an RSA public key to PEM format
func RSAPublicKeyToPEM(key *rsa.PublicKey) ([]byte, error) {
	publicKeyDER, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}

	return pem.EncodeToMemory(publicKeyBlock), nil
}

// RSAPrivateKeyFromPEM loads an RSA private key from PEM format
func RSAPrivateKeyFromPEM(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

// RSAPublicKeyFromPEM loads an RSA public key from PEM format
func RSAPublicKeyFromPEM(pemData []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// RSAEncrypt encrypts data using RSA public key
func RSAEncrypt(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
}

// RSADecrypt decrypts data using RSA private key
func RSADecrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
}

// RSASign signs data using RSA private key
func RSASign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hash[:])
}

// RSAVerify verifies a signature using RSA public key
func RSAVerify(data, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, 0, hash[:], signature)
}

// HMACSHA256 computes HMAC-SHA256
func HMACSHA256(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyHMACSHA256 verifies an HMAC-SHA256 signature
func VerifyHMACSHA256(data, signature, key []byte) bool {
	expectedSignature := HMACSHA256(data, key)
	return hmac.Equal(signature, expectedSignature)
}

// SimpleJWT represents a simple JWT implementation
type SimpleJWT struct {
	SecretKey []byte
}

// NewSimpleJWT creates a new SimpleJWT instance
func NewSimpleJWT(secretKey []byte) *SimpleJWT {
	return &SimpleJWT{SecretKey: secretKey}
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	Issuer    string                 `json:"iss,omitempty"`
	Subject   string                 `json:"sub,omitempty"`
	Audience  string                 `json:"aud,omitempty"`
	ExpiresAt int64                  `json:"exp,omitempty"`
	NotBefore int64                  `json:"nbf,omitempty"`
	IssuedAt  int64                  `json:"iat,omitempty"`
	ID        string                 `json:"jti,omitempty"`
	Custom    map[string]interface{} `json:"-"`
}

// CreateToken creates a JWT token with the given claims
func (j *SimpleJWT) CreateToken(claims JWTClaims) (string, error) {
	// Create header
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": "HS256",
	}

	// Set default times if not provided
	now := time.Now().Unix()
	if claims.IssuedAt == 0 {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = now + 3600 // 1 hour default
	}

	// Convert claims to map
	claimsMap := map[string]interface{}{
		"iss": claims.Issuer,
		"sub": claims.Subject,
		"aud": claims.Audience,
		"exp": claims.ExpiresAt,
		"nbf": claims.NotBefore,
		"iat": claims.IssuedAt,
		"jti": claims.ID,
	}

	// Add custom claims
	for key, value := range claims.Custom {
		claimsMap[key] = value
	}

	// Encode header and claims
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	claimsJSON, err := json.Marshal(claimsMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature
	signingString := headerEncoded + "." + claimsEncoded
	signature := HMACSHA256([]byte(signingString), j.SecretKey)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	return signingString + "." + signatureEncoded, nil
}

// VerifyToken verifies and parses a JWT token
func (j *SimpleJWT) VerifyToken(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerEncoded, claimsEncoded, signatureEncoded := parts[0], parts[1], parts[2]

	// Verify signature
	signingString := headerEncoded + "." + claimsEncoded
	expectedSignature := HMACSHA256([]byte(signingString), j.SecretKey)
	expectedSignatureEncoded := base64.RawURLEncoding.EncodeToString(expectedSignature)

	if signatureEncoded != expectedSignatureEncoded {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claimsMap map[string]interface{}
	if err := json.Unmarshal(claimsJSON, &claimsMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Extract standard claims
	claims := &JWTClaims{
		Custom: make(map[string]interface{}),
	}

	if iss, ok := claimsMap["iss"].(string); ok {
		claims.Issuer = iss
		delete(claimsMap, "iss")
	}
	if sub, ok := claimsMap["sub"].(string); ok {
		claims.Subject = sub
		delete(claimsMap, "sub")
	}
	if aud, ok := claimsMap["aud"].(string); ok {
		claims.Audience = aud
		delete(claimsMap, "aud")
	}
	if exp, ok := claimsMap["exp"].(float64); ok {
		claims.ExpiresAt = int64(exp)
		delete(claimsMap, "exp")
	}
	if nbf, ok := claimsMap["nbf"].(float64); ok {
		claims.NotBefore = int64(nbf)
		delete(claimsMap, "nbf")
	}
	if iat, ok := claimsMap["iat"].(float64); ok {
		claims.IssuedAt = int64(iat)
		delete(claimsMap, "iat")
	}
	if jti, ok := claimsMap["jti"].(string); ok {
		claims.ID = jti
		delete(claimsMap, "jti")
	}

	// Add remaining claims as custom
	for key, value := range claimsMap {
		claims.Custom[key] = value
	}

	// Verify time claims
	now := time.Now().Unix()
	if claims.ExpiresAt != 0 && now > claims.ExpiresAt {
		return nil, fmt.Errorf("token has expired")
	}
	if claims.NotBefore != 0 && now < claims.NotBefore {
		return nil, fmt.Errorf("token not yet valid")
	}

	return claims, nil
}
