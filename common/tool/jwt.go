package tool

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"time"
)

// MyCustomClaims
type MyCustomClaims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

// pkcs1 PKCS #1（Public-Key Cryptography Standards #1）是由RSA实验室定义的一组公钥加密标准，用于描述RSA加密和签名算法的实现细节。PKCS #1标准包括了关于密钥格式、填充方案以及数字签名和加密流程的规范。Go语言中，常用的库如crypto/rsa和crypto/x509提供了对PKCS #1标准的支持。
func parsePriKeyBytes(buf []byte) (*rsa.PrivateKey, error) {
	// 解码 PEM 块
	p, _ := pem.Decode(buf)
	if p == nil {
		return nil, errors.New("parse key error: failed to decode PEM block")
	}

	// 尝试解析 PKCS#1 私钥
	if key, err := x509.ParsePKCS1PrivateKey(p.Bytes); err == nil {
		return key, nil
	}

	// 尝试解析 PKCS#8 私钥
	if key, err := x509.ParsePKCS8PrivateKey(p.Bytes); err == nil {
		switch k := key.(type) {
		case *rsa.PrivateKey:
			return k, nil
		default:
			return nil, errors.New("parse key error: not an RSA private key")
		}
	}

	return nil, errors.New("parse key error: failed to parse private key")
}

// GenerateTokenUsingRS256 生成token
func GenerateTokenUsingRS256(userID, iat, seconds int64) (string, error) {
	claims := MyCustomClaims{
		UserID: int(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			//sub(Subject)	主题，标识 JWT 的主题，通常指用户的唯一标识
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP", "WEB_APP"}, //aud(Audience)	观众，标识 JWT 的接收者
			ExpiresAt: jwt.NewNumericDate(time.Unix(iat+seconds, 0)),         //exp(Expiration Time)	过期时间。标识 JWT 的过期时间，这个时间必须是将来的
			IssuedAt:  jwt.NewNumericDate(time.Unix(iat, 0)),                 //iat(Issued At)	发行时间，标识 JWT 的发行时间
			ID:        GenerateSalt(10),                                      //jti(JWT ID)	JWT 的唯一标识符，用于防止 JWT 被重放（即重复使用）
		},
	}

	//
	rsaPriKey, err := parsePriKeyBytes([]byte(PRI_KEY))
	if err != nil {
		return "", err
	}

	//	生成token
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(rsaPriKey)
	if err != nil {
		return "", err
	}

	return token, nil

}

// parsePubKeyBytes 解析 PEM 编码的 RSA 公钥
func ParsePubKeyBytes(pub_key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub_key)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the key")
	}

	var pub interface{}
	var err error

	// 尝试解析 PKIX 格式的公钥
	pub, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
		return nil, errors.New("key is not of type RSA")
	}

	// 尝试解析 PKCS1 格式的公钥
	pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		if rsaPub, ok := pub.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
		return nil, errors.New("key is not of type RSA")
	}

	return nil, errors.New("failed to parse RSA public key")
}

// ParseTokenRs256 解析token
func ParseTokenRs256(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		pub, err := ParsePubKeyBytes([]byte(PUB_KEY))
		if err != nil {
			//global.Log.Errorf("parsePubKeyBytes error: %v", err)
			return nil, err
		}
		return pub, nil
	})
	if err != nil {
		return nil, err
	}
	//	判断是否失效
	if !token.Valid {
		return nil, errors.New("claim invalid")
	}
	//	断言判断是否成功
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}
	return claims, nil
}

