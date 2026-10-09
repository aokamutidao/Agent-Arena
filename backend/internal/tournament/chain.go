package tournament

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthReader struct {
	client  *ethclient.Client
	cfg     Config
	manager abi.ABI
	abis    map[common.Address]abi.ABI
}

// NewEthReader validates chain and immutable deployment wiring using eth_call.
// No private key, signer, bind.TransactOpts or send-transaction method is used.
func NewEthReader(ctx context.Context, rpcURL string, cfg Config) (*EthReader, error) {
	if !common.IsHexAddress(cfg.ManagerAddress) || common.HexToAddress(cfg.ManagerAddress) == (common.Address{}) {
		return nil, fmt.Errorf("course manager address is required")
	}
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("connect read-only RPC: %w", err)
	}
	r := &EthReader{client: client, cfg: cfg, abis: map[common.Address]abi.ABI{}}
	r.cfg.ManagerAddress = address(cfg.ManagerAddress)
	b, err := os.ReadFile(filepath.Join(cfg.ABIPath, "TournamentManager.json"))
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("read manager ABI: %w", err)
	}
	r.manager, err = abi.JSON(strings.NewReader(string(b)))
	if err != nil {
		client.Close()
		return nil, err
	}
	r.abis[common.HexToAddress(cfg.ManagerAddress)] = r.manager
	chainID, err := client.ChainID(ctx)
	if err != nil || !chainID.IsUint64() || chainID.Uint64() != cfg.ChainID {
		client.Close()
		return nil, fmt.Errorf("RPC chain does not match configured chain_id")
	}
	code, err := client.CodeAt(ctx, common.HexToAddress(cfg.ManagerAddress), nil)
	if err != nil || len(code) == 0 {
		client.Close()
		return nil, fmt.Errorf("manager has no deployed code")
	}
	for method, destination := range map[string]*string{"settlementToken": &r.cfg.TokenAddress, "tournamentRegistry": &r.cfg.RegistryAddress, "prizeEscrow": &r.cfg.EscrowAddress, "resultReview": &r.cfg.ReviewAddress} {
		out, err := r.call(ctx, method)
		if err != nil || len(out) != 1 {
			client.Close()
			return nil, fmt.Errorf("verify deployment %s", method)
		}
		actual, ok := out[0].(common.Address)
		if !ok || actual == (common.Address{}) {
			client.Close()
			return nil, fmt.Errorf("invalid deployment %s", method)
		}
		if *destination != "" && actual != common.HexToAddress(*destination) {
			client.Close()
			return nil, fmt.Errorf("configured %s address does not match manager", method)
		}
		*destination = address(actual.Hex())
	}
	entries := []struct{ name, addr string }{{"TournamentRegistry", r.cfg.RegistryAddress}, {"PrizeEscrow", r.cfg.EscrowAddress}, {"ResultReview", r.cfg.ReviewAddress}}
	for _, item := range entries {
		b, err := os.ReadFile(filepath.Join(cfg.ABIPath, item.name+".json"))
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("read %s ABI: %w", item.name, err)
		}
		parsed, err := abi.JSON(strings.NewReader(string(b)))
		if err != nil {
			client.Close()
			return nil, err
		}
		r.abis[common.HexToAddress(item.addr)] = parsed
	}
	decimalABI, err := abi.JSON(strings.NewReader("[{\"name\":\"decimals\",\"type\":\"function\",\"stateMutability\":\"view\",\"inputs\":[],\"outputs\":[{\"type\":\"uint8\"}]}]"))
	if err != nil {
		client.Close()
		return nil, err
	}
	to := common.HexToAddress(r.cfg.TokenAddress)
	data, err := decimalABI.Pack("decimals")
	if err != nil {
		client.Close()
		return nil, err
	}
	raw, err := client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("read settlement token decimals")
	}
	decimals, err := decimalABI.Unpack("decimals", raw)
	if err != nil || len(decimals) != 1 {
		client.Close()
		return nil, fmt.Errorf("invalid token decimals")
	}
	r.cfg.TokenDecimals = decimals[0].(uint8)
	return r, nil
}
func (r *EthReader) Config() Config { return r.cfg }
func (r *EthReader) Close()         { r.client.Close() }
func (r *EthReader) call(ctx context.Context, method string, args ...any) ([]any, error) {
	return r.callAt(ctx, nil, method, args...)
}
func (r *EthReader) callAt(ctx context.Context, block *big.Int, method string, args ...any) ([]any, error) {
	data, err := r.manager.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	to := common.HexToAddress(r.cfg.ManagerAddress)
	result, err := r.client.CallContract(ctx, ethereum.CallMsg{To: &to, Data: data}, block)
	if err != nil {
		return nil, err
	}
	return r.manager.Unpack(method, result)
}
func (r *EthReader) Head(ctx context.Context) (uint64, error) { return r.client.BlockNumber(ctx) }
func (r *EthReader) BlockHash(ctx context.Context, block uint64) (string, error) {
	header, err := r.client.HeaderByNumber(ctx, new(big.Int).SetUint64(block))
	if err != nil {
		return "", err
	}
	return header.Hash().Hex(), nil
}
func (r *EthReader) Roles(ctx context.Context, wallet string) ([]string, error) {
	roles := []string{}
	for _, item := range []struct {
		name string
		hash common.Hash
	}{{"admin", common.Hash{}}, {"organizer", crypto.Keccak256Hash([]byte("ORGANIZER_ROLE"))}, {"auditor", crypto.Keccak256Hash([]byte("AUDITOR_ROLE"))}} {
		out, err := r.call(ctx, "hasRole", item.hash, common.HexToAddress(wallet))
		if err != nil {
			return nil, err
		}
		if len(out) != 1 {
			return nil, fmt.Errorf("invalid role response")
		}
		has, ok := out[0].(bool)
		if !ok {
			return nil, fmt.Errorf("invalid role type")
		}
		if has {
			roles = append(roles, item.name)
		}
	}
	return roles, nil
}

