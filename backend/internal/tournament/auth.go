package tournament

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	store *Store
	cfg   Config
	chain ChainReader
	now   func() time.Time
}
type NonceResponse struct {
	Nonce     string `json:"nonce"`
	Message   string `json:"message"`
	ExpiresAt int64  `json:"expires_at"`
}
type LoginRequest struct {
	Address   string `json:"address"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
type sessionClaims struct {
	Address        string `json:"address"`
	ChainID        uint64 `json:"chain_id"`
	ManagerAddress string `json:"manager_address"`
	jwt.RegisteredClaims
}

func NewAuth(cfg Config, s *Store, c ChainReader) (*Auth, error) {
	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}
	if strings.ContainsAny(cfg.AuthDomain, "\r\n /") {
		return nil, fmt.Errorf("AUTH_DOMAIN must be an exact host, optionally with a port")
	}
	return &Auth{store: s, cfg: cfg, chain: c, now: time.Now}, nil
}

func (a *Auth) IssueNonce(wallet string) (*NonceResponse, error) {
	if !common.IsHexAddress(wallet) || common.HexToAddress(wallet) == (common.Address{}) {
		return nil, fmt.Errorf("invalid wallet address")
	}
	wallet = address(wallet)
	nonce, err := randomID()
	if err != nil {
		return nil, err
	}
	now := a.now()
	expires := now.Add(5 * time.Minute)
	scheme := "https"
	host := a.cfg.AuthDomain
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	ip := net.ParseIP(host)
	if host == "localhost" || (ip != nil && ip.IsLoopback()) {
		scheme = "http"
	}
	message := fmt.Sprintf("%s wants you to sign in with your Ethereum account:\n%s\n\nSign in to Agent Arena Course. This signature does not authorize a token transfer.\n\nURI: %s://%s\nVersion: 1\nChain ID: %d\nNonce: %s\nIssued At: %s\nExpiration Time: %s", a.cfg.AuthDomain, wallet, scheme, a.cfg.AuthDomain, a.cfg.ChainID, nonce, now.UTC().Format(time.RFC3339), expires.UTC().Format(time.RFC3339)) + loginResource(a.cfg)
	if err := a.store.CreateNonce(wallet, message, nonce, expires.Unix()); err != nil {
		return nil, err
	}
	return &NonceResponse{nonce, message, expires.Unix()}, nil
}

func loginResource(cfg Config) string {
	manager := "unconfigured"
	if cfg.ManagerAddress != "" {
		manager = address(cfg.ManagerAddress)
	}
	return fmt.Sprintf("\nResources:\n- urn:agent-arena:%d:%s", cfg.ChainID, manager)
}

func (a *Auth) Roles(ctx context.Context, wallet string) ([]string, error) {
	if a.chain == nil {
		return []string{"participant"}, nil
	}
	roles, err := a.chain.Roles(ctx, wallet)
	if err != nil {
		return nil, err
	}
	out := []string{"participant"}
	seen := map[string]bool{"participant": true}
	for _, r := range roles {
		if (r == "admin" || r == "organizer" || r == "auditor") && !seen[r] {
			out = append(out, r)
			seen[r] = true
		}
	}
	return out, nil
}

func (a *Auth) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if !common.IsHexAddress(req.Address) {
		return nil, fmt.Errorf("invalid wallet address")
	}
	wallet := address(req.Address)
	message, err := a.store.Nonce(wallet, req.Nonce, a.now().Unix())
	if err != nil {
		return nil, fmt.Errorf("nonce expired, consumed, or unknown")
	}
	// Persisted challenges must still belong to this exact application scope.
	// Changing a domain/deployment must not re-authorize an older signature.
	if !strings.HasPrefix(message, a.cfg.AuthDomain+" wants you to sign in with your Ethereum account:\n"+wallet+"\n") || !strings.Contains(message, fmt.Sprintf("\nChain ID: %d\n", a.cfg.ChainID)) || !strings.HasSuffix(message, loginResource(a.cfg)) {
		return nil, fmt.Errorf("login challenge belongs to a different domain, chain, or deployment")
	}
	sig := common.FromHex(req.Signature)
	if len(sig) != 65 {
		return nil, fmt.Errorf("invalid signature")
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	hash := crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))
	pub, err := crypto.SigToPub(hash, sig)
	if err != nil || address(crypto.PubkeyToAddress(*pub).Hex()) != wallet {
		return nil, fmt.Errorf("invalid signature")
	}
	roles, err := a.Roles(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("chain role verification unavailable: %w", err)
	}
	now := a.now()
	claims := sessionClaims{Address: wallet, ChainID: a.cfg.ChainID, ManagerAddress: address(a.cfg.ManagerAddress), RegisteredClaims: jwt.RegisteredClaims{Subject: wallet, Issuer: "agent-arena-course", Audience: jwt.ClaimStrings{a.cfg.AuthDomain}, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}
	if err := a.store.ConsumeNonce(wallet, req.Nonce, now.Unix()); err != nil {
		return nil, fmt.Errorf("nonce already consumed")
	}
	return &LoginResponse{Token: token, User: User{Address: wallet, Roles: roles}}, nil
}

func (a *Auth) Validate(token string) (string, error) {
	claims := &sessionClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) { return []byte(a.cfg.JWTSecret), nil }, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("agent-arena-course"), jwt.WithAudience(a.cfg.AuthDomain), jwt.WithExpirationRequired(), jwt.WithTimeFunc(a.now))
	if err != nil || !parsed.Valid || !common.IsHexAddress(claims.Address) || claims.Subject != claims.Address || claims.ChainID != a.cfg.ChainID || claims.ManagerAddress != address(a.cfg.ManagerAddress) {
		return "", fmt.Errorf("invalid or expired session")
	}
	return address(claims.Address), nil
}
