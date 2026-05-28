package dto

type APIError struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expiresAt"`
	User      UserResponse `json:"user"`
}

type ProfileResponse struct {
	User        UserResponse      `json:"user"`
	Metrics     ProfileMetricsOut `json:"metrics"`
	Preferences ProfilePrefsOut   `json:"preferences"`
	Budget      ProfileBudgetOut  `json:"budget"`
}

type ProfileMetricsOut struct {
	HeightCM      *float64 `json:"heightCm"`
	WeightKG      *float64 `json:"weightKg"`
	Age           *int     `json:"age"`
	Sex           *string  `json:"sex"`
	ActivityLevel *string  `json:"activityLevel"`
	Goal          *string  `json:"goal"`
}

type ProfilePrefsOut struct {
	DietType           string   `json:"dietType"`
	Allergies          []string `json:"allergies"`
	Dislikes           []string `json:"dislikes"`
	Likes              []string `json:"likes"`
	DailyCalorieTarget *int     `json:"dailyCalorieTarget"`
	Notes              *string  `json:"notes"`
}

type ProfileBudgetOut struct {
	Month       string `json:"month"`
	Currency    string `json:"currency"`
	AmountCents int    `json:"amountCents"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PatchMetricsRequest struct {
	HeightCM      *float64 `json:"heightCm"`
	WeightKG      *float64 `json:"weightKg"`
	Age           *int     `json:"age"`
	Sex           *string  `json:"sex"`
	ActivityLevel *string  `json:"activityLevel"`
	Goal          *string  `json:"goal"`
}

type PatchPreferencesRequest struct {
	DietType           *string  `json:"dietType"`
	Allergies          []string `json:"allergies"`
	Dislikes           []string `json:"dislikes"`
	Likes              []string `json:"likes"`
	DailyCalorieTarget *int     `json:"dailyCalorieTarget"`
	Notes              *string  `json:"notes"`
}

type PatchBudgetRequest struct {
	Month       *string `json:"month"`
	Currency    *string `json:"currency"`
	AmountCents *int    `json:"amountCents"`
}
