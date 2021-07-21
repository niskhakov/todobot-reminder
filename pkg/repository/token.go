package repository

type Bucket string

const (
	AccessTokens Bucket = "access_tokens"
	CodeTokens   Bucket = "code_tokens"
)

type IterateFunc func(chatID int64, accessToken string, accumulator interface{}) error

type TokenRepository interface {
	Save(chatID int64, token string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
	ForEach(bucket Bucket, fnc IterateFunc, accumulator interface{}) error
}
