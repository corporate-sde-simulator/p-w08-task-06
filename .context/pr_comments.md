# PR Review - Load shedding and admission control (by Deepak)

## Reviewer: Nisha Gupta
---

**Overall:** Good foundation but critical bugs need fixing before merge.

### `admissionController.go`

> **Bug #1:** Load estimation uses a stale snapshot instead of moving average and reacts to old data
> This is the higher priority fix. Check the logic carefully and compare against the design doc.

### `loadEstimator.go`

> **Bug #2:** Priority-based shedding drops high-priority requests first instead of low-priority
> This is more subtle but will cause issues in production. Make sure to add a test case for this.

---

**Deepak**
> Acknowledged. I have documented the issues for whoever picks this up.
