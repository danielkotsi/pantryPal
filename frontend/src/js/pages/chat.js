class ChatPageHandler {
    constructor() {
        this.currentProposal = null;
        this.setupDelegatedListeners();
    }

    setupDelegatedListeners() {
        document.addEventListener('click', (e) => {
            const generateBtn = e.target.closest('[data-generate]');
            if (generateBtn) {
                e.preventDefault();
                this.handleGenerate(generateBtn.dataset.generate);
                return;
            }

            const acceptBtn = e.target.closest('.proposal-accept');
            if (acceptBtn) {
                e.preventDefault();
                this.handleAccept(acceptBtn.dataset.proposalId);
                return;
            }

            const declineBtn = e.target.closest('.proposal-decline');
            if (declineBtn) {
                e.preventDefault();
                this.handleDecline(declineBtn.dataset.proposalId);
                return;
            }
        });
    }

    init() {
        this.currentProposal = null;
        this.historyVersion = 0;
        const container = document.getElementById('chatMessages');
        if (container) container.innerHTML = '';
        this.loadChatHistory();

        const form = document.getElementById('chatForm');
        if (form) {
            form.onsubmit = (e) => {
                e.preventDefault();
                const input = document.getElementById('chatInput');
                if (input) this.handleChatSubmit(input);
            };
        }
    }

    async loadChatHistory() {
        const container = document.getElementById('chatMessages');
        if (!container) return;
        const version = ++this.historyVersion;
        try {
            const data = await api.getChatHistory();
            if (version !== this.historyVersion) return;
            const messages = data.messages || [];
            container.innerHTML = messages.map(msg => this.renderMessage(msg)).join('');
            this.scrollToBottom();
        } catch (err) {
            if (version === this.historyVersion) {
                container.innerHTML = '';
            }
        }
    }

    renderMessage(msg) {
        const isUser = msg.role === 'user';
        const time = msg.createdAt
            ? new Date(msg.createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
            : '';
        const actionBadge = msg.action ? `<span class="action-badge">${msg.action}</span>` : '';
        return `
            <div class="chat-message ${isUser ? 'user-message' : 'bot-message'}">
                <p>${this.escapeHtml(msg.content)}</p>
                ${actionBadge}
                <span class="chat-timestamp">${time}</span>
            </div>
        `;
    }

    async handleChatSubmit(input) {
        const message = input.value.trim();
        if (!message) return;

        const container = document.getElementById('chatMessages');
        this.appendMessage(container, 'user', message);
        input.value = '';
        this.scrollToBottom();

        try {
            const response = await api.sendChatMessage(message);
            if (response.botMessage && response.botMessage.content) {
                this.appendMessage(container, 'bot', response.botMessage.content);
            } else {
                this.appendMessage(container, 'bot', 'AI assistant is not available. Use the action buttons above to generate meal plans.');
            }
            this.scrollToBottom();
        } catch (err) {
            this.appendMessage(container, 'bot', 'Sorry, I encountered an error. Please try again.');
            this.scrollToBottom();
            router.setState({ error: err.message });
        }
    }

    appendMessage(container, role, text) {
        const cls = role === 'user' ? 'user-message' : 'bot-message';
        const time = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
        container.innerHTML += `
            <div class="chat-message ${cls}">
                <p>${this.escapeHtml(text)}</p>
                <span class="chat-timestamp">${time}</span>
            </div>
        `;
    }

    async handleGenerate(periodType) {
        const container = document.getElementById('chatMessages');
        this.appendMessage(container, 'user', `Generate ${periodType} plan`);
        this.scrollToBottom();

        const loadingId = 'generating-msg-' + Date.now();
        container.innerHTML += `
            <div class="chat-message bot-message" id="${loadingId}">
                <p>Generating ${periodType} plan...</p>
                <div class="spinner-small"></div>
            </div>
        `;
        this.scrollToBottom();

        try {
            const response = await api.generatePlan(periodType);
            const loadingEl = document.getElementById(loadingId);
            if (loadingEl) loadingEl.remove();

            this.currentProposal = response.proposal;
            const isFallback = response.fallbackActive;

            if (isFallback) {
                this.appendMessage(container, 'bot', 'Fallback mode active — using seeded data instead of AI.');
            }

            this.appendMessage(container, 'bot', `Here's your ${periodType} plan:`);
            container.innerHTML += this.renderProposalPreview(response.proposal, isFallback);
            this.scrollToBottom();
        } catch (err) {
            const loadingEl = document.getElementById(loadingId);
            if (loadingEl) loadingEl.remove();
            router.setState({ error: err.message });
        }
    }

    renderProposalPreview(proposal, isFallback) {
        const plan = proposal.plan || {};
        const dayCount = proposal.days ? proposal.days.length : 0;
        const mealCount = proposal.days
            ? proposal.days.reduce((sum, d) => sum + Object.keys(d.sections || {}).length, 0)
            : 0;

        const weekTotal = proposal.weekTotals || {};

        let html = '<div class="proposal-preview">';
        html += `
            <div class="proposal-header">
                <span class="proposal-type">${plan.periodType} plan</span>
                ${isFallback ? '<span class="fallback-badge">Fallback</span>' : '<span class="ai-badge">AI</span>'}
                <span class="proposal-status">${plan.status || 'pending'}</span>
            </div>
            <div class="proposal-meta">
                <span>${dayCount} days \u00B7 ${mealCount} meals</span>
                ${plan.aiCostCentsTotal != null ? `<span>$${(plan.aiCostCentsTotal / 100).toFixed(2)}</span>` : ''}
                ${plan.source ? `<span>Source: ${plan.source}</span>` : ''}
                ${plan.proposalVersion ? `<span>v${plan.proposalVersion}</span>` : ''}
                ${plan.startDate ? `<span>${plan.startDate}${plan.endDate ? ' \u2192 ' + plan.endDate : ''}</span>` : ''}
            </div>
        `;

        if (proposal.days && proposal.days.length > 0) {
            html += '<div class="proposal-days">';
            proposal.days.forEach(d => {
                const t = d.totals || {};
                html += `<div class="proposal-day">
                    <div class="day-header">
                        <span class="day-date">${d.date || ''}</span>
                        <span class="day-total">${t.calories || 0} kcal  P:${t.proteinG || 0}g  C:${t.carbsG || 0}g  F:${t.fatG || 0}g</span>
                    </div>`;
                const sections = d.sections || {};
                ['breakfast', 'lunch', 'dinner', 'snacks'].forEach(key => {
                    const meal = sections[key];
                    if (!meal) return;
                    const mealCost = meal.estimatedCostCents != null ? `$${(meal.estimatedCostCents / 100).toFixed(2)}` : '';
                    html += `
                        <div class="meal-row">
                            <span class="meal-section-badge">${key}</span>
                            <span class="meal-name">${this.escapeHtml(meal.recipeName)}</span>
                            <span class="meal-macros">${meal.macros.calories || 0} kcal  P:${meal.macros.proteinG || 0}g  C:${meal.macros.carbsG || 0}g  F:${meal.macros.fatG || 0}g</span>
                            <span class="meal-cost">${mealCost}</span>
                        </div>`;
                });
                html += '</div>';
            });
            html += '</div>';
        }

        if (weekTotal.calories || weekTotal.proteinG || weekTotal.carbsG || weekTotal.fatG) {
            html += `
                <div class="proposal-week-total">
                    <h4>Week Total</h4>
                    <div class="macro-row">
                        <span class="macro-val">${weekTotal.calories || 0} kcal</span>
                        <span class="macro p">P: ${weekTotal.proteinG || 0}g</span>
                        <span class="macro c">C: ${weekTotal.carbsG || 0}g</span>
                        <span class="macro f">F: ${weekTotal.fatG || 0}g</span>
                    </div>
                </div>
            `;
        }

        html += `
            <div class="proposal-actions">
                <button class="btn btn-primary proposal-accept" data-proposal-id="${plan.id}">Accept</button>
                <button class="btn btn-secondary proposal-decline" data-proposal-id="${plan.id}">Decline</button>
            </div>
        `;
        html += '</div>';
        return html;
    }

    async handleAccept(proposalId) {
        const container = document.getElementById('chatMessages');
        try {
            await api.acceptProposal(proposalId);
            this.appendMessage(container, 'bot', '✅ Plan accepted! It has been saved to your planner.');
            this.currentProposal = null;
            const preview = container.querySelector('.proposal-preview');
            if (preview) preview.remove();
            this.scrollToBottom();
        } catch (err) {
            router.setState({ error: err.message });
        }
    }

    async handleDecline(proposalId) {
        const container = document.getElementById('chatMessages');
        const reason = prompt('Optional reason for declining:');
        try {
            await api.declineProposal(proposalId, reason || undefined);
            this.appendMessage(container, 'bot', 'Plan declined. You can generate a new one.');
            this.currentProposal = null;
            const preview = container.querySelector('.proposal-preview');
            if (preview) preview.remove();
            this.scrollToBottom();
        } catch (err) {
            router.setState({ error: err.message });
        }
    }

    scrollToBottom() {
        const container = document.getElementById('chatMessages');
        if (container) {
            container.scrollTop = container.scrollHeight;
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

const chatPageHandler = new ChatPageHandler();
