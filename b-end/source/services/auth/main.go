package auth

import (
	"fmt"
	"nft_marketplace/eth/source/database"
	"nft_marketplace/eth/source/database/models"
	"nft_marketplace/eth/source/services/users"
	"nft_marketplace/eth/source/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
)

type CreateAccountDTO struct {
    Signature   string  `json:"Signature"`
    Challenge     string  `json:"Challenge"`
    Username    string  `json:"Username"`
}

func VerifySignatureAndCreateNewAccount(dto *CreateAccountDTO) (models.User, error) {
    var user models.User

    result := database.Postgres.Where("username = ?", dto.Username).First(&user)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return user, fmt.Errorf("User not found") 
        } 
    }

    // decode signature
    decodedSignature, err := hexutil.Decode(dto.Signature)
    if err != nil {
        return user, err
    }

    v, err := utils.NormalizeRecoveryID(decodedSignature[64])
    if err != nil {
        return user, err
    }

    decodedSignature[64] = v

    pubKey, err := utils.ExtractPublicKeyFromSignedChallenge([]byte(dto.Challenge), decodedSignature)
    if err != nil {
        return user, err
    }

    fmt.Println("pub key : " + string(crypto.CompressPubkey(&pubKey)))

    address, err := utils.ExtractAddressFromPubKey(pubKey)
    database.Postgres.Where("Address = ?", address).First(&user)

    if user.ID == 0 {
        return user, fmt.Errorf("Invalid signature")
    }

    user.PubKey = utils.PubKeyECDSAToString(pubKey)
    user.Address = address
    users.AddUser(&user)

    return user, nil
}

func CheckIfUserExists(username string) (bool, error) {
    _, result := users.FindUserByUsername(username)

    if result.Error != nil {
        if result.Error != gorm.ErrRecordNotFound {

        } else {
            return false, result.Error
        }
    }

    return true, nil
}
