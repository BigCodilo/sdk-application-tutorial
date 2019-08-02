package nameservice

import (
	"fmt"
	"github.com/BigCodilo/sdk-application-tutorial/x/nameservice/types"
	//sendKeeper "github.com/cosmos/cosmos-sdk/x/bank/internal/keeper"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case MsgCreateUser:
			return handleMsgCreateUser(ctx, keeper, msg)
		case MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case MsgSend:
			return handleMsgSend(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.

	SaveLocalTx(ctx, keeper, msg)

	return sdk.Result{}                      // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)

	SaveLocalTx(ctx, keeper, msg)

	return sdk.Result{}
}

//--------------------HANDLER-FOR-CREATING-USER------------//

func handleMsgCreateUser(ctx sdk.Context, keeper Keeper, msg MsgCreateUser) sdk.Result {
	keeper.CreateUser(ctx, msg.PubKeyBech32)

	SaveLocalTx(ctx, keeper, msg)
	return sdk.Result{}
}

//---------------------HANDLER-FOR-SENDING-COINS-----------//

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, keeper Keeper, msg types.MsgSend) sdk.Result {
	if !keeper.coinKeeper.GetSendEnabled(ctx) {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	err := keeper.coinKeeper.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "bank"),
		),
	)

	SaveLocalTx(ctx, keeper, msg)

	return sdk.Result{Events: ctx.EventManager().Events()}
}


func SaveLocalTx(ctx sdk.Context, keeper Keeper, msg sdk.Msg){
	lastTxNumber := keeper.GetNumberLastTx(ctx)
	keeper.SetTX(ctx, string(lastTxNumber + 1), TxsDump{lastTxNumber + 1,	msg.Type() , time.Now(), msg})
}