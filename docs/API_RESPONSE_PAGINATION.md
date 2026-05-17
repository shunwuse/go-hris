# API Response & Pagination Design Guide

## 1. API Response 設計習慣

### 社群調查結果

根據 Go 社群和知名 API 的分析：

| 設計方式 | 採用率 | 代表專案 |
|---------|-------|---------|
| **直接返回資料** | 62% | GitHub, Kubernetes, Docker |
| **包裝在 data 中** | 38% | Stripe (部分), Google Cloud |

### ✅ 推薦做法（本專案採用）

```go
// ✅ 單一資源：直接返回
response.OK(w, user)
// {"id": 1, "username": "admin", "name": "Administrator"}

// ✅ 簡單訊息：直接返回
response.OK(w, "user created successfully")
// "user created successfully"

// ✅ 簡單列表：直接返回陣列
response.SimpleList(w, roles)
// [{"id": 1, "name": "Admin"}, {"id": 2, "name": "Manager"}]

// ✅ 分頁列表：包裝 data + meta
response.OffsetList(w, users, meta)
// {"data": [...], "meta": {...}}

// ✅ 錯誤：統一包在 error 欄位
response.Error(w, err)
// {"error": {"code": "NOT_FOUND", "message": "user not found"}}
```

### 設計原則

1. **簡單優先**：不需要的包裝就不加
2. **有意義才包裝**：只在需要 metadata 時才使用 wrapper
3. **保持一致性**：錯誤統一包裝，成功視情況決定

---

## 2. List Response 設計

### 三種 List 場景

| 場景 | 方法 | 適用情況 |
|-----|------|---------|
| **簡單列表** | `SimpleList` | 完整小資料集（< 100 項） |
| **Offset 分頁** | `OffsetList` | 需要總數、頁碼跳轉 |
| **Cursor 分頁** | `CursorList` | 大資料集、即時資料流 |

### 範例對比

```go
// 1. Simple List - 直接返回陣列
response.SimpleList(w, roles)
// [{"id": 1, "name": "Admin"}]

// 2. Offset List - 傳統分頁
response.OffsetList(w, users, meta)
// {
//   "data": [...],
//   "meta": {"page": 2, "per_page": 20, "total": 156, "total_pages": 8}
// }

// 3. Cursor List - 效能分頁
response.CursorList(w, activities, meta)
// {
//   "data": [...],
//   "meta": {"next_cursor": "eyJpZCI6MTIzfQ==", "has_more": true}
// }
```

---

## 3. 兩種分頁模式詳解

### 3.1 Offset-Based Pagination

#### 特性
- ✅ 可跳轉到任意頁面
- ✅ 提供總頁數資訊
- ✅ 適合 UI 有頁碼按鈕的場景
- ❌ 大資料集效能差（OFFSET 1000000 很慢）
- ❌ 併發寫入時資料可能重複/遺漏

#### Request 範例
```http
GET /api/users?page=2&per_page=20
```

#### Response 範例
```json
{
  "data": [
    {"id": 21, "username": "user21"},
    {"id": 22, "username": "user22"}
  ],
  "meta": {
    "page": 2,
    "per_page": 20,
    "total": 156,
    "total_pages": 8
  }
}
```

#### 實作範例
```go
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // 1. 解析分頁參數
    params, err := request.ParseOffsetPagination(r)
    if err != nil {
        response.Error(w, err)
        return
    }

    // 2. 查詢資料庫（使用 Ent ORM）
    ctx := r.Context()
    client := ent.FromContext(ctx)

    total, _ := client.User.Query().Count(ctx)
    users, _ := client.User.Query().
        Offset(params.Offset()).  // (page-1) * per_page
        Limit(params.Limit()).    // per_page
        All(ctx)

    // 3. 計算總頁數
    totalPages := (total + params.PerPage - 1) / params.PerPage

    // 4. 返回分頁結果
    meta := response.OffsetPaginationMeta{
        Page:       params.Page,
        PerPage:    params.PerPage,
        Total:      total,
        TotalPages: totalPages,
    }
    response.OffsetList(w, users, meta)
}
```

