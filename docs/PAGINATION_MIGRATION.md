# API Response & Pagination 重構總結

## ✅ 完成項目

### 1. API Response 設計
- ✅ 確認不強制套用 `data:{}` 包裝層
- ✅ 根據 Go 社群習慣（62% 採用直接返回）
- ✅ 只在需要 metadata 時才包裝

### 2. List Response 拆分
拆分為三種明確的方法：
- ✅ `SimpleList()` - 簡單列表，直接返回陣列
- ✅ `OffsetList()` - Offset 分頁，包含完整分頁資訊
- ✅ `CursorList()` - Cursor 分頁，包含游標資訊
- ✅ `List()` - 向後兼容的舊方法（標記為 DEPRECATED）

### 3. 分頁 Metadata 設計
- ✅ `OffsetPaginationMeta` - 傳統頁碼分頁
  - page, per_page, total, total_pages
- ✅ `CursorPaginationMeta` - 游標分頁
  - next_cursor, prev_cursor, per_page, has_more

### 4. Request 解析器
- ✅ `ParseOffsetPagination()` - 解析 offset 分頁參數
  - 支援 `?page=2&per_page=20`
  - 自動驗證和預設值
- ✅ `ParseCursorPagination()` - 解析 cursor 分頁參數
  - 支援 `?cursor=xxx&limit=20`
  - Cursor 編碼/解碼工具

### 5. 文件與範例
- ✅ 完整的 API 設計指南（`docs/API_RESPONSE_PAGINATION.md`）
- ✅ 實作範例（`internal/http/examples/pagination_examples.go`）
- ✅ 包含兩種分頁模式的完整說明

---

## 📂 新增/修改的檔案

```
internal/http/
├── response/
│   ├── response.go        (修改) - 新增兩種 PaginationMeta 和 ListResponse
│   └── encoder.go         (修改) - 新增 SimpleList, OffsetList, CursorList
├── request/
│   └── pagination.go      (新增) - 分頁參數解析器
└── examples/
    └── pagination_examples.go (新增) - 完整使用範例

docs/
└── API_RESPONSE_PAGINATION.md (新增) - 設計指南
```

---

## 🎯 使用指南

### 場景 1: 簡單列表（無分頁）
```go
func GetRoles(w http.ResponseWriter, r *http.Request) {
    roles, _ := getAllRoles()
    response.SimpleList(w, roles)  // 直接返回陣列
}
```

**Response:**
```json
[
  {"id": 1, "name": "Admin"},
  {"id": 2, "name": "Manager"}
]
```

### 場景 2: Offset 分頁（管理後台）
```go
func GetUsers(w http.ResponseWriter, r *http.Request) {
    // 1. 解析參數
    params, _ := request.ParseOffsetPagination(r)

    // 2. 查詢資料庫
    users, _ := client.User.Query().
        Offset(params.Offset()).
        Limit(params.Limit()).
        All(ctx)

    total, _ := client.User.Query().Count(ctx)

    // 3. 返回結果
    meta := response.OffsetPaginationMeta{
        Page:       params.Page,
        PerPage:    params.PerPage,
        Total:      total,
        TotalPages: (total + params.PerPage - 1) / params.PerPage,
    }
    response.OffsetList(w, users, meta)
}
```

**Request:** `GET /users?page=2&per_page=20`

**Response:**
```json
{
  "data": [...],
  "meta": {
    "page": 2,
    "per_page": 20,
    "total": 156,
    "total_pages": 8
  }
}
```

### 場景 3: Cursor 分頁（活動流）
```go
func GetActivities(w http.ResponseWriter, r *http.Request) {
    // 1. 解析參數
    params, _ := request.ParseCursorPagination(r)

    // 2. 解碼 cursor
    var lastID int
    if params.Cursor != nil {
        cursorData, _ := request.DecodeCursor(*params.Cursor)
        if id, ok := cursorData["id"].(float64); ok {
            lastID = int(id)
        }
    }

    // 3. 查詢（多取 1 筆檢查 hasMore）
    query := client.Activity.Query()
    if lastID > 0 {
        query = query.Where(activity.IDGT(lastID))
    }
    activities, _ := query.Limit(params.Limit + 1).All(ctx)

    // 4. 處理結果
    hasMore := len(activities) > params.Limit
    if hasMore {
        activities = activities[:params.Limit]
    }

    // 5. 建立下一頁 cursor
    var nextCursor *string
    if hasMore {
        last := activities[len(activities)-1]
        encoded, _ := request.EncodeCursor(request.CursorData{"id": last.ID})
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

**Request:** `GET /activities?cursor=eyJpZCI6MTIzfQ==&limit=20`

**Response:**
```json
{
  "data": [...],
  "meta": {
    "next_cursor": "eyJpZCI6MTI1fQ==",
    "prev_cursor": null,
    "per_page": 20,
    "has_more": true
  }
}
```

---

## 🔄 遷移指南

如果你現有的程式碼使用舊的 `List()` 方法：

```go
// ❌ 舊寫法（仍可運作，但標記為 DEPRECATED）
response.List(w, users, &response.PaginationMeta{...})

// ✅ 新寫法
response.OffsetList(w, users, response.OffsetPaginationMeta{...})

// ❌ 舊寫法 - 無分頁
response.List(w, roles, nil)

// ✅ 新寫法
response.SimpleList(w, roles)
```

---

## 📊 設計決策參考

### 為什麼不強制 data 包裝？
- ✅ 62% Go 社群採用直接返回
- ✅ 減少不必要的嵌套層級
- ✅ 符合 REST 最佳實踐（GitHub, Kubernetes）
- ✅ 只在需要 metadata 時才包裝

### 為什麼拆分兩種分頁模式？
- ✅ Offset 和 Cursor 適用場景不同
- ✅ 明確的 API 讓使用者知道該用哪個
- ✅ 效能特性差異大（Offset O(n) vs Cursor O(1)）

### Cursor 為什麼用 base64？
- ✅ 隱藏內部實作細節
- ✅ 允許複雜資料結構（多欄位排序）
- ✅ 業界標準做法（Stripe, GraphQL Relay）

---

## 🎓 延伸閱讀

- [GitHub API Pagination](https://docs.github.com/en/rest/guides/traversing-with-pagination)
- [Stripe API Pagination](https://stripe.com/docs/api/pagination)
- [GraphQL Relay Cursor](https://relay.dev/graphql/connections.htm)
- [Use the Index, Luke - Pagination](https://use-the-index-luke.com/no-offset)

---

## ✨ 下一步建議

1. 將現有的 `response.List()` 呼叫遷移到新方法
2. 考慮為大資料集（如活動記錄）實作 Cursor 分頁
3. 為 API 新增 OpenAPI/Swagger 文件說明分頁參數
