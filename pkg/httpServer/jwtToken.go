package httpServer

import (
	"cine-circle-api/pkg/utils/envUtils"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	ExpirationDate time.Time `json:"expirationDate"`
	TokenString    string    `json:"tokenString"`
}

// GenerateTokenWithUserInfo create new token based on user info.
// UserInfo will be added into token claims.
// Environment variables can be used :
//  - APPLICATION_NAME (mandatory) : name of the application that will sign token
//  - JWT_RSA256_PRIVATE_KEY (mandatory) : secret key used to sign tokens (can be base64 encoded to be used in Goland configurations)
//  - TOKEN_EXPIRATION_HOURS (optional) : define duration before token expiration (default: 24 --> 1day)
func GenerateTokenWithUserInfo(userInfo interface{}) (token Token, err error) {
	// Get expiration duration from env or default
	expirationHoursStr := envUtils.GetFromEnvOrDefault(envTokenExpirationDurationHours, defaultTokenExpirationDurationHours)
	// Calculate expirationDate from expirationDuration or default value if not specified
	expirationHours, err := strconv.Atoi(expirationHoursStr)
	if err != nil {
		return token, errors.WithStack(err)
	}
	expirationDate := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Transform userInfo into []bytes
	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil {
		return token, errors.WithStack(err)
	}

	// Creating token with claims
	applicationName, err := envUtils.GetFromEnvOrError(envApplicationName)
	if err != nil {
		return token, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":          applicationName,
		userInfoClaims: string(userInfoBytes),
		"aud":          "any",
		"exp":          expirationDate.Unix(),
	})

	// Sign token using jwt rsa256 private key
	privateKey, err := GetRsaPrivateKey()
	if err != nil {
		return
	}

	tokenString, err := jwtToken.SignedString(privateKey)
	if err != nil {
		return token, errors.WithStack(err)
	}
	token = Token{
		ExpirationDate: expirationDate,
		TokenString:    tokenString,
	}
	return
}

// ValidateTokenAndGetUserInfo will find bearer token from Authorization header.
// Then token validity will be checked using RSA256 public key.
// Then user will be retrieved from token claims and unmarshalled into user interface.
//  - JWT_RSA256_PUBLIC_KEY (mandatory) : public key used to validate/check tokens (can be base64 encoded to be used in Goland configurations)
func ValidateTokenAndGetUserInfo(req *restful.Request, user interface{}) (err error) {
	// Get token from header
	token, err := GetTokenFromAuthorizationHeader(req)
	if err != nil {
		return
	}

	// Check and get claims from token
	claims, err := ValidateToken(token)
	if err != nil {
		return
	}

	// Get user from claims
	err = GetUserInfoFromTokenClaims(claims, user)
	if err != nil {
		return
	}
	return
}

func GetUserInfoFromTokenClaims(claims jwt.MapClaims, user interface{}) (err error) {
	utilisateurStr, ok := claims[userInfoClaims].(string)
	if !ok {
		return fmt.Errorf("cannot get utilisateurBytes from token with claims = %v", claims)
	}
	err = json.Unmarshal([]byte(utilisateurStr), user)
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

// ValidateToken checks validity of token and return claims.
//  - JWT_RSA256_PUBLIC_KEY (mandatory) : public key used to validate/check tokens (can be base64 encoded to be used in Goland configurations)
func ValidateToken(tokenStr string) (claims jwt.MapClaims, err error) {
	// Get RSA public key
	publicKey, err := GetRsaPublicKey()
	if err != nil {
		return
	}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return claims, errors.WithStack(err)
	}
	if token == nil || !token.Valid {
		return claims, fmt.Errorf("token is not valid")
	}
	return
}

// GetTokenFromAuthorizationHeader return token string from Authorization header request
func GetTokenFromAuthorizationHeader(req *restful.Request) (tokenStr string, err error) {
	// Get Authorization header that contains bearer token string
	tokenHeader := req.HeaderParameter(tokenHeaderName)
	if tokenHeader == "" {
		return tokenStr, fmt.Errorf("cannot find header %s", tokenHeaderName)
	}

	// Cut header value to get token kind (part 1) and token string (part 2)
	res := strings.Split(tokenHeader, tokenDelimiter)
	if len(res) != tokenPart {
		return tokenStr, fmt.Errorf("cannot find token from header %s with value %s", tokenHeaderName, tokenHeader)
	}
	if res[0] != tokenKind {
		return tokenStr, fmt.Errorf("token %s is not typed as %s", tokenHeader, tokenKind)
	}

	// Return token string (string after Bearer delimitation)
	tokenStr = res[1]
	return
}

// GetRsaPublicKey will return public key based on value of environment variable JWT_RSA256_PUBLIC_KEY.
// This value can be base64 encoded (needed it to be used with Goland configurations)
func GetRsaPublicKey() (publicKey *rsa.PublicKey, err error) {
	// Get value from JWT_RSA256_PUBLIC_KEY
	publicKeyFromEnv, err := envUtils.GetFromEnvOrError(envJwtRsa256PublicKey)
	if err != nil {
		return
	}

	var publicKeyStr string
	// Try to decode it using base64, if not working, it means that is not encoded
	publicKeyDecodedBytes, err := base64.StdEncoding.DecodeString(publicKeyFromEnv)
	// If err is not nil, it means that value from env variable is not base64 encoded. In this case, with use value like that without any operation
	if err == nil {
		publicKeyStr = string(publicKeyDecodedBytes)
	} else {
		publicKeyStr = publicKeyFromEnv
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyStr))
	return publicKey, errors.WithStack(err)
}

// GetRsaPrivateKey will return private key based on value of environment variable JWT_RSA256_PRIVATE_KEY.
// This value can be base64 encoded (needed it to be used with Goland configurations)
func GetRsaPrivateKey() (privateKey *rsa.PrivateKey, err error) {
	// Get value from JWT_RSA256_PRIVATE_KEY
	privateKeyFromEnv, err := envUtils.GetFromEnvOrError(envJwtRsa256PrivateKey)
	if err != nil {
		return
	}

	var privateKeyStr string
	// Try to decode it using base64, if not working, it means that is not encoded
	privateKeyDecodedBytes, err := base64.StdEncoding.DecodeString(privateKeyFromEnv)
	// If err is not nil, it means that value from env variable is not base64 encoded. In this case, with use value like that without any operation
	if err == nil {
		privateKeyStr = string(privateKeyDecodedBytes)
	} else {
		privateKeyStr = privateKeyFromEnv
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyStr))
	return privateKey, errors.WithStack(err)
}
