see: https://tech.anti-pattern.co.jp/go-cache-gen/

## Install

```
go get github.com/yuemori/simplecache
```

## Usage

```
type userRepository struct {
	db    *sql.DB
	cache *cache.Cache[*user.User]
}

func NewUserRepository(
	db *sql.DB,
	cacheClient cache.Client,
) user.Repository {
	return &userRepository{
		db: db,
		cache: cache.NewCache[*user.User](
			cacheClient,
			time.Minute,
		),
	}
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id int,
) (*user.User, error) {
	return r.cache.GetOrSet(ctx, makeUserKey(id),
		func(ctx context.Context) (*user.User, error) {
			var u *user.User
			// r.dbを使ってDBから取得する
			return u, nil
		},
	)
}

func makeUserKey(id int) string {
	return fmt.Sprintf("user:id:%d", id)
}
```
