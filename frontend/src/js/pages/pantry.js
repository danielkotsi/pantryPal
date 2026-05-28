class PantryPageHandler {
    constructor() {
        this.items = [];
        this.searchResults = [];
        this.selectedIngredient = null;
        this.searchTimeout = null;
        this.setupEventListeners();
    }

    setupEventListeners() {
        document.addEventListener('input', (e) => {
            if (e.target.id === 'searchPantry') {
                this.handleSearchInput(e.target.value);
            }
        });

        document.addEventListener('click', (e) => {
            if (e.target.id === 'addItemBtn') {
                this.showAddForm();
            }
            const decBtn = e.target.closest('.btn-dec');
            const incBtn = e.target.closest('.btn-inc');
            const delBtn = e.target.closest('.btn-delete-item');
            const submitAdd = e.target.closest('#submitAddItem');
            const cancelAdd = e.target.closest('#cancelAddItem');
            const searchResult = e.target.closest('.search-result-item');

            if (decBtn) this.adjustQuantity(decBtn.dataset.id, -1);
            if (incBtn) this.adjustQuantity(incBtn.dataset.id, 1);
            if (delBtn) this.deleteItem(delBtn.dataset.id);
            if (submitAdd) this.submitAddItem();
            if (cancelAdd) this.clearAddForm();
            if (searchResult) this.selectIngredient(searchResult.dataset);
        });

        document.addEventListener('submit', (e) => {
            if (e.target.id === 'addToPantryForm') {
                e.preventDefault();
                this.submitAddItem();
            }
        });
    }

    handleSearchInput(value) {
        clearTimeout(this.searchTimeout);
        const query = value.trim();
        if (query.length < 2) {
            this.searchResults = [];
            this.renderSearchResults();
            return;
        }
        this.searchTimeout = setTimeout(() => this.searchIngredients(query), 300);
    }

    async searchIngredients(query) {
        try {
            const result = await api.searchIngredients(query);
            this.searchResults = (result && result.items) || [];
            this.renderSearchResults();
        } catch (err) {
            router.setState({ error: 'Failed to search ingredients' });
        }
    }

    renderSearchResults() {
        const container = document.getElementById('pantryContent');
        if (!container) return;
        const existing = container.querySelector('.search-results');
        if (existing) existing.remove();

        if (this.searchResults.length === 0) return;

        const div = document.createElement('div');
        div.className = 'search-results';
        div.innerHTML = this.searchResults.map(ing => `
            <div class="search-result-item" data-fdc-id="${ing.fdcId}" data-name="${this.escapeHtml(ing.description)}">
                <span class="result-name">${this.escapeHtml(ing.description)}</span>
            </div>
        `).join('');
        container.prepend(div);
    }

    selectIngredient(data) {
        this.selectedIngredient = { fdcId: parseInt(data.fdcId), name: data.name };
        this.showAddForm();
    }

    showAddForm() {
        const container = document.getElementById('pantryContent');
        if (!container) return;

        const existing = container.querySelector('.add-pantry-form');
        if (existing) existing.remove();

        const form = document.createElement('div');
        form.className = 'add-pantry-form';
        form.innerHTML = `
            <form id="addToPantryForm">
                <h4>Add to Pantry</h4>
                <div class="form-group">
                    <label>Ingredient</label>
                    <input type="text" id="addIngredientName" value="${this.escapeHtml(this.selectedIngredient?.name || '')}" placeholder="Ingredient name" required>
                </div>
                <div class="form-row">
                    <div class="form-group">
                        <label for="addQuantity">Quantity</label>
                        <input type="number" id="addQuantity" step="0.01" min="0.01" required>
                    </div>
                    <div class="form-group">
                        <label for="addUnit">Unit</label>
                        <input type="text" id="addUnit" placeholder="g, ml, pc..." required>
                    </div>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary" id="submitAddItem">Add</button>
                    <button type="button" class="btn btn-secondary" id="cancelAddItem">Cancel</button>
                </div>
            </form>
        `;
        container.prepend(form);
    }

    clearAddForm() {
        this.selectedIngredient = null;
        const form = document.querySelector('.add-pantry-form');
        if (form) form.remove();
        const input = document.getElementById('searchPantry');
        if (input) input.value = '';
        this.searchResults = [];
        this.renderSearchResults();
    }

    async submitAddItem() {
        const quantity = parseFloat(document.getElementById('addQuantity')?.value);
        const unit = document.getElementById('addUnit')?.value?.trim();

        if (!quantity || !unit) {
            router.setState({ error: 'Please enter quantity and unit' });
            return;
        }

        if (!this.selectedIngredient) {
            router.setState({ error: 'Please select an ingredient from search results first' });
            return;
        }

        const item = {
            fdcId: this.selectedIngredient.fdcId,
            quantity,
            unit,
        };
        this.clearAddForm();
        try {
            const result = await api.addPantryItem(item);
            console.log('Pantry item added:', result);
        } catch (err) {
            console.error('Add pantry item responded with error (item may still be saved):', err.message, 'status:', err.status, 'data:', err.data);
        }
        await this.loadItems();
    }

    async adjustQuantity(id, delta) {
        try {
            await api.updatePantryItem(id, delta);
            await this.loadItems();
        } catch (err) {
            router.setState({ error: err.message });
        }
    }

    async deleteItem(id) {
        try {
            await api.deletePantryItem(id);
            await this.loadItems();
        } catch (err) {
            router.setState({ error: err.message });
        }
    }

    async loadItems() {
        try {
            const result = await api.getPantryItems();
            this.items = (result && result.items) || [];
            this.renderItems();
        } catch (err) {
            router.setState({ error: 'Failed to load pantry items' });
        }
    }

    renderItems() {
        const container = document.getElementById('pantryContent');
        if (!container) return;

        let listEl = container.querySelector('.pantry-list');
        if (!listEl) {
            listEl = document.createElement('div');
            listEl.className = 'pantry-list';
            container.appendChild(listEl);
        }

        if (this.items.length === 0) {
            listEl.innerHTML = '<p class="text-muted text-center">Your pantry is empty. Search and add ingredients above.</p>';
        } else {
            listEl.innerHTML = this.items.map(item => `
                <div class="pantry-item" data-id="${item.id}">
                    <div class="pantry-item-info">
                        <span class="pantry-item-name">${this.escapeHtml(item.food.description)}</span>
                        <span class="pantry-item-qty">${item.quantity} ${item.unit}</span>
                    </div>
                    <div class="pantry-item-actions">
                        <button class="btn btn-small btn-dec" data-id="${item.id}">−</button>
                        <button class="btn btn-small btn-inc" data-id="${item.id}">+</button>
                        <button class="btn btn-small btn-danger btn-delete-item" data-id="${item.id}">&times;</button>
                    </div>
                </div>
            `).join('');
        }
    }

    init() {
        this.clearAddForm();
        this.loadItems();
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

const pantryPageHandler = new PantryPageHandler();
