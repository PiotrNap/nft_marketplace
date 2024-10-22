package validators

import (
	"errors"
	"regexp"
)

func EthereumPubKey(pubKey string) (error){
    validPubKey := regexp.MustCompile("^0x[A-Fa-f0-9]{40}$") 
    if !validPubKey.MatchString(pubKey){
        return errors.New("Public key is incorrect")
    }

    return nil

}

func Username(username string) (error){
    validUsername := regexp.MustCompile("^[a-zA-Z0-9_]+$")
    if !validUsername.MatchString(username){
        return errors.New("Username contains bad characters")
    }

    return nil
}
