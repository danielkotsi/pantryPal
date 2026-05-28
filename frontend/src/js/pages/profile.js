/**
 * Profile Page Handler
 * Handles profile page loading, editing, and form submissions
 */

class ProfilePageHandler {
    constructor() {
        this.currentProfile = null;
        this.setupEventListeners();
    }

    /**
     * Setup event listeners for profile forms
     */
    setupEventListeners() {
        document.addEventListener('submit', async (e) => {
            if (e.target.id === 'metricsForm') {
                e.preventDefault();
                await this.handleMetricsSubmit(e.target);
            }
            if (e.target.id === 'preferencesForm') {
                e.preventDefault();
                await this.handlePreferencesSubmit(e.target);
            }
            if (e.target.id === 'budgetForm') {
                e.preventDefault();
                await this.handleBudgetSubmit(e.target);
            }
        });
    }

    /**
     * Load and display profile
     */
    async loadProfile() {
        try {
            this.currentProfile = await api.getProfile();
            this.renderProfile();
        } catch (error) {
            router.setState({ error: 'Failed to load profile' });
        }
    }

    /**
     * Render profile data
     */
    renderProfile() {
        if (!this.currentProfile) return;

        this.renderPersonalInfo();
        this.renderMetrics();
        this.renderPreferences();
        this.renderBudget();
    }

    /**
     * Render personal information section
     */
    renderPersonalInfo() {
        const container = document.getElementById('personalInfo');
        if (!container) return;

        const user = this.currentProfile.user;
        container.innerHTML = `
            <div class="info-item">
                <label>Email</label>
                <span>${this.escapeHtml(user.email)}</span>
            </div>
            <div class="info-item">
                <label>Display Name</label>
                <span>${this.escapeHtml(user.displayName || 'Not set')}</span>
            </div>
        `;
    }

    /**
     * Render body metrics form
     */
    renderMetrics() {
        const container = document.getElementById('bodyMetrics');
        if (!container) return;

        const metrics = this.currentProfile.metrics || {};
        container.innerHTML = `
            <div class="form-group">
                <label for="heightCm">Height (cm)</label>
                <input type="number" id="heightCm" step="0.1" value="${metrics.heightCm || ''}">
            </div>
            <div class="form-group">
                <label for="weightKg">Weight (kg)</label>
                <input type="number" id="weightKg" step="0.1" value="${metrics.weightKg || ''}">
            </div>
            <div class="form-group">
                <label for="age">Age</label>
                <input type="number" id="age" min="1" value="${metrics.age || ''}">
            </div>
            <div class="form-group">
                <label for="sex">Sex</label>
                <select id="sex">
                    <option value="">Select...</option>
                    <option value="male" ${metrics.sex === 'M' ? 'selected' : ''}>Male</option>
                    <option value="female" ${metrics.sex === 'F' ? 'selected' : ''}>Female</option>
                    <option value="other" ${metrics.sex === 'O' ? 'selected' : ''}>Other</option>
                </select>
            </div>
            <div class="form-group">
                <label for="activityLevel">Activity Level</label>
                <select id="activityLevel">
                    <option value="">Select...</option>
                    <option value="sedentary" ${metrics.activityLevel === 'sedentary' ? 'selected' : ''}>Sedentary</option>
                    <option value="light" ${metrics.activityLevel === 'light' ? 'selected' : ''}>Light</option>
                    <option value="moderate" ${metrics.activityLevel === 'moderate' ? 'selected' : ''}>Moderate</option>
                    <option value="very_active" ${metrics.activityLevel === 'very_active' ? 'selected' : ''}>Very Active</option>
                </select>
            </div>
            <div class="form-group">
                <label for="goal">Goal</label>
                <select id="goal">
                    <option value="">Select...</option>
                    <option value="lose_weight" ${metrics.goal === 'lose_weight' ? 'selected' : ''}>Lose Weight</option>
                    <option value="maintain" ${metrics.goal === 'maintain' ? 'selected' : ''}>Maintain</option>
                    <option value="gain_muscle" ${metrics.goal === 'gain_muscle' ? 'selected' : ''}>Gain Muscle</option>
                </select>
            </div>
        `;
    }

