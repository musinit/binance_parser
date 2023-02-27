# Binance parser

### Getting started
1. Start the service
```azure
    make run
```
3. Check how it works

```azure
POST localhost:9090/parser/subscription (subscribe)
GET localhost:9090/parser/blocks/current (get latest block number)
GET localhost:9090/parser/transactions (check transactions)

```

Limits:
 - Check new blocks once in 500 sec
 - Updating and caching user transactions with the same rate
 - Start from 16717904 block number (can be changed in repository layer)
 - After subscribing start transaction caching with the latest block num
