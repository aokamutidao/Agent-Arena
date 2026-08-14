// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = time.Tick
	_ = context.Background
)

// ArenaCoinBindingsMetaData contains all meta data concerning the ArenaCoinBindings contract.
var ArenaCoinBindingsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DAILY_CLAIM_AMOUNT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canClaim\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimDaily\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastClaimTime\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"timeUntilNextClaim\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DailyClaimed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f80fd5b5060043610610111575f3560e01c806389934d3a1161009e578063a9059cbb1161006e578063a9059cbb1461021d578063b77cf9c614610230578063bf3506c11461024f578063dd62ed3e14610262578063f2fde38b1461029a575f80fd5b806389934d3a146101df5780638da5cb5b146101e757806392d82dbe1461020257806395d89b4114610215575f80fd5b8063313ce567116100e4578063313ce5671461017b57806340c10f191461018a5780636c8cbd0d1461019f57806370a08231146101af578063715018a6146101d7575f80fd5b806306fdde0314610115578063095ea7b31461013357806318160ddd1461015657806323b872dd14610168575b5f80fd5b61011d6102ad565b60405161012a91906108db565b60405180910390f35b610146610141366004610942565b61033d565b604051901515815260200161012a565b6002545b60405190815260200161012a565b61014661017636600461096a565b610356565b6040516012815260200161012a565b61019d610198366004610942565b610379565b005b61015a68056bc75e2d6310000081565b61015a6101bd3660046109a3565b6001600160a01b03165f9081526020819052604090205490565b61019d61038f565b61019d6103a2565b6005546040516001600160a01b03909116815260200161012a565b61015a6102103660046109a3565b610476565b61011d6104bf565b61014661022b366004610942565b6104ce565b61015a61023e3660046109a3565b60066020525f908152604090205481565b61014661025d3660046109a3565b6104db565b61015a6102703660046109bc565b6001600160a01b039182165f90815260016020908152604080832093909416825291909152205490565b61019d6102a83660046109a3565b610509565b6060600380546102bc906109ed565b80601f01602080910402602001604051908101604052809291908181526020018280546102e8906109ed565b80156103335780601f1061030a57610100808354040283529160200191610333565b820191905f5260205f20905b81548152906001019060200180831161031657829003601f168201915b5050505050905090565b5f3361034a818585610546565b60019150505b92915050565b5f33610363858285610558565b61036e8585856105d4565b506001949350505050565b610381610631565b61038b828261065e565b5050565b610397610631565b6103a05f610692565b565b335f908152600660205260409020546103be9062015180610a39565b4210156104125760405162461bcd60e51b815260206004820181905260248201527f43616e206f6e6c7920636c61696d206f6e63652070657220323420686f75727360448201526064015b60405180910390fd5b335f8181526006602052604090204290556104369068056bc75e2d6310000061065e565b60405168056bc75e2d63100000815233907fb924ab64c28bc1c27fc5a56df34d576dfebbf347ea99c0e51abc4a46f458a77c9060200160405180910390a2565b6001600160a01b0381165f90815260066020526040812054819061049d9062015180610a39565b90508042106104ae57505f92915050565b6104b84282610a4c565b9392505050565b6060600480546102bc906109ed565b5f3361034a8185856105d4565b6001600160a01b0381165f908152600660205260408120546105009062015180610a39565b42101592915050565b610511610631565b6001600160a01b03811661053a57604051631e4fbdf760e01b81525f6004820152602401610409565b61054381610692565b50565b61055383838360016106e3565b505050565b6001600160a01b038381165f908152600160209081526040808320938616835292905220545f198110156105ce57818110156105c057604051637dc7a0d960e11b81526001600160a01b03841660048201526024810182905260448101839052606401610409565b6105ce84848484035f6106e3565b50505050565b6001600160a01b0383166105fd57604051634b637e8f60e11b81525f6004820152602401610409565b6001600160a01b0382166106265760405163ec442f0560e01b81525f6004820152602401610409565b6105538383836107b5565b6005546001600160a01b031633146103a05760405163118cdaa760e01b8152336004820152602401610409565b6001600160a01b0382166106875760405163ec442f0560e01b81525f6004820152602401610409565b61038b5f83836107b5565b600580546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b6001600160a01b03841661070c5760405163e602df0560e01b81525f6004820152602401610409565b6001600160a01b03831661073557604051634a1406b160e11b81525f6004820152602401610409565b6001600160a01b038085165f90815260016020908152604080832093871683529290522082905580156105ce57826001600160a01b0316846001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516107a791815260200190565b60405180910390a350505050565b6001600160a01b0383166107df578060025f8282546107d49190610a39565b9091555061084f9050565b6001600160a01b0383165f90815260208190526040902054818110156108315760405163391434e360e21b81526001600160a01b03851660048201526024810182905260448101839052606401610409565b6001600160a01b0384165f9081526020819052604090209082900390555b6001600160a01b03821661086b57600280548290039055610889565b6001600160a01b0382165f9081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040516108ce91815260200190565b60405180910390a3505050565b5f602080835283518060208501525f5b81811015610907578581018301518582016040015282016108eb565b505f604082860101526040601f19601f8301168501019250505092915050565b80356001600160a01b038116811461093d575f80fd5b919050565b5f8060408385031215610953575f80fd5b61095c83610927565b946020939093013593505050565b5f805f6060848603121561097c575f80fd5b61098584610927565b925061099360208501610927565b9150604084013590509250925092565b5f602082840312156109b3575f80fd5b6104b882610927565b5f80604083850312156109cd575f80fd5b6109d683610927565b91506109e460208401610927565b90509250929050565b600181811c90821680610a0157607f821691505b602082108103610a1f57634e487b7160e01b5f52602260045260245ffd5b50919050565b634e487b7160e01b5f52601160045260245ffd5b8082018082111561035057610350610a25565b8181038181111561035057610350610a2556fea2646970667358221220162101176317791d91eb97874fbf60167d35ad1349b49ad189c009b58218671064736f6c63430008180033",
}