// GenerateSalt 为密码生成盐值
func GenerateSalt(length int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyz_0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	rand.Seed(int64(time.Now().UnixNano()))
	for i := range bytes {
		bytes[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(bytes)
}

const (
	PRI_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC0cAnvlPi+Q0AO
/3byxg9eECnK/mj0g+8JBUfMHZiF+Aq3qNyLyzRYhnajpS45O4jqlvHFlRoX5buN
km8+zY0nN1ZgWpw+arxlL2QNk5CdkqeeNRONMKUy5lSoknaogkdeotYNIsptBtUp
hQs6SFS2d1ZrdWeE5Ts7RLQWxUhgliKwZLfaEFGoHePbrJYX/V2FucvF/Y7fbKjB
BL2+7C2RH3ToJPaaMQfSoTkQzfCef3YV5mBU0iEl4szA5hiPfV1/a4Rctp8Bkhcb
b/IhEgFikUhMrxFZF+p0aTdTwAhFDaXZIZDXedrwCw+C3ohUu0qnuQHvnGZRfgDp
0gElAd5RAgMBAAECggEAPcQEgrC6HVcj/H5Sl3RZGlfqFoTUJK9tCed21lZjOajY
8lmpdWkP+CW/QvFuS0Un7zXQeVQ65GdNtn6j+hme8getV7pstakP6Is6crNK08W5
/xwoHzLBuhTCchoA6FoAWFLhdwmrxXqLSXUzjLXp2bQsLwi3cPSPPBCenRTXW8g1
yngdrJTUaFZNSrNu/Lyjhn4lvzlZul81DxMkY2EfhoSA0YuOiCR4n3TUvEoHn66y
CO2LRcHxotlVvoNbw84DcuClmuKt1gpT9nrgsKLTSafYRfHF+CcQeh6l8p1+KPya
1GKpg2ttXW+EscLV6rUh/LnJYC0qFKlVA3oL/ciAgQKBgQDeI+LnLHuceSK/Zo+j
LwQ9QSocoO7ptF28n3i2sbWydjItfPqtXOntkuue9o/ZjHcQ/JNO86USthO3Ay67
KCDKW+FXOuwlTV199Sy15zk8ioaOWZaoLQygoA58kdpI/I8Slb8PEgCrJSxrmquA
tH083i57k4mhKtwCuHsoIBUq6QKBgQDP8OGrA1G8OOxtFw2jdoXHM5XGaHZc0vPo
M819FWljNtDN17RoMKsURSuRvJfURbME6AheGA/IM+82VVYdcRv+36ofTpC/wsOs
JsFQ9E/nO2AVVVKno1ksphpqOadLxeU1yzURh4W66JjqgVf7rBNGzKudxqXzEl+D
w3uBuaCnKQKBgQCAUViD0zVASNUinPsB92nKfHb3/Jqlk1PGXpQbbIIZqZ8ImbYw
KIjUfFbxB1pG/5XT0SLCq4lCSr0OrZ7z65Utb2+2tMmuLod/9/0wwnVUnGxnlCar
1QIDUxGrMZFXMdTvlmK8MNkEA8AqFDlXamshmvJc3ffVim12gNxbbFTt6QKBgQCE
NE8VzkdyFvLiLM0EB3/odXidK59NRuXB1OWpyCo35Qr+RE00DPVILu4Te0dAs4us
6+UeBchK7hIBhmH42AgHlKZxvx6yfJ6xXfZ8hMgkaJCfH58sa+NvSq/yp3Mg7tHa
0LaNzY8NlYJbXh7VKMMcuVXHOxwZHa5SdL+aa62jeQKBgQCqlfG2+Iofn/P53Aiy
y+weau94VvO/s0YmAvrD2fGaXIff5KIG9dnIWGadpPgZz9iFr4w1b5jDCVOrSkd0
n6e9xtjukrRZvnCH9IVEDk6NF7vmAXLuQb1fDlzww+qu4y2fT8jkPaa+05KDmEar
ERiXupfjUS35nvBAGYLmG3LVGA==
-----END PRIVATE KEY-----`
	PUB_KEY = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtHAJ75T4vkNADv928sYP
XhApyv5o9IPvCQVHzB2YhfgKt6jci8s0WIZ2o6UuOTuI6pbxxZUaF+W7jZJvPs2N
JzdWYFqcPmq8ZS9kDZOQnZKnnjUTjTClMuZUqJJ2qIJHXqLWDSLKbQbVKYULOkhU
tndWa3VnhOU7O0S0FsVIYJYisGS32hBRqB3j26yWF/1dhbnLxf2O32yowQS9vuwt
kR906CT2mjEH0qE5EM3wnn92FeZgVNIhJeLMwOYYj31df2uEXLafAZIXG2/yIRIB
YpFITK8RWRfqdGk3U8AIRQ2l2SGQ13na8AsPgt6IVLtKp7kB75xmUX4A6dIBJQHe
UQIDAQAB
-----END PUBLIC KEY-----`
)
