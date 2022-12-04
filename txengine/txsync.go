package txengine

import (
	"fmt"
)

// startTxSync caches transactions by addresses
func (engine *Txengine) startTxSync() {
	engine.logger.Writer().Write([]byte("startTxSync started\n"))
	for {
		select {
		case txs := <-engine.newTransactions:
			addresses := engine.parserUsecase.GetAllAddresses()
			for _, tx := range txs {
				transaction := tx
				if _, fromOk := addresses[tx.From]; fromOk {
					engine.parserUsecase.CacheTransaction(tx.From, &transaction)
					fmt.Printf("startTxSync: cached transaction <<from>> %s with hash %s\n", tx.From, transaction.Hash)
				}
				if _, toOk := addresses[tx.To]; toOk {
					engine.parserUsecase.CacheTransaction(tx.To, &transaction)
					fmt.Printf("startTxSync: cached transaction <<to>> %s with hash %s\n", tx.To, transaction.Hash)
				}
			}

			fmt.Printf("startTxSync: successfully processed %d txs\n", len(txs))
		case <-engine.stopCh:
			break
		}
	}
}
