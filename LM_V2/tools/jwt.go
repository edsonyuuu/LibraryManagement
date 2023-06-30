package tools

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenDuration  = 2 * time.Hour       // 令牌持续时间间隔
	RefreshTokenDuration = 30 * 24 * time.Hour //令牌刷新时间间隔
	TokenIssuer          = "vote-yu"           //令牌发行者
)

type VoteJwt struct {
	Secret []byte
}

type Claim struct {
	jwt.RegisteredClaims
	ID   int64  `json:"user_id"`
	Name string `json:"username"`
}

var Token VoteJwt

func NewToken(s string) {
	b := []byte("uuu")
	if s != "" {
		b = []byte(s)
	}
	Token = VoteJwt{Secret: b}
}

// getTime
func (j *VoteJwt) getTime(t time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(t))
}

func (j *VoteJwt) keyFunc(token *jwt.Token) (interface{}, error) {
	return j.Secret, nil
}

// GetToken 颁发token access token 和 refresh token
func (j *VoteJwt) GetToken(id int64, name string) (aToken, rToken string, err error) {
	rc := jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		ExpiresAt: j.getTime(AccessTokenDuration), //token持续时间是const定义常量中的两个小时

	}
	claim := Claim{
		RegisteredClaims: rc,
		ID:               id,
		Name:             name,
	}
	//jwt.NewWithClaims() 是 Go 中用于创建 JSON Web Tokens (JWTs) 的函数,它需要一个 jwt.Claims 接口的实现对象作为参数，用于描述要生成的 JWT 的声明信息
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(j.Secret)
	//jwt.RegisteredClaims{
	//	Issuer:    "",
	//	Subject:   "",
	//	Audience:  nil,
	//	ExpiresAt: nil,
	//	NotBefore: nil,
	//	IssuedAt:  nil,
	//	ID:        "",
	//}
	//refresh token 不需要保存任何用户信息
	rc.ExpiresAt = j.getTime(RefreshTokenDuration)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString(j.Secret)
	return
}

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

// RefreshToken 通过refresh token 刷新 access token
func (j *VoteJwt) RefreshToken(a, r string) (aToken, rToken string, err error) {
	// r 无效直接返回
	if _, err = jwt.Parse(r, j.keyFunc); err != nil {
		return
	}
	//从旧access token 中解析出claims数据
	claim := &Claim{}
	_, err = jwt.ParseWithClaims(a, claim, j.keyFunc)
	// 判断错误是不是因为access token 正常过期导致的
	if errors.Is(err, jwt.ErrTokenExpired) {
		return j.GetToken(claim.ID, claim.Name)
	}
	return
}
