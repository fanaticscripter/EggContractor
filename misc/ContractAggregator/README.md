This module implements an aggregator for all Egg, Inc. contracts that ever existed.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Mechanism](#mechanism)
- [Invocation](#invocation)
- [Published data](#published-data)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Mechanism

API response for `/first_contact` corresponding to each player only references contracts attempted by that player (moreover, leggacy attempts supercede original attempts). API response for `/get_periodicals` only references currently active events. Both are unsuitable for the stated goal in isolation, but the latter is hopeless (unless we get ahold of an archive) while the former could be pieced together by looking at multiple players. How do we find players? Well, we either know them directly, or we can harvest loads of them from publicly advertised coops.

## Invocation

Once a config file `config.toml` is placed in the current directory, which looks like

```toml
[database]
# Path for contracts database. Required.
path = "data/contracts.db"

[export]
# Path for exported CSV database.
csv_path = "data/contracts.csv"

[aggregator]
# A list of user IDs blacklisted for unreliable contract information, e.g. for
# cheating on contract expiry timestamps.
player_id_blacklist = []

# A complete list of know contract IDs. Leggacy runs SHOULD be duplicated.
# Used for cross-verification with a third party, e.g.
# https://docs.google.com/spreadsheets/d/1JE5OlFG7tKfj-yXw-rN9fhLzvdQv0YmB_olrYFAAeKo/
# from @mikit on Egg, Inc. Discord.
known_contract_ids = []
```

`ContractAggregator` takes a list of player/coop identifiers. Identifiers of players are simply good old player IDs. Identifiers of coops are specified in the form `contract-id@coop-code`. A sample invocation, with fake IDs and codes:

```console
$ ./ContractAggregator G:1234567890 U:5716bdf11cdfc4769e19cef4a2c7b669b ion-production-2021:fake-code xmas-trash:fake-code"
```

The program would then scan for players (not stored) and contract properties, storing them in the database when found, and eventually generate a CSV export.

## Published data

Latest data is available at [data/contracts.csv](data/contracts.csv).

Note that *offering and expiry datetimes are not reliable, especially for early contracts in 2018.*

The columns:

- **ID**: self-explanatory;
- **Name**: self-explanatory;
- **Type**: Original or Leggacy, self-explanatory;
- **Has Leggacy**: whether this contract has been offered as leggacy (true on leggacy contracts themselves);
- **Offering (Estimated)**: estimated date (UTC) when the contract was unveiled; estimation strategy: expiry minus 3 weeks for originals; expiry minus 1 week for leggacies;
- **Expiry**: earliest expiry timestamp (UTC) we've witnessed accross scraped players; may not be accurate, especially for early contracts in 2018;
- **Egg**: egg type, self-explanatory;
- **Duration**: time to complete the contract;
- **Size**: max coop size, N/A when coop not-allowed;
- **Token**: token interval, N/A when boost tokens weren't a thing yet;
- **Std#1,2,3**: standard tier goals; goals for older contracts without tiers are put in the standard columns;
- **Std Rate/hr**: required egg laying rate to complete the final goal on standard tier;
- **Elt#1,2,3**: elite tier goals; N/A for older contracts without tiers;
- **Elt Rate/hr**: required egg laying rate to complete the final goal on elite tier;
- **#PE**: number of prophecy eggs at stake;
- **Std PE Goal**: goal for prophecy egg(s) reward on standard tier;
- **Std PE Rate/hr**: required egg laying rate to acquire prophecy egg(s) reward on standard tier;
- **Elt PE Goal**: goal for prophecy egg(s) reward on elite tier;
- **Elt PE Rate/hr**: required egg laying rate to acquire prophecy egg(s) reward on elite tier;
- **JSON**: JSON serialization of `api.ContractProperties` (see `api/egginc.pb.go`);
- **Protobuf (base64)**: base64-encoded Protobuf serialization of the `ContractProperties` message (see `api/egginc.proto`).
