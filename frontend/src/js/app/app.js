/**
 * Main App
 * Initializes the application
 */

class App {
    constructor() {
        this.init();
    }

    /**
     * Initialize the app
     */
    init() {
        // Subscribe to router state changes
        router.subscribe((state) => {
            this.onStateChange(state);
        });

        // Initial render
        router.render();

        console.log('PantryPal app initialized');
    }

    /**
     * Handle state changes
     */
    onStateChange(state) {
        // Update UI based on state changes
        if (state.error) {
            this.showError(state.error);
        }

        if (state.loading) {
            this.showLoading();
        } else {
            this.hideLoading();
        }
    }

    /**
     * Show error message
     */
    showError(message) {
        const errorContainer = document.getElementById('errorContainer');
        const errorText = document.getElementById('errorText');
        errorText.textContent = message;
        errorContainer.style.display = 'block';

        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            if (errorContainer.style.display !== 'none') {
                router.setState({ error: null });
            }
        }, 5000);
    }

    /**
     * Show loading indicator
     */
    showLoading() {
        const loadingContainer = document.getElementById('loadingContainer');
        loadingContainer.style.display = 'flex';
    }

    /**
     * Hide loading indicator
     */
    hideLoading() {
        const loadingContainer = document.getElementById('loadingContainer');
        loadingContainer.style.display = 'none';
    }
}

// Initialize app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        new App();
    });
} else {
    new App();
}
