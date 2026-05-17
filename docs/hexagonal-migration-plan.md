# Hexagonal Folder Structure Migration Plan

## 1. 目標

將目前偏 layered 的專案結構，逐步整理為 hexagonal architecture，讓核心商業邏輯獨立於 HTTP、Ent、Redis、DB 與其他外部系統，降低耦合、提升可測試性與長期維護性。

## 2. 範圍

本次計劃只處理 folder 結構、依賴方向與程式碼分層，不包含業務需求變更。

### 會涵蓋的內容

- HTTP entrypoint 與路由組裝。
- 用例與商業邏輯。
- Inbound / outbound ports。
- HTTP、Repository、Cache、Worker 等 adapters。
- Wire DI 組裝方式。
- 測試結構與 mock 位置。

### 不會做的內容

- 不改業務規則。
- 不改資料表結構。
- 不做不必要的大規模重寫。
- 不一次性搬完所有功能模組。

## 3. 現況判讀

目前專案已經有一些 hexagonal 的基礎元素，例如：

- `internal/ports`
- `internal/core`
- `internal/adapters`
- `cmd/server`、`cmd/worker`、`cmd/migrate` 作為 entrypoints

但整體仍較像傳統 layered design：

- HTTP controller、service、repository 的責任邊界還不夠明確。
- 一些跨層依賴仍直接指向基礎設施型實作。
- domain / use case / adapter 的命名與資料流還沒有完全對齊 hexagonal。

## 4. 目標結構

以下是建議的目標方向，會以現有資料夾為基礎漸進調整：

| 目錄 | 角色 |
|---|---|
| `cmd/server` | HTTP 程式啟動與組裝入口 |
| `cmd/worker` | Worker 啟動與組裝入口 |
| `cmd/migrate` | migration 工具入口 |
| `internal/core` | domain 與 use case 的核心商業邏輯 |
| `internal/ports` | inbound / outbound interfaces |
| `internal/adapters/http` | HTTP controllers、routes、middleware |
| `internal/adapters/persistence` | Ent repository 與資料存取 |
| `internal/adapters/cache` | Redis、快取、idempotency、lock 相關 adapter |
| `internal/adapters/worker` | 背景工作與排程 adapter |
| `internal/infra` | 共用基礎設施工具與 cross-cutting utilities |
| `internal/dtos` | request / response schema，必要時可逐步下放到 adapters |
| `internal/domains` | 純 domain model 與 entity |

## 5. 遷移原則

- 核心規則只往內依賴，不往外依賴。
- use case 不直接依賴 HTTP、Ent、Redis client 或 chi router。
- port 先定義契約，再實作 adapter。
- 先搬邊界，再搬核心，避免大爆炸式重構。
- 每次遷移都要保留可編譯、可測試狀態。
- Wire 組裝只負責把依賴接起來，不放業務邏輯。

## 6. Git 記錄保留策略

這次遷移不只要改結構，也要盡量保留 `git blame` 與可追溯性。Git 的 rename 判斷本質上是相似度推測，不是絕對保證，所以流程要刻意設計成「讓 Git 容易看出這是搬移，而不是刪除後重寫」。

### 原則

- 先搬檔案，再改內容。
- 一個 commit 只做一種性質的變更。
- 盡量讓 rename 與內容重構分開。
- 大型模組搬遷時，先保留舊介面與包裝層，再逐步縮小。
- 每個 commit 都要能獨立通過 build / test。

### 實作方式

- 使用 `git mv` 或 IDE 的 move / refactor 功能來移動檔案，不要用刪除後重建。
- 移動時盡量維持檔案內容接近原樣，只先改路徑與 import。
- 若某個檔案要大改，先分成兩步：第一步只搬位置，第二步再改邏輯。
- 避免在同一個 commit 中同時做大量改名、拆檔、重寫內容。
- 若一次要拆成多個檔案，先保留一個薄的轉接層，讓舊路徑能持續工作，再慢慢移除。

