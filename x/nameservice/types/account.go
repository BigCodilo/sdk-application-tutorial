package types


//------------------НЕНУЖНЫЙ КУСОК


//
//import (
//	"fmt"
//	"github.com/tendermint/tendermint/crypto"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"strings"
//	"time"
//)
//
//type Account struct{
//	Address       sdk.AccAddress `json:"address" yaml:"address"`
//	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
//	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
//	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
//	Sequence      uint64         `json:"sequence" yaml:"sequence"`
//}
//
//func (account Account) GetAddress() sdk.AccAddress{
//	return account.Address
//}
//
//func (account Account) SetAddress(addres sdk.AccAddress) error{
//	account.Address = addres
//	return nil
//} // errors if already set.
//
//func (account Account) GetPubKey() crypto.PubKey{
//	return account.PubKey
//} // can return nil.
//
//func (account Account) SetPubKey(pubKey crypto.PubKey) error{
//	account.PubKey = pubKey
//	return nil
//}
//
//func (account Account) GetAccountNumber() uint64{
//	return account.AccountNumber
//}
//
//func (account Account) SetAccountNumber(number uint64) error{
//	account.AccountNumber = number
//	return nil
//}
//
//func (account Account) GetSequence() uint64{
//	return account.Sequence
//}
//
//func (account Account) SetSequence(sequense uint64) error{
//	account.Sequence = sequense
//	return nil
//}
//
//func (account Account) GetCoins() sdk.Coins{
//	return account.Coins
//}
//
//func (account Account) SetCoins(coins sdk.Coins) error{
//	account.Coins = coins
//	return nil
//}
//
//// Calculates the amount of coins that can be sent to other accounts given
//// the current time.
//func (account Account) SpendableCoins(blockTime time.Time) sdk.Coins{
//	return account.Coins
//}
//
//// Ensure that account implements stringer
//func (account Account) String() string{
//	return strings.TrimSpace(fmt.Sprintf(`Address: %s
//	Public key: %s
//	Account number: %s
//	Sequense: %s
//	Coins: %s`, account.Address.String(), account.PubKey, account.AccountNumber, account.Sequence, account.Coins.String()))
//}
//