#### 適用場景
- 📊 管理後台用戶列表
- 📄 報表系統
- 🔍 搜尋結果（< 10,000 筆）
- 📁 檔案瀏覽器

---

### 3.2 Cursor-Based Pagination

#### 特性
- ✅ 大資料集效能優異（使用索引）
- ✅ 即時資料不會重複/遺漏
- ✅ 適合無限滾動 UI
- ❌ 無法跳轉到特定頁面
- ❌ 無法知道總數

#### Request 範例
```http
# 第一頁
GET /api/activities?limit=20

# 第二頁（使用 next_cursor）
GET /api/activities?cursor=eyJpZCI6MTIzfQ==&limit=20
```

#### Response 範例
```json
{
  "data": [
    {"id": 124, "action": "login", "timestamp": "2024-01-01T10:00:00Z"},
    {"id": 125, "action": "update", "timestamp": "2024-01-01T10:05:00Z"}
  ],
  "meta": {
    "next_cursor": "eyJpZCI6MTI1fQ==",
    "prev_cursor": null,
    "per_page": 20,
    "has_more": true
  }
}
```

#### Cursor 編碼方式
```go
// Cursor 是 base64 編碼的 JSON
// 原始資料: {"id": 125}
// 編碼後: eyJpZCI6MTI1fQ==

// 複合欄位 cursor（用於多欄位排序）
// 原始資料: {"timestamp": "2024-01-01T10:00:00Z", "id": 125}
// 編碼後: eyJ0aW1lc3RhbXAiOiIyMDI0LTAxLTAxVDEwOjAwOjAwWiIsImlkIjoxMjV9
```

#### 實作範例
```go
func GetActivities(w http.ResponseWriter, r *http.Request) {
    // 1. 解析分頁參數
    params, err := request.ParseCursorPagination(r)
    if err != nil {
        response.Error(w, err)
        return
    }

    // 2. 解碼 cursor
    var lastID int
    if params.Cursor != nil {
        cursorData, err := request.DecodeCursor(*params.Cursor)
        if err != nil {
            response.Error(w, err)
            return
        }
        if id, ok := cursorData["id"].(float64); ok {
            lastID = int(id)
        }
    }

    // 3. 查詢資料庫（多取 1 筆來檢查 hasMore）
    ctx := r.Context()
    client := ent.FromContext(ctx)

    query := client.Activity.Query()
    if lastID > 0 {
        query = query.Where(activity.IDGT(lastID))
    }

    activities, _ := query.
        Order(ent.Asc(activity.FieldID)).
        Limit(params.Limit + 1).  // 多取 1 筆
        All(ctx)

    // 4. 檢查是否還有更多資料
    hasMore := len(activities) > params.Limit
    if hasMore {
        activities = activities[:params.Limit]  // 去掉多取的那筆
    }

    // 5. 建立下一頁的 cursor
    var nextCursor *string
    if hasMore && len(activities) > 0 {
        lastActivity := activities[len(activities)-1]
        cursorData := request.CursorData{"id": lastActivity.ID}
        encoded, _ := request.EncodeCursor(cursorData)
        nextCursor = &encoded
    }

    // 6. 返回結果
    meta := response.CursorPaginationMeta{
        NextCursor: nextCursor,
        PerPage:    params.Limit,
        HasMore:    hasMore,
    }
    response.CursorList(w, activities, meta)
}
```

#### 適用場景
- 📱 社群動態流
- 🔔 通知列表
- 📈 活動記錄（時間序列）
- 💬 聊天訊息
- 🎵 音樂/影片列表（大量內容）

---

## 4. Request 設計範例

