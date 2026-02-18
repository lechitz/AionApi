package dto

// StartRegistrationRequest represents registration step 1 payload.
type StartRegistrationRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateRegistrationProfileRequest represents registration step 2 payload.
type UpdateRegistrationProfileRequest struct {
	Locale   string `json:"locale"`
	Timezone string `json:"timezone"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
}

// UpdateRegistrationAvatarRequest represents registration step 3 payload.
type UpdateRegistrationAvatarRequest struct {
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// RegistrationSessionResponse represents staged registration state.
type RegistrationSessionResponse struct {
	RegistrationID string  `json:"registration_id"`
	CurrentStep    int     `json:"current_step"`
	Name           string  `json:"name"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Locale         *string `json:"locale,omitempty"`
	Timezone       *string `json:"timezone,omitempty"`
	Location       *string `json:"location,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	ExpiresAt      string  `json:"expires_at"`
	Status         string  `json:"status"`
}