### 建議的 commit 節奏

1. `chore`: 只搬檔案與調整 import，盡量不改行為。
2. `refactor`: 調整 package 邊界與介面，但保持測試通過。
3. `feat` / `fix`: 只在結構穩定後再做真正的行為變更。

### 驗證方式

- 每個 commit 後先跑對應模組的測試，再跑全域 `make test` 或至少相關子測試。
- 在 code review 時用 `git diff -M -C` 檢查 rename 偵測結果。
- 若 Git 沒有自動認出 rename，就代表這個 commit 的內容變化太大，應拆小。
- 追查歷史時可用 `git log --follow`，但不要把它當成完全可靠的保證。

## 7. 分階段計劃

### Phase 0：盤點與切界

- 盤點目前所有 package 的依賴方向。
- 標記哪些是 core、哪些是 adapter、哪些其實應該移到 port。
- 找出最容易切開的 bounded context，優先從 auth、user、approval 這類模組開始。
- 建立 folder 對照表，避免搬遷時命名漂移。

### Phase 1：定義核心邊界

- 整理 `internal/core` 的責任，讓它只保留商業邏輯與 use case。
- 把跨層通訊改成 ports 介面。
- 明確區分 inbound ports 與 outbound ports。
- 清理 core 對外部 framework 的直接依賴。

### Phase 2：抽出 adapters

- 把 HTTP controller、route、middleware 收斂到 HTTP adapter。
- 把 Ent repository 收斂到 persistence adapter。
- 把 Redis、cache、idempotency、lock 收斂到 infra 或 cache adapter。
- 讓 adapter 只做轉換、包裝、呼叫，不做核心決策。

### Phase 3：調整依賴注入

- 重新整理 Wire provider 的分組。
- 讓 entrypoint 只組裝，不承擔業務判斷。
- 檢查是否有重複或過度耦合的 provider。
- 確保新增功能時只需要新增 port + adapter + use case，而不是跨很多層改動。

### Phase 4：逐模組搬遷

- 先搬一個低風險模組做樣板。
- 以同樣模式搬遷 user、auth、approval、worker job。
- 每搬一個模組就補齊對應測試。
- 每次遷移完成都跑一次 build / test 驗證。

### Phase 5：文件與規範同步

- 更新 README 的專案結構說明。
- 更新 `docs/testing.md` 的測試分類與位置。
- 新增一份 hexagonal 架構導覽文件，說明各層責任與依賴方向。
- 補上貢獻者的目錄規則與新增功能範本。

## 8. 建議的任務拆解

- 建立現況依賴圖。
- 定義 hexagonal 的目錄責任表。
- 先抽出 ports 契約。
- 重整第一個模組的 use case 與 adapter。
- 調整 Wire provider 與 entrypoint 組裝。
- 補測試並修正命名。
- 同步 README 與架構文件。
- 逐步搬遷剩餘模組。

## 9. 驗收標準

- core 不直接依賴 HTTP、Ent、Redis 或 chi。
- 所有外部互動都經過 ports。
- adapter 可替換，不影響核心邏輯。
- 新模組加入時有一致的新增路徑。
- 測試能在沒有真實外部依賴的情況下覆蓋核心邏輯。
- README 與 docs 能正確反映新結構。

## 10. 風險與對策

- 風險：一次搬太多會造成編譯中斷。對策：採漸進式遷移，每次只搬一個模組。
- 風險：ports 定義不穩會反覆改動。對策：先用最小契約，避免過度抽象。
- 風險：Wire 組裝變複雜。對策：保持 provider 分層清楚，entrypoint 只做組裝。
- 風險：文件和程式碼不同步。對策：每次搬遷都同步更新 docs。

## 11. 推薦執行順序

1. 先盤點現況與依賴圖。
2. 先定義目標 folder 與責任邊界。
3. 先搬一個模組做範例。
4. 再批次遷移其餘模組。
5. 最後統一更新文件與測試說明。
