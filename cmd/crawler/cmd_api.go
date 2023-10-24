package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/node-crawler/pkg/api"
	"github.com/ethereum/node-crawler/pkg/apidb"
	"github.com/ethereum/node-crawler/pkg/crawlerdb"
)

var (
	apiCommand = cli.Command{
		Name:   "api",
		Usage:  "API server for the crawler",
		Action: startAPI,
		Flags: []cli.Flag{
			apiDBFlag,
			apiListenAddrFlag,
			crawlerDBFlag,
			dropNodesTimeFlag,
		},
	}
)

func startAPI(ctx *cli.Context) error {
	crawlerDBPath := ctx.String(crawlerDBFlag.Name)
	crawlerDB, err := openDB(
		crawlerDBPath,
	)
	if err != nil {
		return err
	}

	apiDBPath := ctx.String(apiDBFlag.Name)
	nodeDB, err := openDB(
		apiDBPath,
	)
	if err != nil {
		return err
	}
	if err := apidb.InitDB(nodeDB); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(3)

	// Start reading deamon
	go newNodeDeamon(&wg, crawlerDB, nodeDB)
	go dropDeamon(&wg, nodeDB, ctx.Duration(dropNodesTimeFlag.Name))

	// Start the API deamon
	apiAddress := ctx.String(apiListenAddrFlag.Name)
	apiDeamon := api.New(apiAddress, nodeDB)
	go apiDeamon.HandleRequests(&wg)
	wg.Wait()

	return nil
}

func transferNewNodes(crawlerDB, nodeDB *sql.DB) error {
	crawlerDBTx, err := crawlerDB.Begin()
	if err != nil {
		// Sometimes error occur trying to read the crawler database, but
		// they are normally recoverable, and a lot of the time, it's
		// because the database is locked by the crawler.
		return fmt.Errorf("error starting transaction to read nodes: %w", err)
	}
	defer crawlerDBTx.Rollback()

	nodes, err := crawlerdb.ReadAndDeleteUnseenNodes(crawlerDBTx)
	if err != nil {
		// Simiar to nodeDB.Begin() error
		return fmt.Errorf("error reading nodes: %w", err)
	}

	if len(nodes) > 0 {
		err := apidb.InsertCrawledNodes(nodeDB, nodes)
		if err != nil {
			// This shouldn't happen because the database is not shared in this
			// instance, so there shouldn't be lock errors, but anything can
			// happen. We will still try again.
			return fmt.Errorf("error inserting nodes: %w", err)
		}
		log.Info("Nodes inserted", "len", len(nodes))
	}

	crawlerDBTx.Commit()
	return nil
}

// newNodeDeamon reads new nodes from the crawler and puts them in the db
// Might trigger the invalidation of caches for the api in the future
func newNodeDeamon(wg *sync.WaitGroup, crawlerDB, nodeDB *sql.DB) {
	defer wg.Done()

	// This is so that we can make some kind of exponential backoff for the
	// retries.
	retryTimeout := time.Minute

	for {
		err := transferNewNodes(crawlerDB, nodeDB)
		if err != nil {
			log.Error("Failure in transferring new nodes", "err", err)
			time.Sleep(retryTimeout)
			retryTimeout *= 2
			continue
		}

		retryTimeout = time.Minute
		time.Sleep(time.Second)
	}
}

func dropDeamon(wg *sync.WaitGroup, db *sql.DB, dropTimeout time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := apidb.DropOldNodes(db, dropTimeout)
		if err != nil {
			panic(err)
		}
	}
}
