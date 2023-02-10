package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	AccountKeeper   AccountKeeper
	BankKeeper      types.BankKeeper
	FeegrantKeeper  FeegrantKeeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = DefaultSigVerificationGasConsumer
	}

	anteDecorators := []sdk.AnteDecorator{
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewPrintDebugInfoDecorator("ante.1"),
		NewRejectExtensionOptionsDecorator(),
		NewPrintDebugInfoDecorator("ante.2"),
		NewMempoolFeeDecorator(),
		NewPrintDebugInfoDecorator("ante.3"),
		NewValidateBasicDecorator(),
		NewPrintDebugInfoDecorator("ante.4"),
		NewTxTimeoutHeightDecorator(),
		NewPrintDebugInfoDecorator("ante.5"),
		NewValidateMemoDecorator(options.AccountKeeper),
		NewPrintDebugInfoDecorator("ante.6"),
		NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewPrintDebugInfoDecorator("ante.7"),
		NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		NewPrintDebugInfoDecorator("ante.8"),
		NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		NewPrintDebugInfoDecorator("ante.9"),
		NewValidateSigCountDecorator(options.AccountKeeper),
		NewPrintDebugInfoDecorator("ante.10"),
		NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		NewPrintDebugInfoDecorator("ante.11"),
		NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewPrintDebugInfoDecorator("ante.12"),
		NewIncrementSequenceDecorator(options.AccountKeeper),
		NewPrintDebugInfoDecorator("ante.13"),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
