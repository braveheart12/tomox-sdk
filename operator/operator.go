package operator

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/tomochain/tomoxsdk/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/rabbitmq"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/utils"
)

var logger = utils.Logger

// Operator manages the transaction queue that will eventually be
// sent to the exchange contract. The Operator Wallet must be equal to the
// account that initially deployed the exchange contract or an address with operator rights
// on the contract
type Operator struct {
	// AccountService     interfaces.AccountService
	WalletService     interfaces.WalletService
	TradeService      interfaces.TradeService
	OrderService      interfaces.OrderService
	TokenService      interfaces.TokenService
	EthereumProvider  interfaces.EthereumProvider
	Exchange          interfaces.Exchange
	TxQueues          []*TxQueue
	QueueAddressIndex map[common.Address]*TxQueue
	Broker            *rabbitmq.Connection
	mutex             *sync.Mutex
}

type OperatorInterface interface {
	SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error
	QueueTrade(o *types.Order, t *types.Trade) error
	GetShortestQueue() (*TxQueue, int, error)
}

// NewOperator creates a new operator struct. It creates an exchange contract instance from the
// provided address. The error and trade events are received in the ErrorChannel and TradeChannel.
// Upon receiving errors and trades in their respective channels, event payloads are sent to the
// associated order maker and taker sockets through the through the event channel on the Order and Trade struct.
// In addition, an error event cancels the trade in the trading engine and makes the order available again.
func NewOperator(
	walletService interfaces.WalletService,
	tradeService interfaces.TradeService,
	orderService interfaces.OrderService,
	provider interfaces.EthereumProvider,
	exchange interfaces.Exchange,
	conn *rabbitmq.Connection,
	accountService interfaces.AccountService,
	tokenService interfaces.TokenService,
) (*Operator, error) {
	txqueues := []*TxQueue{}
	addressIndex := make(map[common.Address]*TxQueue)

	wallets, err := walletService.GetOperatorWallets()
	if err != nil {
		panic(err)
	}

	for i, w := range wallets {
		name := strconv.Itoa(i) + w.Address.Hex()
		ch := conn.GetChannel("TX_QUEUES:" + name)

		err := conn.DeclareThrottledQueue(ch, "TX_QUEUES:"+name)
		if err != nil {
			panic(err)
		}

		txq, err := NewTxQueue(
			name,
			tradeService,
			provider,
			orderService,
			w,
			exchange,
			conn,
			accountService,
			tokenService,
		)

		if err != nil {
			panic(err)
		}

		txqueues = append(txqueues, txq)
	}

	op := &Operator{
		WalletService:     walletService,
		TradeService:      tradeService,
		OrderService:      orderService,
		EthereumProvider:  provider,
		Exchange:          exchange,
		TxQueues:          txqueues,
		QueueAddressIndex: addressIndex,
		mutex:             &sync.Mutex{},
	}

	go op.HandleEvents()
	return op, nil
}

// SubscribeOperatorMessages
func (op *Operator) SubscribeOperatorMessages(fn func(*types.OperatorMessage) error) error {
	ch := op.Broker.GetChannel("OPERATOR_SUB")
	q := op.Broker.GetQueue(ch, "TX_MESSAGES")

	go func() {
		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			logger.Error("Failed to register a consumer", err)
		}

		forever := make(chan bool)

		go func() {
			for m := range msgs {
				om := &types.OperatorMessage{}

				err := json.Unmarshal(m.Body, &om)
				if err != nil {
					logger.Error(err)
					continue
				}

				logger.Info(om)

				go fn(om)
			}
		}()

		<-forever
	}()
	return nil
}

func (op *Operator) HandleError(m *types.Matches) {
	err := op.Broker.PublishErrorMessage(m, "Server error")
	if err != nil {
		logger.Error(err)
	}
}

func (op *Operator) HandleTxError(m *types.Matches, id int) {
	errType := getErrorType(id)
	err := op.Broker.PublishTxErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}
}

func (op *Operator) HandleTxSuccess(m *types.Matches, receipt *eth.Receipt) {
	err := op.Broker.PublishTradeSuccessMessage(m)
	if err != nil {
		logger.Error(err)
	}
}

// Bug: In certain cases, the trade channel seems to be receiving additional unexpected trades.
// In the case TestSocketExecuteOrder (in file socket_test.go) is run on its own, everything is working correctly.
// However, in the case TestSocketExecuteOrder is run among other tests, some tradeLogs do not correspond to an
// order hash in the ordertrade mapping. I suspect this is because the event listener catches events from previous
// tests. It might be helpful to see how to listen to events from up to a certain block.
func (op *Operator) HandleEvents() error {
	errorEvents, err := op.Exchange.ListenToErrors()
	if err != nil {
		logger.Error(err)
		return err
	}

	for {
		select {
		case event := <-errorEvents:
			logger.Error("Receiving error event", utils.JSON(event))
			makerOrderHash := event.MakerOrderHash
			takerOrderHash := event.TakerOrderHash
			errID := int(event.ErrorId)

			trades, err := op.TradeService.GetByTakerOrderHash(takerOrderHash)
			if err != nil {
				logger.Error(err)
			}

			to, err := op.OrderService.GetByHash(takerOrderHash)
			if err != nil {
				logger.Error(err)
			}

			makerOrders, err := op.OrderService.GetByHash(makerOrderHash)
			if err != nil {
				logger.Error(err)
			}

			matches := &types.Matches{
				MakerOrders: []*types.Order{makerOrders},
				TakerOrder:  to,
				Trades:      trades,
			}

			go op.HandleTxError(matches, errID)
		}
	}
}

func (op *Operator) HandleTrades(msg *types.OperatorMessage) error {
	err := op.QueueTrade(msg.Matches)
	if err != nil {
		logger.Error(err)
		op.HandleError(msg.Matches)
		return err
	}

	return nil
}

// QueueTrade
func (op *Operator) QueueTrade(m *types.Matches) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	txq, len, err := op.GetShortestQueue()
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("Queuing matches on queue: %v", txq.Name)

	if len > 10 {
		logger.Warning("Transaction queue is full")
		return errors.New("Transaction queue is full")
	}

	logger.Infof("Queuing Trade on queue: %v (previous queue length = %v)", txq.Name, len)

	err = txq.PublishPendingTrades(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetShortestQueue
func (op *Operator) GetShortestQueue() (*TxQueue, int, error) {
	shortest := &TxQueue{}
	min := 1000

	for _, txq := range op.TxQueues {
		if shortest == nil {
			shortest = txq
			min = txq.Length()
		}

		ln := txq.Length()
		if ln < min {
			shortest = txq
			min = ln
		}
	}

	return shortest, min, nil
}

func (op *Operator) PurgeQueues() error {
	for _, txq := range op.TxQueues {
		err := txq.PurgePendingTrades()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (op *Operator) GetTxSendOptions() (*bind.TransactOpts, error) {
	wallet, err := op.WalletService.GetDefaultAdminWallet()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return bind.NewKeyedTransactor(wallet.PrivateKey), nil
}
