# PLATFORM-2984: Refactor load shedding and admission control

**Status:** In Progress · **Priority:** Medium
**Sprint:** Sprint 30 · **Story Points:** 5
**Reporter:** Suresh Kumar (Infra Lead) · **Assignee:** You (Intern)
**Due:** End of sprint (Friday)
**Labels:** `backend`, `golang`, `resilience`, `performance`
**Task Type:** Code Maintenance

---

## Description

The load shedding system works but has code quality issues from the last review. Refactor without changing external behavior.

## Acceptance Criteria

- [ ] Magic numbers replaced with named constants
- [ ] Priority logic simplified and documented
- [ ] Redundant mutex double-locking fixed
- [ ] All tests still pass