type chainTournament struct {
	Id                    *big.Int
	Organizer             common.Address
	Auditor               common.Address
	RegistrationDeadline  uint64
	SettlementDeadline    uint64
	ReviewPeriod          uint64
	Status                uint8
	Participant1          common.Address
	Participant2          common.Address
	Participant1AgentHash [32]byte
	Participant2AgentHash [32]byte
	ProposedWinner        common.Address
	Winner                common.Address
	ResultHash            [32]byte
	ReviewDeadline        uint64
	Title                 string
	MetadataURI           string
}
type chainEscrow struct {
	TotalFunded  *big.Int
	TotalPaid    *big.Int
	Winner       common.Address
	Cancelled    bool
	PrizeClaimed bool
}

func hashString(h [32]byte) string {
	if h == ([32]byte{}) {
		return ""
	}
	return common.BytesToHash(h[:]).Hex()
}

func (r *EthReader) Tournament(ctx context.Context, id uint64) (*Tournament, error) {
	head, err := r.Head(ctx)
	if err != nil {
		return nil, err
	}
	return r.TournamentAt(ctx, id, head)
}
func (r *EthReader) TournamentAt(ctx context.Context, id, block uint64) (*Tournament, error) {
	out, err := r.callAt(ctx, new(big.Int).SetUint64(block), "getTournament", new(big.Int).SetUint64(id))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "revert") {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if len(out) != 1 {
		return nil, fmt.Errorf("invalid tournament response")
	}
	v := abi.ConvertType(out[0], new(chainTournament)).(*chainTournament)
	esc, err := r.callAt(ctx, new(big.Int).SetUint64(block), "getEscrow", new(big.Int).SetUint64(id))
	if err != nil {
		return nil, err
	}
	if len(esc) != 1 {
		return nil, fmt.Errorf("invalid escrow response")
	}
	fund := abi.ConvertType(esc[0], new(chainEscrow)).(*chainEscrow)
	states := []string{"open", "running", "result_submitted", "disputed", "finalized", "cancelled"}
	if int(v.Status) >= len(states) {
		return nil, fmt.Errorf("unknown contract status")
	}
	t := &Tournament{ID: id, Title: v.Title, MetadataURI: v.MetadataURI, Organizer: address(v.Organizer.Hex()), Auditor: address(v.Auditor.Hex()), Token: address(r.cfg.TokenAddress), EntryFee: "0", PrizePool: fund.TotalFunded.String(), PrizePaid: fund.TotalPaid.String(), PrizeClaimed: fund.PrizeClaimed, MaxParticipants: 2, RegistrationDeadline: int64(v.RegistrationDeadline), SettlementDeadline: int64(v.SettlementDeadline), ReviewPeriod: int64(v.ReviewPeriod), ReviewDeadline: int64(v.ReviewDeadline), Status: states[v.Status], Participants: []string{}, ParticipantAgents: []ParticipantAgent{}, ProposedWinner: nonzeroAddress(v.ProposedWinner), Winner: nonzeroAddress(v.Winner), EvidenceHash: hashString(v.ResultHash), ChainUpdatedAt: time.Now().Unix()}
	for _, p := range []struct {
		wallet common.Address
		hash   [32]byte
	}{{v.Participant1, v.Participant1AgentHash}, {v.Participant2, v.Participant2AgentHash}} {
		if p.wallet != (common.Address{}) {
			a := address(p.wallet.Hex())
			t.Participants = append(t.Participants, a)
			t.ParticipantAgents = append(t.ParticipantAgents, ParticipantAgent{a, hashString(p.hash)})
		}
	}
	t.ReadBlock = block
	return t, nil
}

