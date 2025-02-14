package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/suite"
)

type ParamsTestSuite struct {
	suite.Suite
}

func TestParamsTestSuite(t *testing.T) {
	suite.Run(t, new(ParamsTestSuite))
}

func (suite *ParamsTestSuite) TestParamKeyTable() {
	suite.Require().IsType(paramtypes.KeyTable{}, ParamKeyTable())
}

func (suite *ParamsTestSuite) TestParamsValidate() {
	validExponentialCalculation := ExponentialCalculation{
		A:             sdk.NewDec(int64(16_304_348)),
		R:             sdk.NewDecWithPrec(35, 2),
		C:             sdk.ZeroDec(),
		BondingTarget: sdk.NewDecWithPrec(66, 2),
		MaxVariance:   sdk.ZeroDec(),
	}

	validInflationDistribution := InflationDistribution{
		StakingRewards: sdk.NewDecWithPrec(1000000, 6),
		CommunityPool:  sdk.ZeroDec(),
	}

	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{
			"default",
			DefaultParams(),
			false,
		},
		{
			"valid",
			NewParams(
				"cvnt",
				validExponentialCalculation,
				validInflationDistribution,
				true,
			),
			false,
		},
		{
			"valid param literal",
			Params{
				MintDenom:              "cvnt",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution:  validInflationDistribution,
				EnableInflation:        true,
			},
			false,
		},
		{
			"invalid - denom",
			NewParams(
				"/cvnt",
				validExponentialCalculation,
				validInflationDistribution,
				true,
			),
			true,
		},
		{
			"invalid - denom",
			Params{
				MintDenom:              "",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution:  validInflationDistribution,
				EnableInflation:        true,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative A",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(-1)),
					R:             sdk.NewDecWithPrec(5, 1),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - R greater than 1",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(5, 0),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative R",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(-5, 1),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative C",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(5, 1),
					C:             sdk.NewDec(int64(-9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - BondingTarget greater than 1",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(5, 1),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 1),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative BondingTarget",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(5, 1),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2).Neg(),
					MaxVariance:   sdk.NewDecWithPrec(20, 2),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - exponential calculation - negative max Variance",
			Params{
				MintDenom: "cvnt",
				ExponentialCalculation: ExponentialCalculation{
					A:             sdk.NewDec(int64(300_000_000)),
					R:             sdk.NewDecWithPrec(5, 1),
					C:             sdk.NewDec(int64(9_375_000)),
					BondingTarget: sdk.NewDecWithPrec(50, 2),
					MaxVariance:   sdk.NewDecWithPrec(20, 2).Neg(),
				},
				InflationDistribution: validInflationDistribution,
				EnableInflation:       true,
			},
			true,
		},
		{
			"invalid - inflation distribution - negative staking rewards",
			Params{
				MintDenom:              "cvnt",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards: sdk.OneDec().Neg(),
					CommunityPool:  sdk.NewDecWithPrec(133333, 6),
				},
				EnableInflation: true,
			},
			true,
		},
		{
			"invalid - inflation distribution - negative usage incentives",
			Params{
				MintDenom:              "cvnt",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards: sdk.NewDecWithPrec(533334, 6),
					CommunityPool:  sdk.NewDecWithPrec(133333, 6),
				},
				EnableInflation: true,
			},
			true,
		},
		{
			"invalid - inflation distribution - negative community pool rewards",
			Params{
				MintDenom:              "cvnt",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards: sdk.NewDecWithPrec(533334, 6),
					CommunityPool:  sdk.OneDec().Neg(),
				},
				EnableInflation: true,
			},
			true,
		},
		{
			"invalid - inflation distribution - total distribution ratio unequal 1",
			Params{
				MintDenom:              "cvnt",
				ExponentialCalculation: validExponentialCalculation,
				InflationDistribution: InflationDistribution{
					StakingRewards: sdk.NewDecWithPrec(533333, 6),
					CommunityPool:  sdk.NewDecWithPrec(133333, 6),
				},
				EnableInflation: true,
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			suite.Require().Error(err, tc.name)
		} else {
			suite.Require().NoError(err, tc.name)
		}
	}
}
