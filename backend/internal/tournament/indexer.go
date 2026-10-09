package tournament

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Indexer struct {
	store      *Store
	chain      ChainReader
	cfg        Config
	source     string
	mu         sync.RWMutex
	indexed    uint64
	head       uint64
	lastError  string
	observed   int64
	hasIndexed bool
	work       sync.Mutex
	reorg      bool
}

func NewIndexer(cfg Config, s *Store, c ChainReader) *Indexer {
	return &Indexer{store: s, chain: c, cfg: cfg, source: fmt.Sprintf("%d:%s", cfg.ChainID, address(cfg.ManagerAddress))}
}
func (i *Indexer) State() (indexed, head uint64, lastError string, observed int64, ready bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.indexed, i.head, i.lastError, i.observed, i.chain != nil && i.hasIndexed && i.lastError == ""
}
func (i *Indexer) error(err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if err == nil {
		i.lastError = ""
	} else {
		i.lastError = "Chain synchronization failed; verify deployment/RPC and retry."
		if strings.Contains(err.Error(), "checkpoint changed") {
			i.lastError = "Confirmed chain checkpoint changed; indexing paused until the operator rebuilds the projection."
		}
	}
	i.observed = time.Now().Unix()
}
func (i *Indexer) SyncOnce(ctx context.Context) error {
	i.work.Lock()
	defer i.work.Unlock()
	if i.chain == nil {
		return fmt.Errorf("blockchain deployment is not configured")
	}
	if i.reorg {
		return fmt.Errorf("chain checkpoint changed: indexing paused; operator must rebuild the confirmed projection")
	}
	if checker, ok := i.chain.(CheckpointReader); ok {
		height, expected, err := i.store.Checkpoint(i.source)
		if err != nil {
			i.error(err)
			return err
		}
		if expected != "" {
			actual, err := checker.BlockHash(ctx, height)
			if err != nil {
				i.error(err)
				return err
			}
			if actual != expected {
				i.reorg = true
				err := fmt.Errorf("confirmed chain checkpoint changed: indexing paused; rebuild confirmed projection")
				i.error(err)
				return err
			}
		}
	}
	head, err := i.chain.Head(ctx)
	if err != nil {
		i.error(err)
		return err
	}
	i.mu.Lock()
	i.head = head
	i.observed = time.Now().Unix()
	i.mu.Unlock()
	if head < i.cfg.Confirmations {
		i.error(fmt.Errorf("waiting for confirmation depth"))
		return nil
	}
	safe := head - i.cfg.Confirmations
	next, err := i.store.Cursor(i.source, i.cfg.DeploymentBlock)
	if err != nil {
		i.error(err)
		return err
	}
	for next <= safe {
		end := next + 1999
		if end > safe {
			end = safe
		}
		events, err := i.chain.Events(ctx, next, end)
		if err != nil {
			i.error(err)
			return err
		}
		ids := map[uint64]bool{}
		for _, e := range events {
			if e.TournamentID > 0 {
				ids[e.TournamentID] = true
			}
		}
		ordered := make([]uint64, 0, len(ids))
		for id := range ids {
			ordered = append(ordered, id)
		}
		sort.Slice(ordered, func(a, b int) bool { return ordered[a] < ordered[b] })
		snapshots := make([]*Tournament, 0, len(ordered))
		for _, id := range ordered {
			var t *Tournament
			var err error
			if pinned, ok := i.chain.(PinnedTournamentReader); ok {
				t, err = pinned.TournamentAt(ctx, id, end)
			} else {
				t, err = i.chain.Tournament(ctx, id)
			}
			if err != nil {
				i.error(err)
				return err
			}
			t.ReadBlock = end
			t.Confirmed = true
			snapshots = append(snapshots, t)
		}
		checkpoint := ""
		if checker, ok := i.chain.(CheckpointReader); ok {
			checkpoint, err = checker.BlockHash(ctx, end)
			if err != nil {
				i.error(err)
				return err
			}
		}
		if err := i.store.CommitBatch(i.source, end+1, events, snapshots, checkpoint); err != nil {
			i.error(err)
			return err
		}
		next = end + 1
		i.mu.Lock()
		i.indexed = end
		i.hasIndexed = true
		i.mu.Unlock()
	}
	if next > i.cfg.DeploymentBlock {
		i.mu.Lock()
		i.indexed = next - 1
		i.hasIndexed = true
		i.mu.Unlock()
	}
	i.error(nil)
	return nil
}
func (i *Indexer) Run(ctx context.Context) {
	ticker := time.NewTicker(i.cfg.PollInterval)
	defer ticker.Stop()
	for {
		runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		_ = i.SyncOnce(runCtx)
		cancel()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