// ArenaCoinBindingsABI is the input ABI used to generate the binding from.
// Deprecated: Use ArenaCoinBindingsMetaData.ABI instead.
var ArenaCoinBindingsABI = ArenaCoinBindingsMetaData.ABI

// ArenaCoinBindingsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ArenaCoinBindingsMetaData.Bin instead.
var ArenaCoinBindingsBin = ArenaCoinBindingsMetaData.Bin

// DeployArenaCoinBindings deploys a new Ethereum contract, binding an instance of ArenaCoinBindings to it.
func DeployArenaCoinBindings(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ArenaCoinBindings, error) {
	parsed, err := ArenaCoinBindingsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ArenaCoinBindingsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ArenaCoinBindings{ArenaCoinBindingsCaller: ArenaCoinBindingsCaller{contract: contract}, ArenaCoinBindingsTransactor: ArenaCoinBindingsTransactor{contract: contract}, ArenaCoinBindingsFilterer: ArenaCoinBindingsFilterer{contract: contract}}, nil
}

// ArenaCoinBindings is an auto generated Go binding around an Ethereum contract.
type ArenaCoinBindings struct {
	ArenaCoinBindingsCaller     // Read-only binding to the contract
	ArenaCoinBindingsTransactor // Write-only binding to the contract
	ArenaCoinBindingsFilterer   // Log filterer for contract events
}

// ArenaCoinBindingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArenaCoinBindingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArenaCoinBindingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArenaCoinBindingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArenaCoinBindingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArenaCoinBindingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArenaCoinBindingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArenaCoinBindingsSession struct {
	Contract     *ArenaCoinBindings // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ArenaCoinBindingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArenaCoinBindingsCallerSession struct {
	Contract *ArenaCoinBindingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ArenaCoinBindingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArenaCoinBindingsTransactorSession struct {
	Contract     *ArenaCoinBindingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ArenaCoinBindingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArenaCoinBindingsRaw struct {
	Contract *ArenaCoinBindings // Generic contract binding to access the raw methods on
}

// ArenaCoinBindingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArenaCoinBindingsCallerRaw struct {
	Contract *ArenaCoinBindingsCaller // Generic read-only contract binding to access the raw methods on
}

// ArenaCoinBindingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArenaCoinBindingsTransactorRaw struct {
	Contract *ArenaCoinBindingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArenaCoinBindings creates a new instance of ArenaCoinBindings, bound to a specific deployed contract.
func NewArenaCoinBindings(address common.Address, backend bind.ContractBackend) (*ArenaCoinBindings, error) {
	contract, err := bindArenaCoinBindings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindings{ArenaCoinBindingsCaller: ArenaCoinBindingsCaller{contract: contract}, ArenaCoinBindingsTransactor: ArenaCoinBindingsTransactor{contract: contract}, ArenaCoinBindingsFilterer: ArenaCoinBindingsFilterer{contract: contract}}, nil
}

