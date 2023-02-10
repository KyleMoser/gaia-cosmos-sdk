package util

import (
	b64 "encoding/base64"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func PrintTxInfo(tx sdk.Tx, caller string) {
	signingTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		fmt.Println("NOT a SigVerifiableTx")
		return
	}

	id := 0
	allTx, ok := tx.(authsigning.Tx)
	if ok {
		id, _ = strconv.Atoi(allTx.GetMemo())
	}

	pubKeys, err := signingTx.GetPubKeys()
	if err != nil {
		fmt.Printf("[SDK:%d] %s: pubkey err: %s \n", id, caller, err.Error())
		return
	}
	for _, curr := range pubKeys {
		if curr != nil {
			sEnc := b64.StdEncoding.EncodeToString(curr.Bytes())
			fmt.Printf("[SDK:%d] %s: public key b64: %s \n", id, caller, sEnc)
		} else {
			fmt.Printf("[SDK:%d] %s: public key nil \n", id, caller)
		}
	}

	sigs, err := signingTx.GetSignaturesV2()
	if err != nil {
		fmt.Printf("[SDK:%d] %s: sigs err: %s \n", id, caller, err.Error())
		return
	}
	for _, curr := range sigs {
		if curr.PubKey != nil {
			sEnc := b64.StdEncoding.EncodeToString(curr.PubKey.Bytes())
			fmt.Printf("[SDK:%d] %s: SIG public key: %s \n", id, caller, sEnc)
		} else {
			fmt.Printf("[SDK:%d] %s: SIG public key nil \n", id, caller)
		}
	}

	signers := signingTx.GetSigners()
	for _, curr := range signers {
		if curr != nil {
			fmt.Printf("[SDK:%d] %s: signer: %s \n", id, caller, curr.String())
		} else {
			fmt.Printf("[SDK:%d] %s signer nil \n", id, caller)
		}
	}
}
