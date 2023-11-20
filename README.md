see: https://tech.anti-pattern.co.jp/go-cache-gen/

## Install

```
go get github.com/yuemori/simplecache
```
## Usage

```golang
addr := os.GetEnv("REDIS_ADDRESS")
pass := os.GetEnv("REDIS_PASSWORD")
client := simplecache.NewRedisClient(addr, pass)

expiration := time.Minute
cache := cache.NewCache[*User](cacheClient, expiration)

dsn := os.GetEnv("DSN")

db, err := sqlx.Open("mysql", dsn)
if err != nil {
  panic(err)
}

ctx := context.Background()
id := 1

cache.GetOrSet(ctx, fmt.Sprintf("user:id:%d", id), func(ctx context.Context) (*User, error) {
    var u *User

    if err := db.GetContext(ctx, "SELECT * FROM users WHERE id = ?", &u, id); err != nil {
      return nil, err
    }

    return u, nil
  },
)
```