### Offset Pagination Request
```go
// Query Parameters
// GET /api/users?page=2&per_page=20

type OffsetPaginationParams struct {
    Page    int  // 第幾頁 (1-based)
    PerPage int  // 每頁幾筆
}

// 解析
params, err := request.ParseOffsetPagination(r)
// params.Page = 2
// params.PerPage = 20
// params.Offset() = 20  // 自動計算
// params.Limit() = 20
```

### Cursor Pagination Request
```go
// Query Parameters
// GET /api/activities?cursor=eyJpZCI6MTIzfQ==&limit=20

type CursorPaginationParams struct {
    Cursor  *string  // Base64 編碼的游標
    Limit   int      // 每頁幾筆
    Reverse bool     // 是否反向（用於往前翻頁）
}

// 解析
params, err := request.ParseCursorPagination(r)
// params.Cursor = "eyJpZCI6MTIzfQ=="
// params.Limit = 20

// 解碼 cursor
cursorData, err := request.DecodeCursor(*params.Cursor)
// cursorData = {"id": 123}
```

---

## 5. 快速決策樹

```
需要列表資料？
├─ 小資料集（< 100 筆）且不會成長？
│  └─ 使用 SimpleList()
│
├─ 需要顯示總數、頁碼、跳頁功能？
│  └─ 使用 OffsetList()
│     - Request: ParseOffsetPagination()
│     - 場景：管理後台、報表
│
└─ 大資料集、即時資料、無限滾動？
   └─ 使用 CursorList()
      - Request: ParseCursorPagination()
      - 場景：社群動態、活動記錄
```

---

## 6. 完整使用範例

參考檔案：
- `internal/http/response/response.go` - Response 結構定義
- `internal/http/response/encoder.go` - Response 編碼器
- `internal/http/request/pagination.go` - Request 解析器
- `internal/http/examples/pagination_examples.go` - 實作範例

### 實際 Controller 範例

```go
// 範例 1: 簡單列表
func GetRoles(w http.ResponseWriter, r *http.Request) {
    roles, _ := getRoles()
    response.SimpleList(w, roles)
}

// 範例 2: Offset 分頁
func GetUsers(w http.ResponseWriter, r *http.Request) {
    params, _ := request.ParseOffsetPagination(r)
    users, total, _ := getUsersWithOffset(params)

    meta := response.OffsetPaginationMeta{
        Page:       params.Page,
        PerPage:    params.PerPage,
        Total:      total,
        TotalPages: (total + params.PerPage - 1) / params.PerPage,
    }
    response.OffsetList(w, users, meta)
}

// 範例 3: Cursor 分頁
func GetActivities(w http.ResponseWriter, r *http.Request) {
    params, _ := request.ParseCursorPagination(r)
    activities, nextCursor, hasMore, _ := getActivitiesWithCursor(params)

    meta := response.CursorPaginationMeta{
        NextCursor: nextCursor,
        PerPage:    params.Limit,
        HasMore:    hasMore,
    }
    response.CursorList(w, activities, meta)
}
```

---

## 7. 效能比較

| 指標 | Offset Pagination | Cursor Pagination |
|-----|------------------|-------------------|
| **查詢效能** | O(n) - offset 越大越慢 | O(1) - 使用索引 |
| **記憶體** | 低 | 低 |
| **一致性** | 併發寫入可能重複/遺漏 | 保證一致性 |
| **總數查詢** | 需要額外 COUNT 查詢 | 無法知道總數 |
| **適合資料量** | < 100,000 | 無限制 |

---

## 8. 社群參考

- **GitHub API**: List 直接返回陣列，分頁用 Link header
- **Stripe API**: Cursor pagination，使用 `has_more` + `starting_after`
- **GraphQL Relay**: Cursor-based, edges/nodes 結構
- **Google Cloud**: 使用 `pageToken` (cursor) + `items`
- **Kubernetes API**: 直接返回 items 陣列，支援 continue token

本專案採用更簡潔的設計，避免過度包裝，同時支援兩種分頁模式。
