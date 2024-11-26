package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"nft_marketplace/eth/source/database/models"
	"nft_marketplace/eth/source/database/models/validators"
	"nft_marketplace/eth/source/handlers"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
)

func PubKeyECDSAToString(p ecdsa.PublicKey) (string){
    pubKeyBytes := crypto.FromECDSAPub(&p)
    return hex.EncodeToString(pubKeyBytes)
}

func ExtractAddressFromPubKey(p ecdsa.PublicKey)(string ,error){
    addr := crypto.PubkeyToAddress(p)
    if addr.String() == "" {
        return "", fmt.Errorf("Error while converting public key to address")
    }
    return addr.String(), nil
}

func ExtractPublicKeyFromSignedChallenge(message []byte, signature []byte)(ecdsa.PublicKey, error) {
    // hash the message
    messageHash := crypto.Keccak256Hash(message)
    fmt.Println(1)
    
    fmt.Println(2)

    // check if signature has correct length (65)
    if len(signature) != 65 {
        return ecdsa.PublicKey{}, fmt.Errorf("Signature is invalid")
    }

    fmt.Println(3, signature,  messageHash.Bytes())

    // recover pub key
    pubKey, err := crypto.SigToPub(messageHash.Bytes(), signature)
    if err != nil {
        return ecdsa.PublicKey{}, err
    }
    fmt.Println(4)

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

    if user.Username != "" {
        err = validators.Username(user.Username)
        if err != nil {
            return user, err
        }
    }

    if user.PubKey != "" {
        err = validators.EthereumPubKey(user.PubKey)
        if err != nil {
            return user, err
        }
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
    randomString := base64.URLEncoding.EncodeToString(random)
    challengeString := ttl + "_" + randomString

    return challengeString, nil
}

func HashValues(fst []byte) string {
    return crypto.Keccak256Hash(fst).String()
}

func NormalizeRecoveryID(v byte) (byte,error) {
    if v == 28 || v == 27 {
        return v - 27, nil
    }
    if v == 0 || v == 1 {
        return v, nil
    }
    return 0, errors.New("Unable to normalize recovery ID")
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
