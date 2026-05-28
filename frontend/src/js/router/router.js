/**
 * Router & State Management
 * Handles routing between pages and manages app state
 */

class Router {
    constructor() {
        this.currentRoute = 'auth';
        this.routes = new Map();
        this.state = {
            user: null,
            isAuthenticated: false,
            loading: false,
            error: null,
        };
        this.listeners = [];
        this.pageCache = new Map();

        this.initializeRoutes();
        this.setupEventListeners();
        this.restoreState();
    }

    /**
     * Register a route
     */
    register(name, handler) {
        this.routes.set(name, handler);
    }

    /**
     * Subscribe to state changes
     */
    subscribe(listener) {
        this.listeners.push(listener);
        return () => {
            this.listeners = this.listeners.filter(l => l !== listener);
        };
    }

    /**
     * Update state and notify listeners
     */
    setState(updates) {
        this.state = { ...this.state, ...updates };
        this.listeners.forEach(listener => listener(this.state));
        this.persistState();
    }

    /**
     * Get current state
     */
    getState() {
        return this.state;
    }

    /**
     * Navigate to a route
     */
    navigate(routeName) {
        if (!this.routes.has(routeName)) {
            console.error(`Route not found: ${routeName}`);
            return;
        }

        this.currentRoute = routeName;
        this.render();
    }

    /**
     * Initialize built-in routes
     */
    initializeRoutes() {
        this.register('auth', () => this.renderPage('authPage'));
        this.register('profile', () => {
            if (!this.state.isAuthenticated) {
                this.navigate('auth');
                return;
            }
            this.renderPage('profilePage');
            // Load profile data after page is rendered
            setTimeout(() => profilePageHandler.loadProfile(), 0);
        });
        this.register('planner', () => {
            if (!this.state.isAuthenticated) {
                this.navigate('auth');
                return;
            }
            this.renderPage('plannerPage');
            setTimeout(() => plannerPageHandler.init(), 0);
        });
        this.register('pantry', () => {
            if (!this.state.isAuthenticated) {
                this.navigate('auth');
                return;
            }
            this.renderPage('pantryPage');
        });
        this.register('chat', () => {
            if (!this.state.isAuthenticated) {
                this.navigate('auth');
                return;
            }
            this.renderPage('chatPage');
            setTimeout(() => chatPageHandler.init(), 0);
        });
    }

    /**
     * Setup navigation event listeners
     */
    setupEventListeners() {
        // Nav link clicks
        document.addEventListener('click', (e) => {
            const navLink = e.target.closest('[data-route]');
            if (navLink) {
                e.preventDefault();
                const route = navLink.dataset.route;
                this.navigate(route);
            }
        });

        // Error container close button
        document.addEventListener('click', (e) => {
            if (e.target.closest('.close-error')) {
                this.setState({ error: null });
            }
        });

        // API error handler
        api.onError((error) => {
            this.setState({ error: error.message });
        });

        // API loading handler
        api.onLoadingChange((isLoading) => {
            this.setState({ loading: isLoading });
        });

        // Auth form handler
        document.addEventListener('submit', async (e) => {
            if (e.target.id === 'authForm') {
                e.preventDefault();
                await this.handleAuth(e.target);
            }

        });

        // Sign up toggle
        document.addEventListener('click', (e) => {
            if (e.target.id === 'toggleSignup') {
                this.toggleAuthMode();
            }
        });
    }

    /**
     * Render page from template
     */
    renderPage(templateId) {
        const pageContent = document.getElementById('pageContent');
        const template = document.getElementById(templateId);

        if (!template) {
            console.error(`Template not found: ${templateId}`);
            return;
        }

        const clone = template.content.cloneNode(true);
        pageContent.innerHTML = '';
        pageContent.appendChild(clone);

        this.updateNavigation();
        this.updateErrorDisplay();
        this.updateLoadingDisplay();
        this.updateUserDisplay();
    }

    /**
     * Render the current route
     */
    render() {
        const handler = this.routes.get(this.currentRoute);
        if (handler) {
            handler();
        }
    }

