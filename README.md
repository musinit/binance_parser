# Binance parser

### Getting started
1. Start the service
```azure
    make run
```
3. Check how it works

```azure
GET localhost:9090/parser/blocks/current
GET localhost:9090/parser/transactions
POST localhost:9090/parser/subscription
```

Limits:
 - Check address once per ~5 seconds
 - Checks <= 100 latest transaction in ~5 seconds

