package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

//types for txs history
type TxsDump struct{
	Number  int
	Type 	string
	Time 	time.Time
	Message sdk.Msg
}

//func NewTxsDump (name string, bid sdk.Coins, buyer sdk.AccAddress) TxsDump{
//	return TxsDump{
//		Name: name,
//		Bid: bid,
//		Buyer: buyer,
//	}
//}
//
//func (td TxsDump) String() string{
//	return strings.TrimSpace(fmt.Sprintf(`Name: %s
//		Bid: %s
//		Buyer: %s`, td.Name, td.Bid, td.Buyer))
//}