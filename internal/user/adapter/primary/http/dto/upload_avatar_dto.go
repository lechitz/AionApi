package dto

// UploadAvatarResponse is returned after successful avatar upload.
type UploadAvatarResponse struct {
	AvatarURL   string `json:"avatar_url"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
}