// NewArenaCoinBindingsCaller creates a new read-only instance of ArenaCoinBindings, bound to a specific deployed contract.
func NewArenaCoinBindingsCaller(address common.Address, caller bind.ContractCaller) (*ArenaCoinBindingsCaller, error) {
	contract, err := bindArenaCoinBindings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsCaller{contract: contract}, nil
}

// NewArenaCoinBindingsTransactor creates a new write-only instance of ArenaCoinBindings, bound to a specific deployed contract.
func NewArenaCoinBindingsTransactor(address common.Address, transactor bind.ContractTransactor) (*ArenaCoinBindingsTransactor, error) {
	contract, err := bindArenaCoinBindings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsTransactor{contract: contract}, nil
}

// NewArenaCoinBindingsFilterer creates a new log filterer instance of ArenaCoinBindings, bound to a specific deployed contract.
func NewArenaCoinBindingsFilterer(address common.Address, filterer bind.ContractFilterer) (*ArenaCoinBindingsFilterer, error) {
	contract, err := bindArenaCoinBindings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsFilterer{contract: contract}, nil
}

// bindArenaCoinBindings binds a generic wrapper to an already deployed contract.
func bindArenaCoinBindings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ArenaCoinBindingsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArenaCoinBindings *ArenaCoinBindingsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArenaCoinBindings.Contract.ArenaCoinBindingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArenaCoinBindings *ArenaCoinBindingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.ArenaCoinBindingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArenaCoinBindings *ArenaCoinBindingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.ArenaCoinBindingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ArenaCoinBindings *ArenaCoinBindingsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ArenaCoinBindings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.contract.Transact(opts, method, params...)
}

// DAILYCLAIMAMOUNT is a free data retrieval call binding the contract method 0x6c8cbd0d.
//
// Solidity: function DAILY_CLAIM_AMOUNT() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) DAILYCLAIMAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "DAILY_CLAIM_AMOUNT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DAILYCLAIMAMOUNT is a free data retrieval call binding the contract method 0x6c8cbd0d.
//
// Solidity: function DAILY_CLAIM_AMOUNT() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) DAILYCLAIMAMOUNT() (*big.Int, error) {
	return _ArenaCoinBindings.Contract.DAILYCLAIMAMOUNT(&_ArenaCoinBindings.CallOpts)
}

