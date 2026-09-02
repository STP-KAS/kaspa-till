# Vision — stables landing on Kaspa L1

Avoid L2. Do not bill in grams. Do not mint a dollar from a script.

Kaspa Till is a shop whose **unit of account is already the future L1 stable**. The asset is reserved (`kUSD`, 6 decimals, scheme `kaspa-l1-stable`). Until a real KCC-20 (or overcollateral vault) exists, the only live spend is **KAS** at a number the merchant typed.

## Why a second dApp

Gramlane (the other dApp) solved **fee volatility** by invoicing **work**. That is possible today with covenants. It is not a dollar.

A coffee, a subscription, a payroll cannot honestly be “50,000 grams.” Those want a unit of account that tracks **purchasing power**. That object, on Kaspa L1, is not live. Building the till now means:

- the catalog never stores KAS prices as the meaning of the good
- HTTP 402 already has a hole for the future asset
- we refuse the shortcut of an L2 bridged USDC

## How it can land (L1 only)

1. **Issued KCC-20** — a minter covenant on L1. Still needs **capital** (the issuer’s reserves). Wallets learn the template hash. This till’s `accepts[1]` becomes live.
2. **Overcollateral vault** — users lock KAS, mint kUSD under a ratio, oracle + liquidation. Still **capital** (crypto). Still L1.
3. **KIP that changes miner fees** — not this app.

Refused: algorithmic kas-USD, L2 bridges, calling Work Credits a dollar.

## What this process does today

- Prices: integer micro-units of reserved kUSD
- Settlement: KAS sompi = micro × (merchant sompi-per-unit) / 1e6
- 402: live KAS + reserved kUSD
- Orders: memory only, HTTP receipt, not a chain transfer of kUSD

The merchant rate is a **sign on the counter**. It is not Circle and not a DEX TWAP.