    /**
     * Update navigation active state
     */
    updateNavigation() {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
            if (link.dataset.route === this.currentRoute) {
                link.classList.add('active');
            }
        });
    }

    /**
     * Update error display
     */
    updateErrorDisplay() {
        const errorContainer = document.getElementById('errorContainer');
        const errorText = document.getElementById('errorText');

        if (this.state.error) {
            errorText.textContent = this.state.error;
            errorContainer.style.display = 'block';
        } else {
            errorContainer.style.display = 'none';
        }
    }

    /**
     * Update loading display
     */
    updateLoadingDisplay() {
        const loadingContainer = document.getElementById('loadingContainer');
        if (this.state.loading) {
            loadingContainer.style.display = 'flex';
        } else {
            loadingContainer.style.display = 'none';
        }
    }

    /**
     * Update user display in navbar
     */
    updateUserDisplay() {
        const navbarUser = document.getElementById('navbarUser');
        if (this.state.isAuthenticated && this.state.user) {
            navbarUser.innerHTML = `
                <div class="user-info">
                    <span>${this.state.user.displayName || this.state.user.email}</span>
                    <button class="btn btn-secondary" id="logoutBtn">Logout</button>
                </div>
            `;

            document.getElementById('logoutBtn')?.addEventListener('click', () => {
                this.logout();
            });
        } else {
            navbarUser.innerHTML = '';
        }
    }

    /**
     * Handle authentication (login/signup)
     */
    async handleAuth(form) {
        const email = form.querySelector('#email').value;
        const password = form.querySelector('#password').value;
        const isSignup = form.classList.contains('signup-mode');

        try {
            let result;
            if (isSignup) {
                const displayName = form.querySelector('#displayName')?.value || email.split('@')[0];
                result = await api.signup(email, password, displayName);
            } else {
                result = await api.login(email, password);
            }

            this.setState({
                user: result.user,
                isAuthenticated: true,
                error: null,
            });

            this.navigate('profile');
        } catch (error) {
            // Extract error message from API error structure
            let errorMessage = error.message;
            if (error.data?.error?.message) {
                errorMessage = error.data.error.message;
            }
            this.setState({ error: errorMessage });
        }
    }

    /**
     * Toggle between login and signup modes
     */
    toggleAuthMode() {
        const form = document.getElementById('authForm');
        const isSignup = form.classList.contains('signup-mode');

        if (!isSignup) {
            // Switch to signup
            form.classList.add('signup-mode');
            form.innerHTML = `
                <div class="form-group">
                    <label for="displayName">Display Name:</label>
                    <input type="text" id="displayName" name="displayName" required>
                </div>
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" name="email" required>
                </div>
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary">Sign Up</button>
                    <button type="button" class="btn btn-secondary" id="toggleSignup">Back to Login</button>
                </div>
            `;
        } else {
            // Switch to login
            form.classList.remove('signup-mode');
            form.innerHTML = `
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" id="email" name="email" required>
                </div>
                <div class="form-group">
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary">Login</button>
                    <button type="button" class="btn btn-secondary" id="toggleSignup">Sign Up</button>
                </div>
            `;
        }
    }

    /**
     * Logout user
     */
    async logout() {
        await api.logout();
        this.setState({
            user: null,
            isAuthenticated: false,
            error: null,
        });
        this.navigate('auth');
    }

    /**
     * Persist state to localStorage
     */
    persistState() {
        const stateToSave = {
            user: this.state.user,
            isAuthenticated: this.state.isAuthenticated,
        };
        localStorage.setItem('appState', JSON.stringify(stateToSave));
    }

    /**
     * Restore state from localStorage
     */
    restoreState() {
        const saved = localStorage.getItem('appState');
        if (saved) {
            try {
                const { user, isAuthenticated } = JSON.parse(saved);
                this.setState({ user, isAuthenticated });
            } catch (e) {
                console.error('Failed to restore state:', e);
            }
        }
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

// Create global router instance
const router = new Router();
