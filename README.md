https://user-images.githubusercontent.com/30017432/226439153-c0429cd6-747d-4cf6-be84-320738a1ef72.mp4

### Q1 .Explain what an eventually consistent ledger would need to look out for, what are some of the CAP theorem and database considerations that are relevant when designing a bank ledger.

An eventually consistent ledger means that there are risks of overdraft.

### Q2. Explain your solution, how does matching work? What will scaling look like? How would you improve the API beyond a toy implementation?





<details>
<summary>Solution</summary>

Prevent customer from overspent using `debits_must_not_exceed_credits` flag for new customer account in `tigerBeetle`

Uses `Custom Search Attribute (CSA)` from temporal.

Whenever a request for `authorization` came in, the following actions are made

1. Make a `pending transfer` in `tigerBettle`
2. Start a long running `ExpireAuthorization` workflow

   2.1 attach `TransactionPendingAmount` and `TransactionUserId` as `CSA` into the workflow

   2.2 Start workflow timer, the `expiration activity` will only run as a `Future`

   2.3 Workflow constantly listen for `cancellation signal`

When a `presenment` is made, we first search using `CSA`,

- if there is a match, then we take the earliest workflow, and send a `cancellation signal`, then `release payment`
- else, we consider it as `instant transfer`
  </details>

For Scaling, we will need to abstract the exposed API to consumer of the ledger. We need to follow the philosophy of `batching` from `tigerBeetle`

The parameters for batching will depends largely on the volume that we want to serve, probably by fixed interval + batch size + priority of transaction

Move more business logic into workflows, authorization should probably be workflow. This would make handling of business logic very granular, and there is full observability.

## Test Steps

```bash
# Create A new account with some id
curl 'http://localhost:4000/ledger' -d '{"Id":"888"}'

#  Add some balance into the account
curl 'http://localhost:4000/ledger/888/addBalance' -d '{"Amount":1000}'

#  Authorize Payment
curl 'http://localhost:4000/ledger/888/authorize' -d '{"Amount":25}'


# At this stage, the balance should looks like
curl 'http://localhost:4000/ledger/888/balance'

'
{
	"Id": "888",
	"Balance": {
		"Available": 975,
		"Reserved":  25
	}
}
'

# If left for expiration, available balance will return to 1000

# If presentment is made via,
curl 'http://localhost:4000/ledger/888/presentment' -d '{"Amount":25}'

# New balance will looks like
'
{
	"Id": "888",
	"Balance": {
		"Available": 975,
		"Reserved":  0
	}
}
'

```

## Temporal Add Search Attribute

```bash
tctl adm cl asa -n TransactionPendingAmount -t Int

tctl adm cl asa -n TransactionUserId -t Keyword

tctl admin cl gsa
```
