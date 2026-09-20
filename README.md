> **Experimental only. Not a product.** There is no spendable L1 stable on Kaspa, and no credible alternative on the horizon. Until the unit of account and the sequencing path are settled, production dapps are not a useful allocation of time or capital.
>
> Do not use wallet integrations on this GitHub. STP remains a clown. [DISCLAIMER.md](DISCLAIMER.md)

# Kaspa Till — dApp 2 (L1 stable vision)

**project delusional** · [@StppStp](https://x.com/StppStp)

A **Kaspa L1** merchant. The catalog is priced in a **reserved native stable** (`kUSD`). That asset is **not live**. Today you settle in **KAS** at a merchant-posted rate.

**No L2. No work credits. No fake peg.**

This is the dApp you run if the product is “stables will land on Kaspa,” not “replace the dollar with grams.”

Sister dApp (grams / Work Credits): `C:\Users\<user>\Documents\kaspa\superapp` — Gramlane on `:8081`.

```powershell
cd C:\Users\<user>\Documents\kaspa\superappstablesalternative
go test ./...
go build -o kastill.exe ./cmd/kastill
.\kastill.exe
```

http://localhost:8082

If the browser says “localhost refused to connect”:

```powershell
powershell -File C:\Users\<user>\Documents\kaspa\start-local.ps1
```

Index: [STP-KAS/project-delusional](https://github.com/STP-KAS/project-delusional).

| Path | What |
| --- | --- |
| `/idea` | What this URL is |
| `/why` | Beyond the chain: shelf in money, chain not yet |
| `/shop` | Shelf in reserved kUSD |
| `/item/cup` | Dual invoice: kUSD reserved + KAS due now |
| `/vision` | How a stable lands on L1 (KCC-20 or vault). L2 refused. |
| `/rate` | Merchant sign (sompi per 1.00). Not an oracle. |
| `/api/order?item=cup` | 402 with `kaspa` (live) and `kaspa-l1-stable` (not live) |
| `/wallets` | Same Kaspa wallet catalog. Connect + log out. |
| `/safety` | Never DMs, never seeds. |
| `/234` | Why we do not `readInputState` a kUSD UTXO. Amount 1 → vault 264. |
| `/feedback` | Stored on this PC under `Documents\kaspa\feedback\kastill` |

`KasInvoice.sil` compiled with official silverc v1-rc1 (KAS due now + reserved micro-amount in the constructor). Not deployed.

Read [VISION.md](VISION.md) next.

---

> **Standard disclaimer.** This GitHub, not the topic above.
>
> Intentions are good; thought process is questionable. STP remains delusional. Si vis pacem, para bellum.
>
> Intern at https://sixpack.wtf/  
> X: https://x.com/StppStp · GitHub: https://github.com/STP-KAS
