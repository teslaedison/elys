package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// return exitingCoins, weightBalanceBonus, slippage, swapFee, slippageCoins, nil
func (p *Pool) ExitPool(ctx sdk.Context, oracleKeeper OracleKeeper, accountedPoolKeeper AccountedPoolKeeper, exitingShares math.Int, tokenOutDenom string, params Params, takerFees math.LegacyDec, applyWeightBreakingFee bool) (exitingCoins sdk.Coins, weightBalanceBonus math.LegacyDec, slippage math.LegacyDec, swapFee math.LegacyDec, takerFeesFinal math.LegacyDec, slippageCoins sdk.Coins, err error) {
	exitingCoins, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, err = p.CalcExitPoolCoinsFromShares(ctx, oracleKeeper, accountedPoolKeeper, exitingShares, tokenOutDenom, params, takerFees, applyWeightBreakingFee)
	if err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), sdk.Coins{}, err
	}

	if err := p.processExitPool(ctx, exitingCoins, exitingShares); err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), sdk.Coins{}, err
	}

	return exitingCoins, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, nil
}

// exitPool exits the pool given exitingCoins and exitingShares.
// updates the pool's liquidity and totalShares.
func (p *Pool) processExitPool(_ sdk.Context, exitingCoins sdk.Coins, exitingShares math.Int) error {
	balances := p.GetTotalPoolLiquidity().Sub(exitingCoins...)
	if err := p.UpdatePoolAssetBalances(balances); err != nil {
		return err
	}

	totalShares := p.GetTotalShares().Amount
	p.TotalShares = sdk.NewCoin(p.TotalShares.Denom, totalShares.Sub(exitingShares))

	return nil
}
