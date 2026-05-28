/**
 * API Client
 * Handles all HTTP requests with auth headers, JSON handling, and error management
 */

class APIClient {
    constructor(baseURL = 'http://localhost:8080') {
        this.baseURL = baseURL;
        this.token = this.getStoredToken();
        this.errorCallbacks = [];
        this.loadingCallbacks = [];
    }

    /**
     * Set auth token
     */
    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('authToken', token);
        } else {
            localStorage.removeItem('authToken');
        }
    }

    /**
     * Get stored token from localStorage
     */
    getStoredToken() {
        return localStorage.getItem('authToken');
    }

    /**
     * Check if user is authenticated
     */
    isAuthenticated() {
        return !!this.token;
    }

    /**
     * Build headers with auth token
     */
    buildHeaders(options = {}) {
        const headers = {
            'Content-Type': 'application/json',
            ...options,
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        return headers;
    }

    /**
     * Register error callback
     */
    onError(callback) {
        this.errorCallbacks.push(callback);
    }

    /**
     * Trigger error callbacks
     */
    triggerError(error) {
        this.errorCallbacks.forEach(cb => cb(error));
    }

    /**
     * Register loading callback
     */
    onLoadingChange(callback) {
        this.loadingCallbacks.push(callback);
    }

    /**
     * Trigger loading callbacks
     */
    triggerLoadingChange(isLoading) {
        this.loadingCallbacks.forEach(cb => cb(isLoading));
    }

    /**
     * Make HTTP request
     */
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const method = options.method || 'GET';
        const headers = this.buildHeaders(options.headers);

        const config = {
            method,
            headers,
            ...options,
        };

        // Add body for methods that support it
        if (options.body && typeof options.body === 'object') {
            config.body = JSON.stringify(options.body);
        }

        try {
            this.triggerLoadingChange(true);

            const response = await fetch(url, config);

            // Handle response
            let data = null;
            const contentType = response.headers.get('content-type');

            if (contentType && contentType.includes('application/json')) {
                data = await response.json();
            } else {
                data = await response.text();
            }

            // Handle HTTP errors
            if (!response.ok) {
                let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
                
                // Extract error message from backend response format
                if (data?.error?.message) {
                    errorMessage = data.error.message;
                } else if (data?.message) {
                    errorMessage = data.message;
                }

                const error = new APIError(
                    errorMessage,
                    response.status,
                    data
                );
                this.triggerError(error);
                throw error;
            }

            this.triggerLoadingChange(false);
            return data;
        } catch (error) {
            this.triggerLoadingChange(false);

            if (error instanceof APIError) {
                throw error;
            }

            const apiError = new APIError(
                error.message || 'Network error',
                null,
                error
            );
            this.triggerError(apiError);
            throw apiError;
        }
    }

    /**
     * GET request
     */
    get(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'GET' });
    }

    /**
     * POST request
     */
    post(endpoint, body, options = {}) {
        return this.request(endpoint, { ...options, method: 'POST', body });
    }

    /**
     * PUT request
     */
    put(endpoint, body, options = {}) {
        return this.request(endpoint, { ...options, method: 'PUT', body });
    }

    /**
     * PATCH request
     */
    patch(endpoint, body, options = {}) {
        return this.request(endpoint, { ...options, method: 'PATCH', body });
    }

    /**
     * DELETE request
     */
    delete(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'DELETE' });
    }

    /**
     * Auth endpoints
     */
    async login(email, password) {
        const data = await this.post('/auth/login', { email, password });
        if (data.token) {
            this.setToken(data.token);
        }
        return data;
    }

    async signup(email, password, displayName) {
        const data = await this.post('/auth/register', { email, password, displayName });
        if (data.token) {
            this.setToken(data.token);
        }
        return data;
    }

    async getMe() {
        const data = await this.get('/me');
        return data.user;
    }

    async logout() {
        this.setToken(null);
        return { success: true };
    }

    /**
     * Profile endpoints
     */
    getProfile() {
        return this.get('/profile');
    }

    updateMetrics(metrics) {
        return this.patch('/profile/metrics', metrics);
    }

    updatePreferences(preferences) {
        return this.patch('/profile/preferences', preferences);
    }

    updateBudget(budget) {
        return this.patch('/profile/budget', budget);
    }

    /**
     * Pantry endpoints
     */
    getPantryItems() {
        return this.get('/pantry/items');
    }

    addPantryItem(item) {
        return this.post('/pantry/items', item);
    }

    updatePantryItem(itemId, quantityDelta) {
        return this.patch(`/pantry/items/${itemId}`, { quantityDelta });
    }

    deletePantryItem(itemId) {
        return this.delete(`/pantry/items/${itemId}`);
    }

    searchIngredients(query) {
        return this.get(`/ingredients/search?q=${encodeURIComponent(query)}`);
    }

    /**
     * Meal plan endpoints
     */
    getWeekPlan(startDate) {
        return this.get(`/plans/week?start=${startDate}`);
    }

    createProposal(proposal) {
        return this.post('/plans/proposal', proposal);
    }

    acceptProposal(planId) {
        return this.post(`/plans/${planId}/accept`);
    }

    declineProposal(planId, reason) {
        return this.post(`/plans/${planId}/decline`, { reason });
    }

    /**
     * Recipe endpoints
     */
    getRecipe(recipeId) {
        return this.get(`/recipes/${recipeId}`);
    }

    /**
     * Chat endpoints
     */
    sendChatMessage(message, action) {
        const body = { message };
        if (action) body.action = action;
        return this.post('/chat', body);
    }

    getChatHistory() {
        return this.get('/chat');
    }

    generatePlan(periodType, message) {
        const body = { periodType };
        if (message) body.message = message;
        return this.post('/plans/generate', body);
    }

    /**
     * Consumption log endpoints
     */
    getConsumptionLog(date) {
        const params = date ? `?date=${date}` : '';
        return this.get(`/consumption-log${params}`);
    }

    logConsumption(consumption) {
        return this.post('/consumption-log', consumption);
    }

    /**
     * Budget endpoints
     */
    getBudgets() {
        return this.get('/budgets');
    }

    setBudget(budget) {
        return this.post('/budgets', budget);
    }
}

/**
 * Custom API Error class
 */
class APIError extends Error {
    constructor(message, status, data) {
        super(message);
        this.name = 'APIError';
        this.status = status;
        this.data = data;
    }
}

// Create global API client instance
const api = new APIClient();
