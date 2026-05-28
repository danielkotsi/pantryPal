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
            const response = await api.sendChatMessage(`Generate ${periodType} plan`, periodType);
            const loadingEl = document.getElementById(loadingId);
            if (loadingEl) loadingEl.remove();

            if (response.botMessage && response.botMessage.content) {
                this.appendMessage(container, 'bot', response.botMessage.content);
            }
            this.scrollToBottom();
        } catch (err) {
            const loadingEl = document.getElementById(loadingId);
            if (loadingEl) loadingEl.remove();
            router.setState({ error: err.message });
        }
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