func eventValue(v any) any {
	switch n := v.(type) {
	case *big.Int:
		return n.String()
	case common.Address:
		return address(n.Hex())
	case [32]byte:
		return common.BytesToHash(n[:]).Hex()
	default:
		return v
	}
}
func (r *EthReader) Events(ctx context.Context, from, to uint64) ([]ChainEvent, error) {
	addresses := []common.Address{}
	for a := range r.abis {
		addresses = append(addresses, a)
	}
	logs, err := r.client.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: new(big.Int).SetUint64(from), ToBlock: new(big.Int).SetUint64(to), Addresses: addresses})
	if err != nil {
		return nil, err
	}
	events := []ChainEvent{}
	timestamps := map[uint64]int64{}
	for _, entry := range logs {
		if entry.Removed || len(entry.Topics) == 0 {
			continue
		}
		contractABI, ok := r.abis[entry.Address]
		if !ok {
			continue
		}
		definition, err := contractABI.EventByID(entry.Topics[0])
		if err != nil {
			continue
		}
		values := map[string]any{}
		if err := definition.Inputs.NonIndexed().UnpackIntoMap(values, entry.Data); err != nil {
			return nil, fmt.Errorf("decode %s: %w", definition.Name, err)
		}
		indexed := abi.Arguments{}
		for _, input := range definition.Inputs {
			if input.Indexed {
				indexed = append(indexed, input)
			}
		}
		if err := abi.ParseTopicsIntoMap(values, indexed, entry.Topics[1:]); err != nil {
			return nil, err
		}
		var tournamentID uint64
		if n, ok := values["tournamentId"].(*big.Int); ok {
			if !n.IsUint64() {
				return nil, fmt.Errorf("tournament ID exceeds backend limit")
			}
			tournamentID = n.Uint64()
		}
		actor := ""
		for _, key := range []string{"participant", "sponsor", "winner", "organizer", "auditor", "submitter", "challenger", "caller", "account"} {
			if a, ok := values[key].(common.Address); ok {
				actor = address(a.Hex())
				break
			}
		}
		amount := ""
		if entry.Address == common.HexToAddress(r.cfg.ManagerAddress) && (definition.Name == "TournamentFunded" || definition.Name == "PrizeClaimed" || definition.Name == "SponsorRefunded") {
			if n, ok := values["amount"].(*big.Int); ok {
				amount = n.String()
			} else {
				return nil, fmt.Errorf("financial event missing amount")
			}
		}
		for key, v := range values {
			values[key] = eventValue(v)
		}
		data, err := json.Marshal(values)
		if err != nil {
			return nil, err
		}
		ts, ok := timestamps[entry.BlockNumber]
		if !ok {
			header, err := r.client.HeaderByNumber(ctx, new(big.Int).SetUint64(entry.BlockNumber))
			if err != nil {
				return nil, err
			}
			ts = int64(header.Time)
			timestamps[entry.BlockNumber] = ts
		}
		events = append(events, ChainEvent{ID: fmt.Sprintf("%d:%s:%d", r.cfg.ChainID, strings.ToLower(entry.TxHash.Hex()), entry.Index), ChainID: r.cfg.ChainID, TournamentID: tournamentID, ContractAddress: address(entry.Address.Hex()), EventName: definition.Name, TxHash: entry.TxHash.Hex(), LogIndex: entry.Index, BlockNumber: entry.BlockNumber, Timestamp: ts, Actor: actor, AmountRaw: amount, Data: data})
	}
	return events, nil
}
