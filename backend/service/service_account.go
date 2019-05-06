package service

import (
	"github.com/irisnet/explorer/backend/conf"
	"github.com/irisnet/explorer/backend/lcd"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/orm/document"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
)

type AccountService struct {
	BaseService
}

func (service *AccountService) GetModule() Module {
	return Account
}

func (service *AccountService) Query(address string) (result model.AccountVo) {
	prefix, _, _ := utils.DecodeAndConvert(address)
	if prefix == conf.Get().Hub.Prefix.ValAddr {
		valinfo := delegatorService.QueryDelegation(address)
		result.Amount = document.Coins{valinfo.selfBond}
		result.Deposits = valinfo.delegated
	} else {
		res, err := lcd.Account(address)
		if err == nil {
			var amount document.Coins
			for _, coinStr := range res.Coins {
				coin := utils.ParseCoin(coinStr)
				amount = append(amount, coin)
			}
			result.Amount = amount
		}
		result.Deposits = delegatorService.GetDeposits(address)
	}

	result.WithdrawAddress = lcd.QueryWithdrawAddr(address)
	result.IsProfiler = isProfiler(address)
	result.Address = address
	return result
}

func (service *AccountService) QueryRichList() interface{} {
	var result []document.Account
	if err := queryAll(document.CollectionNmAccount, nil, nil, desc(document.AccountFieldBalance), 100, &result); err != nil {
		panic(types.CodeNotFound)
	}
	var accList []model.AccountInfo
	var totalAmt = float64(0)

	for _, acc := range result {
		totalAmt += acc.Balance
	}

	for _, acc := range result {
		rate, _ := utils.RoundFloat(acc.Balance/totalAmt, 6)
		accList = append(accList, model.AccountInfo{
			Address: acc.Address,
			Balance: document.Coins{
				{
					Amount: acc.Balance,
					Denom:  "iris-atto",
				},
			},
			Percent: rate,
		})
	}
	return accList
}

func isProfiler(address string) bool {
	genesis := commonService.GetGenesis()
	for _, profiler := range genesis.Result.Genesis.AppState.Guardian.Profilers {
		if profiler.Address == address {
			return true
		}
	}
	return false
}
