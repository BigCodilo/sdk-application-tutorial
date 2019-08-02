package nameservice

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	keyTHistory sdk.StoreKey // KVStore for saving transactions

	accKeeper auth.AccountKeeper //Keeper for working with accounts
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec, keyTHistory sdk.StoreKey, accKeeper auth.AccountKeeper) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
		keyTHistory:	keyTHistory,
		accKeeper: accKeeper,
	}
}

//-------------------PART FOR SAVING TRANSACTIONS--------------///

//Add transaction to KVStore (from handlrs)
func (k Keeper) SetTX(ctx sdk.Context, number string, txDump TxsDump){
	store := ctx.KVStore(k.keyTHistory)
	store.Set([]byte(number), k.cdc.MustMarshalBinaryBare(txDump))
}

//Get transaction by number
func (k Keeper) GetTX(ctx sdk.Context, number string) TxsDump{
	store := ctx.KVStore(k.keyTHistory)
	txDump := TxsDump{}
	bz := store.Get([]byte(number))
	k.cdc.MustUnmarshalBinaryBare(bz, &txDump)
	return txDump
}

//Get number last tx
func (k Keeper) GetNumberLastTx(ctx sdk.Context) int{
	store := ctx.KVStore(k.keyTHistory)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	count := 0
	for ; iterator.Valid(); iterator.Next(){
		count++
	}
	return count
}

//Get all transaction from KVStore
func (k Keeper) GetAllTxs(ctx sdk.Context) []TxsDump{
	store := ctx.KVStore(k.keyTHistory)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	txDumps := []TxsDump{}
	for ; iterator.Valid(); iterator.Next() {
		txDump := TxsDump{}
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &txDump)
		txDumps = append(txDumps, txDump)
	}
	return txDumps
}

//-------------------END OF PART FOR SAVING TRANSACTIONS----------------//
//-------------------PART FOR ADDING ACCOUNTS-----------------//

//Func that will creat account
func (k Keeper) CreateUser (ctx sdk.Context, pubKeyBech32 string){
	newAccount := types.BaseAccount{}
	pubKey, err := sdk.GetAccPubKeyBech32(pubKeyBech32)
	if err != nil{
		fmt.Println("ERROR IN CREATING ACCOUNT (can't get account pubkey from bach32): ", err)
	}
	address, err := sdk.AccAddressFromHex(pubKey.Address().String())
	if err != nil{
		fmt.Println("ERROR IN CREATING ACCOUNT (can't get account addres from bech32): ", err)
	}
	err = newAccount.SetPubKey(pubKey)
	if err != nil{
		fmt.Println("something wrong: ", err)
	}
	err = newAccount.SetAddress(address)
	if err != nil{
		fmt.Println("something wrong: ", err)
	}
	err = newAccount.SetCoins(sdk.Coins{
		sdk.Coin{"BTC", sdk.NewInt(100000)},
		sdk.Coin{"USD", sdk.NewInt(100000)},
	})
	if err != nil{
		fmt.Println("something wrong: ", err)
	}
	err = newAccount.SetSequence(0)
	if err != nil{
		fmt.Println("something wrong: ", err)
	}
	err = newAccount.SetAccountNumber(k.accKeeper.GetNextAccountNumber(ctx))
	if err != nil{
		fmt.Println("something wrong: ", err)
	}
	k.accKeeper.SetAccount(ctx, &newAccount)
}

//----------------END-OF-PART-FOR-ADDING-ACCOUNT--------------//

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

// Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	//fmt.Println(k.GetAllTxs(ctx))
	fmt.Println("we are here")
	//k.CreateUser(ctx, "cosmospub1addwnpepqw5s0r974af843allmdcy6pae8rnnzjqqsnq2qkjsccy8axgk6aewm2qkr6")
	fmt.Println(k.accKeeper.GetAllAccounts(ctx))
	//k.GetAccount(ctx)
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
