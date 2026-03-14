package output

import "context"

// AvatarUploadInput defines avatar object payload sent to storage adapters.
type AvatarUploadInput struct {
	ObjectKey   string
	ContentType string
	Body        []byte
}

// AvatarStorage defines storage operations for user avatar uploads.
type AvatarStorage interface {
	UploadAvatar(ctx context.Context, input AvatarUploadInput) (string, error)
}
