# Frontend Pinia Memory Cache Design

**Date:** 2026-05-31
**Status:** Approved

## Problem

Every page navigation triggers a fresh API call via `onMounted`. Even though the backend responds in <3ms, the lifecycle of route-change → mount → fetch → render spans multiple microtask frames, producing a visible "loading flash". Three Pinia stores already exist but are unused — all pages bypass them and call the API layer directly.

## Solution

Activate the existing Pinia stores with a consistent `ensureLoaded / invalidate / refresh` pattern. First visit fetches from API and caches in memory; subsequent visits serve instantly from store. Write operations invalidate the relevant cache so the next visit fetches fresh data.

## Architecture

```
Page component
  │ calls store.ensureLoaded()
Pinia Store
  ├─ loaded=true  → return cached data (no network)
  └─ loaded=false → API call → store data → loaded=true

Write operation (create/edit/delete) succeeds
  │ calls store.invalidate()
  └─ sets loaded=false, clears data

Logout
  └─ invalidates all stores
```

## Store Interface (uniform across all 4 stores)

| Method | Signature | Behavior |
|--------|-----------|----------|
| `ensureLoaded` | `() => Promise<void>` | If `loaded`, return immediately. Otherwise fetch, populate state, set `loaded=true`. |
| `invalidate` | `() => void` | Set `loaded=false`, clear state arrays/totals. |
| `refresh` | `() => Promise<void>` | Call `invalidate()` then `ensureLoaded()`. |

## Stores to Modify or Create

| Store | Current State | Action |
|-------|--------------|--------|
| `useCategoryStore` | Exists, has `loaded` + `fetchCategories` | Add `invalidate()`, `refresh()`; rename `fetchCategories` to `ensureLoaded` |
| `useSubscriptionStore` | Exists, no `loaded` | Add `loaded`, `invalidate()`, `refresh()`; adapt `fetchSubscriptions` into `ensureLoaded` |
| `useAssetStore` | Exists, no `loaded` | Same pattern as subscription store |
| `useDashboardStore` | Does not exist | Create new, same pattern, cache `DashboardData` |

## Pages to Modify

| Page | Change |
|------|--------|
| `dashboard/Index.vue` | Replace direct `getDashboard()` with `useDashboardStore().ensureLoaded()`; bind template to store state |
| `category/Index.vue` | Replace direct `listCategories()` with `useCategoryStore().ensureLoaded()`; bind template to store state |
| `subscription/Index.vue` | Replace direct `listSubscriptions()` with `useSubscriptionStore().ensureLoaded()`; bind template to store state |
| `asset/Index.vue` | Replace direct `listAssets()` with `useAssetStore().ensureLoaded()`; bind template to store state |
| `reminder/Index.vue` | **No change** — reminders are dynamically generated server-side, always fetch fresh |
| `subscription/Edit.vue` | On save/delete success → `useSubscriptionStore().invalidate()` |
| `asset/Edit.vue` | On save/delete success → `useAssetStore().invalidate()` |
| `category/Index.vue` (modal save/delete) | On success → `store.invalidate()` |
| `subscription/Detail.vue` | On delete success → `useSubscriptionStore().invalidate()` |
| `asset/Detail.vue` | On delete success → `useAssetStore().invalidate()` |
| `stores/user.ts` | In `logout()` → invalidate all 4 stores |

## Cache Invalidation Triggers

| Trigger | Stores Invalidated |
|---------|-------------------|
| Category create/edit/delete | `useCategoryStore` |
| Subscription create/edit/delete | `useSubscriptionStore` |
| Asset create/edit/delete | `useAssetStore` |
| User logout | All 4 stores |
| Page refresh (F5) | All 4 stores (memory-only, wiped on reload) |

## What We Are NOT Doing

- No localStorage/sessionStorage persistence
- No TTL or time-based expiry
- No caching of individual detail pages (single-record fetch is already instant)
- No caching of reminders (server-generated, must be fresh)
- No new npm dependencies
- No Service Worker or HTTP-level caching

## Expected Outcome

| Scenario | Before | After |
|----------|--------|-------|
| First visit to page | API call (~3ms) | API call (~3ms) — same |
| Return to already-visited page | API call + loading flash | **0ms, instant render from store** |
| After create/edit/delete | Next visit shows stale data until manual reload | Next visit fetches fresh data (cache invalidated) |
| Page refresh | Re-fetches everything | Re-fetches everything (memory cleared) — same |
