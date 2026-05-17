# Go 文字風格指南 (Text Style Guide)

本文件定義專案中註解、錯誤訊息、日誌的撰寫規範,遵循 Go 社群最佳實踐。

## 📝 目錄

1. [註解規範 (Comments)](#註解規範)
2. [錯誤訊息規範 (Error Messages)](#錯誤訊息規範)
3. [日誌規範 (Logging)](#日誌規範)
4. [命名規範 (Naming)](#命名規範)
5. [實例對照](#實例對照)

---

## 註解規範

### 基本原則

1. **使用完整句子**
   - 以大寫字母開頭
   - 以句號結尾
   - 使用正確的標點符號

2. **註解應解釋「為什麼」而非「是什麼」**
   - ❌ 不好: `// Create user` (陳述顯而易見的事實)
   - ✅ 好: `// Create user with default staff role if no role specified.`

3. **公開函數/型別必須有文件註解**
   - 註解應緊接在聲明之前
   - 註解開頭使用聲明的名稱

### 範例

#### ✅ 正確的註解

```go
// NewUserService creates a new user service with required dependencies.
// It initializes repositories and logger for user management operations.
func NewUserService(
    logger infra.Logger,
    userRepository repositories.UserRepository,
) service.UserService {
    // ...
}

// GetUserByUsername retrieves a user and their permissions by username.
// Returns nil if user is not found.
func (s userService) GetUserByUsername(ctx context.Context, username string) (*domains.UserWithPermissions, error) {
    // ...
}

// Convert payload to JSON for JWT claims.
payloadJSON, err := json.Marshal(payload)

// Collect permissions from all roles.
for _, role := range user.Edges.Roles {
    // ...
}
```

#### ❌ 錯誤的註解

```go
// new user service
func NewUserService(...) { } // 小寫開頭,無句號

// get user
func GetUser(...) { } // 太簡短,無意義

// convert payload to json
payloadJson, err := json.Marshal(payload) // 小寫開頭

// Get permissions by role
for _, role := range user.Edges.Roles {
    // Add permissions to user
    // ... // 陳述顯而易見的事實
}
```

### 特殊註解類型

#### TODO 註解
```go
// TODO(username): Implement caching for frequently accessed users.
// TODO: Add retry logic for database connection failures.
```

#### FIXME 註解
```go
// FIXME: This calculation may overflow for large numbers.
```

#### 已棄用標記
```go
// Deprecated: Use GetUserByID instead.
func GetUser(id int) (*User, error) {
    // ...
}
```

---

## 錯誤訊息規範

### 基本原則

1. **使用小寫開頭** (除非是專有名詞)
2. **不使用標點符號結尾**
3. **簡潔明確,描述失敗的操作**
4. **動詞使用過去式或現在完成式**

### 格式規範

#### ✅ 推薦格式

```go
// 格式: "failed to <action> <object>"
errors.New("failed to connect to database")
errors.New("failed to parse JWT token")
errors.New("failed to create user")

// 格式: "invalid <object>"
errors.New("invalid username")
errors.New("invalid token format")

// 格式: "<object> not found"
errors.New("user not found")
errors.New("role not found")

// 格式: "unable to <action>"
errors.New("unable to establish connection")
```

#### ❌ 避免的格式

```go
errors.New("Error creating user")      // 不要用 "Error" 開頭
errors.New("Failed to create user.")   // 不要大寫開頭
errors.New("Cannot create user!")      // 不要用感嘆號
errors.New("Creating user failed")     // 不要用現在進行式
errors.New("User creation error")      // 太抽象
```

### 錯誤包裝

使用 `fmt.Errorf` 包裝錯誤時:

```go
// ✅ 正確
if err != nil {
    return fmt.Errorf("failed to query users: %w", err)
}

// ✅ 提供上下文
if err != nil {
    return fmt.Errorf("failed to create user %s: %w", username, err)
}

// ❌ 錯誤 - 沒有包裝原始錯誤
if err != nil {
    return fmt.Errorf("failed to query users: %v", err)
}
```

### 常見錯誤訊息範本

```go
// 資料庫操作
"failed to query users"
"failed to create user"
"failed to update user"
"failed to delete user"
"failed to connect to database"

// 驗證錯誤
"invalid username or password"
"invalid email format"
"username already exists"

// 授權錯誤
"user not authorized"
"insufficient permissions"
"token expired"

// 資源未找到
"user not found"
"resource not found"

// 解析錯誤
"failed to parse request body"
"failed to decode JSON"
"failed to unmarshal data"
```

---

## 日誌規範

### 基本原則

1. **使用小寫開頭**
2. **簡潔描述性強**
3. **包含適當的結構化欄位**
4. **選擇正確的日誌級別**

### 日誌級別使用指南

#### Debug
```go
logger.Debug("user authentication started",
    zap.String("username", username))
```
- 開發階段的詳細資訊
- 生產環境通常不記錄

#### Info
```go
logger.Info("database connected successfully")
logger.Info("starting HTTP server", zap.String("port", port))
logger.Info("graceful shutdown completed successfully")
```
- 正常運作的重要里程碑
- 系統啟動/關閉
- 重要業務流程完成

#### Warn
```go
logger.Warn("shutdown timeout exceeded")
logger.Warn("retry attempt failed",
    zap.Int("attempt", retryCount),
    zap.Error(err))
```
- 可恢復的錯誤
- 不影響主要功能的問題
- 需要注意但不緊急的情況

#### Error
```go
logger.Error("failed to connect to database", zap.Error(err))
logger.Error("failed to create user",
    zap.Error(err),
    zap.String("username", username))
```
- 操作失敗但系統繼續運行
- 需要調查的問題
- 影響單一請求/操作

#### Fatal
```go
logger.Fatal("failed to start HTTP server", zap.Error(err))
logger.Fatal("configuration file not found", zap.Error(err))
```
- 導致應用程式無法啟動/繼續的嚴重錯誤
- 會終止程式執行
- 謹慎使用

### 日誌訊息格式

#### ✅ 推薦格式

```go
// 操作完成
logger.Info("user created successfully",
    zap.Uint("user_id", userID))

// 操作失敗
logger.Error("failed to query users", zap.Error(err))

// 狀態變化
logger.Info("server shutting down gracefully")

// 帶上下文的錯誤
logger.Error("failed to authenticate user",
    zap.Error(err),
    zap.String("username", username))
```

#### ❌ 避免的格式

```go
logger.Info("User Created Successfully") // 大寫開頭
logger.Error("Error querying users", zap.Error(err)) // "Error" 開頭
logger.Info("Starting to run server...") // 太囉嗦
logger.Error("DB Error") // 太簡短,無意義
```

### 結構化欄位使用

#### 推薦的欄位命名

```go
// ✅ 使用 snake_case
zap.String("user_id", userID)
zap.String("trace_id", traceID)
zap.Duration("elapsed_time", duration)

// ❌ 避免 camelCase
zap.String("userId", userID)
zap.String("traceId", traceID)
```

#### 常用欄位

```go
// 識別資訊
zap.String("user_id", userID)
zap.String("username", username)
zap.String("request_id", requestID)
zap.String("trace_id", traceID)

// 時間相關
zap.Duration("elapsed_time", elapsed)
zap.Duration("timeout", timeout)
zap.Time("created_at", createdAt)

// 錯誤資訊
zap.Error(err)
zap.String("error_type", errorType)

// 系統資訊
zap.String("host", hostname)
zap.String("port", port)
zap.String("version", version)

// 業務資訊
zap.String("action", action)
zap.String("status", status)
zap.Int("count", count)
```

### 敏感資訊處理

```go
// ❌ 不要記錄敏感資訊
logger.Info("user login",
    zap.String("password", password)) // 危險!

// ✅ 記錄非敏感資訊
logger.Info("user login attempt",
    zap.String("username", username),
    zap.String("ip_address", ipAddr))

// ✅ 對敏感資訊脫敏
logger.Info("processing payment",
    zap.String("card_last4", lastFourDigits),
    zap.Float64("amount", amount))
```

---

## 命名規範

### 變數命名

#### ✅ 推薦

```go
// JSON 相關使用大寫縮寫
payloadJSON := marshal(data)
userJSON := encode(user)

// URL 使用大寫
requestURL := buildURL(path)

// ID 使用大寫
userID := getID()
traceID := generateID()

// 普通縮寫遵循 camelCase
ctx := context.Background()
err := doSomething()
cfg := loadConfig()
```

#### ❌ 避免

```go
payloadJson := marshal(data)  // json 應全大寫
userId := getID()              // ID 應全大寫
Url := buildURL(path)          // 不應大寫開頭(非導出)
```

### 常見縮寫規範

```go
// ✅ 正確的縮寫
API, HTTP, HTTPS, JSON, XML, HTML, URL, URI, ID, UUID, SQL
JWT, AWS, GCP, TCP, UDP, SSH, TLS, SSL

// 使用範例
apiKey := getKey()
httpClient := &http.Client{}
userID := 123
tokenUUID := uuid.New()
```

---

## 實例對照

### Service 層

#### Before (修正前)
```go
func (s userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
    users, err := s.userRepository.Client.User.Query().All(ctx)
    if err != nil {
        s.logger.WithContext(ctx).Error("Error getting users", zap.Error(err))
        return nil, err
    }
    return users, nil
}

func (s authService) GenerateToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
    // convert payload to json
    payloadJson, err := json.Marshal(payload)
    if err != nil {
        s.logger.WithContext(ctx).Error("marshalling payload failed", zap.Error(err))
        return "", err
    }
    // ...
}
```

#### After (修正後)
```go
func (s userService) GetUsers(ctx context.Context) ([]*entgen.User, error) {
    users, err := s.userRepository.Client.User.Query().All(ctx)
    if err != nil {
        s.logger.WithContext(ctx).Error("failed to query users", zap.Error(err))
        return nil, err
    }
    return users, nil
}

func (s authService) GenerateToken(ctx context.Context, user *domains.UserWithPermissions) (string, error) {
    // Convert payload to JSON for JWT claims.
    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        s.logger.WithContext(ctx).Error("failed to marshal token payload", zap.Error(err))
        return "", err
    }
    // ...
}
```

### Controller 層

#### Before (修正前)
```go
func (c ApprovalController) GetApprovals(w http.ResponseWriter, r *http.Request) {
    // check all permissions
    if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
        c.logger.WithContext(r.Context()).Error("Error user not authorized to get approvals")
        render.JSON(w, r, map[string]string{
            "error": "User not authorized to get approvals",
        })
        return
    }

    approvals, err := c.approvalService.GetApprovals(r.Context())
    if err != nil {
        c.logger.WithContext(r.Context()).Error("Error getting approvals", zap.Error(err))
        render.JSON(w, r, map[string]string{
            "error": "Error getting approvals",
        })
        return
    }
}
```

#### After (修正後)
```go
func (c ApprovalController) GetApprovals(w http.ResponseWriter, r *http.Request) {
    // Check if user has permission to read approvals.
    if hasPermission := permissions.Contains(constants.PermissionReadApproval); !hasPermission {
        c.logger.WithContext(r.Context()).Error("user not authorized to get approvals")
        render.JSON(w, r, map[string]string{
            "error": "user not authorized to get approvals",
        })
        return
    }

    approvals, err := c.approvalService.GetApprovals(r.Context())
    if err != nil {
        c.logger.WithContext(r.Context()).Error("failed to get approvals", zap.Error(err))
        render.JSON(w, r, map[string]string{
            "error": "failed to get approvals",
        })
        return
    }
}
```

### Infra 層

#### Before (修正前)
```go
func newDatabase(config Config, logger Logger) Database {
    client, err := ent.Open("sqlite3", config.SqliteDBPath)
    if err != nil {
        logger.Fatal("Error connecting to database", zap.Error(err))
    }
    logger.Info("Database connected")
    return Database{Client: client}
}
```

#### After (修正後)
```go
func newDatabase(config Config, logger Logger) Database {
    client, err := ent.Open("sqlite3", config.SqliteDBPath)
    if err != nil {
        logger.Fatal("failed to connect to database", zap.Error(err))
    }
    logger.Info("database connected successfully")
    return Database{Client: client}
}
```

---

## 快速檢查清單

### 註解檢查
- [ ] 大寫開頭
- [ ] 句號結尾
- [ ] 使用完整句子
- [ ] 公開 API 有文件註解
- [ ] 註解解釋「為什麼」而非「是什麼」

### 錯誤訊息檢查
- [ ] 小寫開頭(除非專有名詞)
- [ ] 無標點符號結尾
- [ ] 使用 "failed to" 或 "invalid" 等標準格式
- [ ] 簡潔明確
- [ ] 適當使用 `%w` 包裝錯誤

### 日誌檢查
- [ ] 小寫開頭
- [ ] 選擇正確的日誌級別
- [ ] 包含適當的結構化欄位
- [ ] 使用 snake_case 命名欄位
- [ ] 不記錄敏感資訊

### 命名檢查
- [ ] JSON/URL/ID 等縮寫全大寫
- [ ] 遵循 Go 的 camelCase 慣例
- [ ] 一致性的命名風格

---

## 參考資源

- [Effective Go - Commentary](https://go.dev/doc/effective_go#commentary)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Zap Logging Best Practices](https://pkg.go.dev/go.uber.org/zap)

---

## 維護說明

此文件應隨著專案演進持續更新。當發現新的模式或反模式時,請添加到相應章節。

最後更新: 2025-10-23
