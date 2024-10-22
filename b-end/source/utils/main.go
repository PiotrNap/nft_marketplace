package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"nft_marketplace/eth/source/database/models"
	"nft_marketplace/eth/source/database/models/validators"
	"nft_marketplace/eth/source/handlers"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractAddressFromPubKey(p ecdsa.PublicKey)(string ,error){
    addr := crypto.PubkeyToAddress(p)
    if addr.String() == "" {
        return "", fmt.Errorf("Error while converting public key to address")
    }
    return addr.String(), nil
}

func ExtractPublicKeyFromSignedMessage(message []byte, signatureHex string)(ecdsa.PublicKey, error) {
    // hash the message
    messageHash := crypto.Keccak256Hash(message)
    
    // decode signature
    decodedSignature, err := hexutil.Decode(signatureHex)
    if err != nil {
        return ecdsa.PublicKey{}, err
    }

    // check if signature has correct length (65)
    if len(decodedSignature) != 65 {
        return ecdsa.PublicKey{}, fmt.Errorf("Signature is invalid")
    }

    // recover pub key
    pubKey, err := crypto.SigToPub(decodedSignature, messageHash.Bytes())
    if err != nil {
        return ecdsa.PublicKey{}, err
    }

    return *pubKey, nil
}

func ExtractUserFromBody(r *http.Request) (models.User, error) {
    var user models.User

    body, err := io.ReadAll(r.Body)
    defer r.Body.Close()

    if err != nil {
        return user, err
    }

    err = json.Unmarshal(body, &user)
    if err != nil {
        return user, err
    }

    err = validators.Username(user.Username)
    if err != nil {
        return user, err
    }

    err = validators.EthereumPubKey(user.PubKey)
    if err != nil {
        return user, err
    }

    return user, nil
}

func GenerateChallengeString() (string , error){
    now := time.Now()
    notAfter := now.Add(time.Minute * 15)
    ttl := strconv.FormatInt(notAfter.Unix(), 10)

    random := make([]byte, 32)
    _ , err:= rand.Read(random)
    if err != nil {
        return "", fmt.Errorf("Error while generating unique string")
    }
    string := base64.URLEncoding.EncodeToString(random)

    challenge := ttl + "_" + string 

    return challenge, nil
}

func GenerateJWT(username string) (string, error) {
    var jwtSecretString = os.Getenv("JWT_SECRET")
    jwtSecretBytes := []byte(jwtSecretString)

    expirationTime := time.Now().Add(time.Hour * 24)

    claims := &handlers.Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtSecretBytes)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