// DAILYCLAIMAMOUNT is a free data retrieval call binding the contract method 0x6c8cbd0d.
//
// Solidity: function DAILY_CLAIM_AMOUNT() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) DAILYCLAIMAMOUNT() (*big.Int, error) {
	return _ArenaCoinBindings.Contract.DAILYCLAIMAMOUNT(&_ArenaCoinBindings.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.Allowance(&_ArenaCoinBindings.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.Allowance(&_ArenaCoinBindings.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.BalanceOf(&_ArenaCoinBindings.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.BalanceOf(&_ArenaCoinBindings.CallOpts, account)
}

// CanClaim is a free data retrieval call binding the contract method 0xbf3506c1.
//
// Solidity: function canClaim(address user) view returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) CanClaim(opts *bind.CallOpts, user common.Address) (bool, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "canClaim", user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanClaim is a free data retrieval call binding the contract method 0xbf3506c1.
//
// Solidity: function canClaim(address user) view returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) CanClaim(user common.Address) (bool, error) {
	return _ArenaCoinBindings.Contract.CanClaim(&_ArenaCoinBindings.CallOpts, user)
}

// CanClaim is a free data retrieval call binding the contract method 0xbf3506c1.
//
// Solidity: function canClaim(address user) view returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) CanClaim(user common.Address) (bool, error) {
	return _ArenaCoinBindings.Contract.CanClaim(&_ArenaCoinBindings.CallOpts, user)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Decimals() (uint8, error) {
	return _ArenaCoinBindings.Contract.Decimals(&_ArenaCoinBindings.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) Decimals() (uint8, error) {
	return _ArenaCoinBindings.Contract.Decimals(&_ArenaCoinBindings.CallOpts)
}

// LastClaimTime is a free data retrieval call binding the contract method 0xb77cf9c6.
//
// Solidity: function lastClaimTime(address ) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) LastClaimTime(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "lastClaimTime", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastClaimTime is a free data retrieval call binding the contract method 0xb77cf9c6.
//
// Solidity: function lastClaimTime(address ) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) LastClaimTime(arg0 common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.LastClaimTime(&_ArenaCoinBindings.CallOpts, arg0)
}

// LastClaimTime is a free data retrieval call binding the contract method 0xb77cf9c6.
//
// Solidity: function lastClaimTime(address ) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) LastClaimTime(arg0 common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.LastClaimTime(&_ArenaCoinBindings.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Name() (string, error) {
	return _ArenaCoinBindings.Contract.Name(&_ArenaCoinBindings.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) Name() (string, error) {
	return _ArenaCoinBindings.Contract.Name(&_ArenaCoinBindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Owner() (common.Address, error) {
	return _ArenaCoinBindings.Contract.Owner(&_ArenaCoinBindings.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) Owner() (common.Address, error) {
	return _ArenaCoinBindings.Contract.Owner(&_ArenaCoinBindings.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Symbol() (string, error) {
	return _ArenaCoinBindings.Contract.Symbol(&_ArenaCoinBindings.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) Symbol() (string, error) {
	return _ArenaCoinBindings.Contract.Symbol(&_ArenaCoinBindings.CallOpts)
}

// TimeUntilNextClaim is a free data retrieval call binding the contract method 0x92d82dbe.
//
// Solidity: function timeUntilNextClaim(address user) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) TimeUntilNextClaim(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "timeUntilNextClaim", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeUntilNextClaim is a free data retrieval call binding the contract method 0x92d82dbe.
//
// Solidity: function timeUntilNextClaim(address user) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) TimeUntilNextClaim(user common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.TimeUntilNextClaim(&_ArenaCoinBindings.CallOpts, user)
}

// TimeUntilNextClaim is a free data retrieval call binding the contract method 0x92d82dbe.
//
// Solidity: function timeUntilNextClaim(address user) view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) TimeUntilNextClaim(user common.Address) (*big.Int, error) {
	return _ArenaCoinBindings.Contract.TimeUntilNextClaim(&_ArenaCoinBindings.CallOpts, user)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ArenaCoinBindings.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) TotalSupply() (*big.Int, error) {
	return _ArenaCoinBindings.Contract.TotalSupply(&_ArenaCoinBindings.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ArenaCoinBindings *ArenaCoinBindingsCallerSession) TotalSupply() (*big.Int, error) {
	return _ArenaCoinBindings.Contract.TotalSupply(&_ArenaCoinBindings.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Approve(&_ArenaCoinBindings.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Approve(&_ArenaCoinBindings.TransactOpts, spender, value)
}

// ClaimDaily is a paid mutator transaction binding the contract method 0x89934d3a.
//
// Solidity: function claimDaily() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) ClaimDaily(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "claimDaily")
}

// ClaimDaily is a paid mutator transaction binding the contract method 0x89934d3a.
//
// Solidity: function claimDaily() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsSession) ClaimDaily() (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.ClaimDaily(&_ArenaCoinBindings.TransactOpts)
}

// ClaimDaily is a paid mutator transaction binding the contract method 0x89934d3a.
//
// Solidity: function claimDaily() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) ClaimDaily() (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.ClaimDaily(&_ArenaCoinBindings.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Mint(&_ArenaCoinBindings.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Mint(&_ArenaCoinBindings.TransactOpts, to, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsSession) RenounceOwnership() (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.RenounceOwnership(&_ArenaCoinBindings.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.RenounceOwnership(&_ArenaCoinBindings.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Transfer(&_ArenaCoinBindings.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.Transfer(&_ArenaCoinBindings.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.TransferFrom(&_ArenaCoinBindings.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.TransferFrom(&_ArenaCoinBindings.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ArenaCoinBindings.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.TransferOwnership(&_ArenaCoinBindings.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ArenaCoinBindings *ArenaCoinBindingsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ArenaCoinBindings.Contract.TransferOwnership(&_ArenaCoinBindings.TransactOpts, newOwner)
}

// ArenaCoinBindingsApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsApprovalIterator struct {
	Event *ArenaCoinBindingsApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArenaCoinBindingsApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArenaCoinBindingsApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArenaCoinBindingsApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArenaCoinBindingsApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArenaCoinBindingsApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArenaCoinBindingsApproval represents a Approval event raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ArenaCoinBindingsApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsApprovalIterator{contract: _ArenaCoinBindings.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ArenaCoinBindingsApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArenaCoinBindingsApproval)
				if err := _ArenaCoinBindings.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) ParseApproval(log types.Log) (*ArenaCoinBindingsApproval, error) {
	event := new(ArenaCoinBindingsApproval)
	if err := _ArenaCoinBindings.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArenaCoinBindingsDailyClaimedIterator is returned from FilterDailyClaimed and is used to iterate over the raw logs and unpacked data for DailyClaimed events raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsDailyClaimedIterator struct {
	Event *ArenaCoinBindingsDailyClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArenaCoinBindingsDailyClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArenaCoinBindingsDailyClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArenaCoinBindingsDailyClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArenaCoinBindingsDailyClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArenaCoinBindingsDailyClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArenaCoinBindingsDailyClaimed represents a DailyClaimed event raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsDailyClaimed struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDailyClaimed is a free log retrieval operation binding the contract event 0xb924ab64c28bc1c27fc5a56df34d576dfebbf347ea99c0e51abc4a46f458a77c.
//
// Solidity: event DailyClaimed(address indexed user, uint256 amount)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) FilterDailyClaimed(opts *bind.FilterOpts, user []common.Address) (*ArenaCoinBindingsDailyClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.FilterLogs(opts, "DailyClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsDailyClaimedIterator{contract: _ArenaCoinBindings.contract, event: "DailyClaimed", logs: logs, sub: sub}, nil
}

// WatchDailyClaimed is a free log subscription operation binding the contract event 0xb924ab64c28bc1c27fc5a56df34d576dfebbf347ea99c0e51abc4a46f458a77c.
//
// Solidity: event DailyClaimed(address indexed user, uint256 amount)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) WatchDailyClaimed(opts *bind.WatchOpts, sink chan<- *ArenaCoinBindingsDailyClaimed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.WatchLogs(opts, "DailyClaimed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArenaCoinBindingsDailyClaimed)
				if err := _ArenaCoinBindings.contract.UnpackLog(event, "DailyClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDailyClaimed is a log parse operation binding the contract event 0xb924ab64c28bc1c27fc5a56df34d576dfebbf347ea99c0e51abc4a46f458a77c.
//
// Solidity: event DailyClaimed(address indexed user, uint256 amount)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) ParseDailyClaimed(log types.Log) (*ArenaCoinBindingsDailyClaimed, error) {
	event := new(ArenaCoinBindingsDailyClaimed)
	if err := _ArenaCoinBindings.contract.UnpackLog(event, "DailyClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArenaCoinBindingsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsOwnershipTransferredIterator struct {
	Event *ArenaCoinBindingsOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArenaCoinBindingsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArenaCoinBindingsOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArenaCoinBindingsOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArenaCoinBindingsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArenaCoinBindingsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArenaCoinBindingsOwnershipTransferred represents a OwnershipTransferred event raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ArenaCoinBindingsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsOwnershipTransferredIterator{contract: _ArenaCoinBindings.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ArenaCoinBindingsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArenaCoinBindingsOwnershipTransferred)
				if err := _ArenaCoinBindings.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) ParseOwnershipTransferred(log types.Log) (*ArenaCoinBindingsOwnershipTransferred, error) {
	event := new(ArenaCoinBindingsOwnershipTransferred)
	if err := _ArenaCoinBindings.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ArenaCoinBindingsTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsTransferIterator struct {
	Event *ArenaCoinBindingsTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ArenaCoinBindingsTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArenaCoinBindingsTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ArenaCoinBindingsTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ArenaCoinBindingsTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArenaCoinBindingsTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArenaCoinBindingsTransfer represents a Transfer event raised by the ArenaCoinBindings contract.
type ArenaCoinBindingsTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ArenaCoinBindingsTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ArenaCoinBindingsTransferIterator{contract: _ArenaCoinBindings.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ArenaCoinBindingsTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ArenaCoinBindings.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArenaCoinBindingsTransfer)
				if err := _ArenaCoinBindings.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ArenaCoinBindings *ArenaCoinBindingsFilterer) ParseTransfer(log types.Log) (*ArenaCoinBindingsTransfer, error) {
	event := new(ArenaCoinBindingsTransfer)
	if err := _ArenaCoinBindings.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
