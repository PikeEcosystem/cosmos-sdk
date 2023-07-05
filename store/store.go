package store

import (
	"github.com/PikeEcosystem/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/PikeEcosystem/cosmos-sdk/store/cache"
	"github.com/PikeEcosystem/cosmos-sdk/store/rootmulti"
	"github.com/PikeEcosystem/cosmos-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB) types.CommitMultiStore {
	return rootmulti.NewStore(db, log.NewNopLogger())
}

func NewCommitKVStoreCacheManager(cacheSize int, metricsProvider cache.MetricsProvider) types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cacheSize, metricsProvider)
}
