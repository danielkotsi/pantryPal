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

type FoodSearchItem struct {
	FDCID       int64  `json:"fdcId"`
	Description string `json:"description"`
	FoodClass   string `json:"foodClass,omitempty"`
}

type FoodSearchResponse struct {
	Items []FoodSearchItem `json:"items"`
}

type PantryItemRequest struct {
	FDCID    int64   `json:"fdcId"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type PantryItemPatchRequest struct {
	QuantityDelta float64 `json:"quantityDelta"`
}

type PantryFood struct {
	FDCID       int64  `json:"fdcId"`
	Description string `json:"description"`
	FoodClass   string `json:"foodClass,omitempty"`
}

type PantryItemResponse struct {
	ID       string     `json:"id"`
	Quantity float64    `json:"quantity"`
	Unit     string     `json:"unit"`
	Food     PantryFood `json:"food"`
}

type PantryItemsResponse struct {
	Items []PantryItemResponse `json:"items"`
}

type RecipeIngredientResponse struct {
	FDCID       int64   `json:"fdcId"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

type RecipeMacrosResponse struct {
	Calories float64 `json:"calories"`
	ProteinG float64 `json:"proteinG"`
	CarbsG   float64 `json:"carbsG"`
	FatG     float64 `json:"fatG"`
}

type RecipeResponse struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	MealType           string                     `json:"mealType"`
	Servings           int                        `json:"servings"`
	Instructions       string                     `json:"instructions,omitempty"`
	EstimatedCostCents int                        `json:"estimatedCostCents"`
	Macros             RecipeMacrosResponse       `json:"macros"`
	Ingredients        []RecipeIngredientResponse `json:"ingredients"`
}

type PlanMealInput struct {
	ScheduledDate      string               `json:"scheduledDate"`
	MealSection        string               `json:"mealSection"`
	RecipeID           string               `json:"recipeId,omitempty"`
	RecipeName         string               `json:"recipeName"`
	Servings           float64              `json:"servings"`
	EstimatedCostCents int                  `json:"estimatedCostCents"`
	Macros             RecipeMacrosResponse `json:"macros"`
}

type PlanProposalRequest struct {
	PeriodType       string          `json:"periodType"`
	StartDate        string          `json:"startDate"`
	EndDate          string          `json:"endDate,omitempty"`
	Source           string          `json:"source,omitempty"`
	ProposalVersion  int             `json:"proposalVersion,omitempty"`
	AICostCentsTotal int             `json:"aiCostCentsTotal,omitempty"`
	Notes            string          `json:"notes,omitempty"`
	Meals            []PlanMealInput `json:"meals"`
}

type DeclinePlanRequest struct {
	Reason string `json:"reason,omitempty"`
}

type PlanSummaryResponse struct {
	ID               string `json:"id"`
	PeriodType       string `json:"periodType"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	Status           string `json:"status"`
	Source           string `json:"source"`
	ProposalVersion  int    `json:"proposalVersion"`
	AICostCentsTotal int    `json:"aiCostCentsTotal"`
	Notes            string `json:"notes,omitempty"`
}

type PlanMealResponse struct {
	ID                 string               `json:"id"`
	RecipeID           string               `json:"recipeId,omitempty"`
	RecipeName         string               `json:"recipeName"`
	ScheduledDate      string               `json:"scheduledDate"`
	MealSection        string               `json:"mealSection"`
	Servings           float64              `json:"servings"`
	EstimatedCostCents int                  `json:"estimatedCostCents"`
	Macros             RecipeMacrosResponse `json:"macros"`
	IsConsumed         bool                 `json:"isConsumed"`
	ConsumedAt         string               `json:"consumedAt,omitempty"`
}

type PlanDaySections struct {
	Breakfast *PlanMealResponse `json:"breakfast,omitempty"`
	Lunch     *PlanMealResponse `json:"lunch,omitempty"`
	Dinner    *PlanMealResponse `json:"dinner,omitempty"`
	Snacks    *PlanMealResponse `json:"snacks,omitempty"`
}

type PlanDayResponse struct {
	Date     string               `json:"date"`
	Sections PlanDaySections      `json:"sections"`
	Totals   RecipeMacrosResponse `json:"totals"`
}

type ProposalResponse struct {
	Plan       PlanSummaryResponse  `json:"plan"`
	Days       []PlanDayResponse    `json:"days"`
	WeekTotals RecipeMacrosResponse `json:"weekTotals"`
}

type WeekPlanResponse struct {
	Plans      []PlanSummaryResponse `json:"plans"`
	Days       []PlanDayResponse     `json:"days"`
	WeekTotals RecipeMacrosResponse  `json:"weekTotals"`
}

type ChatSendRequest struct {
	Message string `json:"message"`
	Action  string `json:"action,omitempty"`
}

type ChatMessageResponse struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Action    string `json:"action,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type ChatHistoryResponse struct {
	Messages []ChatMessageResponse `json:"messages"`
}

type GeneratePlanRequest struct {
	PeriodType string `json:"periodType"`
	Message    string `json:"message,omitempty"`
}

type GeneratePlanResponse struct {
	Proposal       ProposalResponse `json:"proposal"`
	FallbackActive bool             `json:"fallbackActive"`
}

type ConsumedItemInfo struct {
	FDCID            int64   `json:"fdcId"`
	Description      string  `json:"description"`
	PantryItemID     string  `json:"pantryItemId,omitempty"`
	QuantityDeducted float64 `json:"quantityDeducted"`
	Unit             string  `json:"unit"`
	BeforeQuantity   float64 `json:"beforeQuantity"`
	AfterQuantity    float64 `json:"afterQuantity"`
	Warning          string  `json:"warning,omitempty"`
}

type ConsumeMealResponse struct {
	MealID    string             `json:"mealId"`
	Consumed  bool               `json:"consumed"`
	ConsumedAt string            `json:"consumedAt"`
	Items     []ConsumedItemInfo `json:"items"`
}
