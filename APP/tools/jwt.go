package tools

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenDuration  = 2 * time.Hour
	RefreshTokenDuration = 30 * 24 * time.Hour
	TokenIssure          = "zy"
)

type VoteJwt struct {
	Secret []byte
}

var Token VoteJwt

func NewToken(s string) {
	b := []byte("咕噜咕噜")
	if s != "" {
		b = []byte(s)
	}
	Token = VoteJwt{Secret: b}
}

type Claim struct {
	jwt.RegisteredClaims
	ID   int64  `json:"user_id"`
	Name string `json:"username"`
}

func (j *VoteJwt) getTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}
func (j *VoteJwt) keyFunc(token *jwt.Token) (interface{}, error) {
	return j.Secret, nil
}

// 生成token
func (j *VoteJwt) GetToken(id int64, name string) (aToken, rToken string, err error) {
	rc := jwt.RegisteredClaims{
		ExpiresAt: j.getTime(AccessTokenDuration),
		Issuer:    TokenIssure,
	}
	claim := Claim{
		ID:               id,
		Name:             name,
		RegisteredClaims: rc,
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(j.Secret)
	rc.ExpiresAt = j.getTime(RefreshTokenDuration)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString(j.Secret)
	return
}

// 验证token
func (j *VoteJwt) VerifyToken(tokenID string) (*Claim, error) {
	claim := &Claim{}
	token, err := jwt.ParseWithClaims(tokenID, claim, j.keyFunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("access token 验证失败")
	}
	return claim, nil
}

// 通过refreshtoken验证accesstoken
func (j *VoteJwt) RefreshToken(a, r string) (aToken, rToken string, err error) {
	if _, err = jwt.Parse(r, j.keyFunc); err != nil {
		return
	}
	claim := &Claim{}
	_, err = jwt.ParseWithClaims(a, claim, j.keyFunc)
	if errors.Is(err, jwt.ErrTokenExpired) {
		return j.GetToken(claim.ID, claim.Name)
	}
	return
}