    /**
     * Render preferences form
     */
    renderPreferences() {
        const container = document.getElementById('preferences');
        if (!container) return;

        const prefs = this.currentProfile.preferences || {};
        container.innerHTML = `
            <div class="form-group">
                <label for="dietType">Diet Type</label>
                <select id="dietType">
                    <option value="">Select...</option>
                    <option value="omnivore" ${prefs.dietType === 'omnivore' ? 'selected' : ''}>Omnivore</option>
                    <option value="vegetarian" ${prefs.dietType === 'vegetarian' ? 'selected' : ''}>Vegetarian</option>
                    <option value="vegan" ${prefs.dietType === 'vegan' ? 'selected' : ''}>Vegan</option>
                    <option value="pescatarian" ${prefs.dietType === 'pescatarian' ? 'selected' : ''}>Pescatarian</option>
                </select>
            </div>
            <div class="form-group">
                <label for="allergies">Allergies (comma-separated)</label>
                <input type="text" id="allergies" value="${(prefs.allergies || []).join(', ')}">
            </div>
            <div class="form-group">
                <label for="dislikes">Dislikes (comma-separated)</label>
                <input type="text" id="dislikes" value="${(prefs.dislikes || []).join(', ')}">
            </div>
            <div class="form-group">
                <label for="likes">Likes (comma-separated)</label>
                <input type="text" id="likes" value="${(prefs.likes || []).join(', ')}">
            </div>
            <div class="form-group">
                <label for="dailyCalorieTarget">Daily Calorie Target</label>
                <input type="number" id="dailyCalorieTarget" value="${prefs.dailyCalorieTarget || ''}">
            </div>
            <div class="form-group">
                <label for="notes">Notes</label>
                <textarea id="notes" rows="3">${prefs.notes || ''}</textarea>
            </div>
        `;
    }

    /**
     * Render budget form
     */
    renderBudget() {
        const container = document.getElementById('budget');
        if (!container) return;

        const budget = this.currentProfile.budget || {};
        const currentMonth = new Date().toISOString().substring(0, 7); // YYYY-MM

        container.innerHTML = `
            <div class="form-group">
                <label for="month">Month</label>
                <input type="month" id="month" value="${budget.month || currentMonth}">
            </div>
            <div class="form-group">
                <label for="currency">Currency</label>
                <input type="text" id="currency" maxlength="3" value="${budget.currency || 'USD'}">
            </div>
            <div class="form-group">
                <label for="amountCents">Amount (in cents)</label>
                <input type="number" id="amountCents" min="0" value="${budget.amountCents || ''}">
            </div>
        `;
    }

    /**
     * Handle metrics form submission
     */
    async handleMetricsSubmit(form) {
        try {
            const metrics = {
                heightCm: this.getInputValue(form, 'heightCm', 'float'),
                weightKg: this.getInputValue(form, 'weightKg', 'float'),
                age: this.getInputValue(form, 'age', 'int'),
                sex: form.querySelector('#sex').value || null,
                activityLevel: form.querySelector('#activityLevel').value || null,
                goal: form.querySelector('#goal').value || null,
            };

            // Remove null values
            Object.keys(metrics).forEach(key => {
                if (metrics[key] === null || metrics[key] === '') {
                    delete metrics[key];
                }
            });

            await api.updateMetrics(metrics);
            router.setState({ error: null });
            this.loadProfile();
        } catch (error) {
            router.setState({ error: error.message });
        }
    }

    /**
     * Handle preferences form submission
     */
    async handlePreferencesSubmit(form) {
        try {
            const preferences = {
                dietType: form.querySelector('#dietType').value || null,
                allergies: this.parseCommaList(form.querySelector('#allergies').value),
                dislikes: this.parseCommaList(form.querySelector('#dislikes').value),
                likes: this.parseCommaList(form.querySelector('#likes').value),
                dailyCalorieTarget: this.getInputValue(form, 'dailyCalorieTarget', 'int'),
                notes: form.querySelector('#notes').value || null,
            };

            // Remove null and empty values
            Object.keys(preferences).forEach(key => {
                if (preferences[key] === null || preferences[key] === '' || 
                    (Array.isArray(preferences[key]) && preferences[key].length === 0)) {
                    delete preferences[key];
                }
            });

            await api.updatePreferences(preferences);
            router.setState({ error: null });
            this.loadProfile();
        } catch (error) {
            router.setState({ error: error.message });
        }
    }

    /**
     * Handle budget form submission
     */
    async handleBudgetSubmit(form) {
        try {
            const budget = {
                month: form.querySelector('#month').value || null,
                currency: form.querySelector('#currency').value || null,
                amountCents: this.getInputValue(form, 'amountCents', 'int'),
            };

            // Remove null values
            Object.keys(budget).forEach(key => {
                if (budget[key] === null || budget[key] === '') {
                    delete budget[key];
                }
            });

            await api.updateBudget(budget);
            router.setState({ error: null });
            this.loadProfile();
        } catch (error) {
            router.setState({ error: error.message });
        }
    }

    /**
     * Get input value with type conversion
     */
    getInputValue(form, id, type) {
        const value = form.querySelector(`#${id}`).value;
        if (!value) return null;

        if (type === 'int') return parseInt(value);
        if (type === 'float') return parseFloat(value);
        return value;
    }

    /**
     * Parse comma-separated list
     */
    parseCommaList(str) {
        if (!str) return [];
        return str.split(',').map(item => item.trim()).filter(item => item.length > 0);
    }

    /**
     * Escape HTML to prevent XSS
     */
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Create global profile handler instance
const profilePageHandler = new ProfilePageHandler();
