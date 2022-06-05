package auth

import "github.com/golang-jwt/jwt/v4"

type Service interface {
	// data token user_id
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

// Signature key
var SECRET_KEY = []byte("BWASTARTUP_s3cr3T_k3Y")

// instance untuk akses GenerateToken
func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	// data apa yg mau di sisipkan ke dalam token <"user_id" = userID>
	// Payload
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	// Header algorithm and token type
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Signature
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil
}
