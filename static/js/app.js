// ============================================
// APP INITIALIZATION & STATE MANAGEMENT
// ============================================

class URLShortenerApp {
    constructor() {
        this.currentResult = null;
        this.history = [];
        this.currentChart = null;
        this.elements = this.cacheDOM();
        this.init();
    }

    cacheDOM() {
        return {
            form: document.getElementById('shorten-form'),
            originalUrl: document.getElementById('original-url'),
            customShort: document.getElementById('custom-short'),
            submitBtn: document.getElementById('submit-btn'),
            resultCard: document.getElementById('result-card'),
            shortUrlDisplay: document.getElementById('short-url-display'),
            loading: document.getElementById('loading'),
            toastContainer: document.getElementById('toast-container'),
            analyticsModal: document.getElementById('analytics-modal'),
            historyList: document.getElementById('history-list'),
            navLinks: document.querySelectorAll('.nav-link'),
            sections: document.querySelectorAll('.section')
        };
    }

    init() {
        this.setupEventListeners();
        this.loadHistory();
        this.animatePageLoad();
    }

    animatePageLoad() {
        document.body.style.opacity = '0';
        setTimeout(() => {
            document.body.style.transition = 'opacity 0.5s ease-in-out';
            document.body.style.opacity = '1';
        }, 100);
    }

    // ============================================
    // EVENT LISTENERS
    // ============================================

    setupEventListeners() {
        this.elements.form.addEventListener('submit', e => this.handleSubmit(e));

        this.elements.navLinks.forEach(link => {
            link.addEventListener('click', e => this.handleNavClick(e, link));
        });

        this.elements.analyticsModal.addEventListener('click', e => {
            if (e.target === this.elements.analyticsModal) this.closeAnalytics();
        });

        document.addEventListener('keydown', e => this.handleKeyboard(e));
        this.elements.originalUrl.focus();
    }

    handleNavClick(e, link) {
        e.preventDefault();
        const section = link.dataset.section;
        this.switchSection(section);

        this.elements.navLinks.forEach(l => l.classList.remove('active'));
        link.classList.add('active');
    }

    handleKeyboard(e) {
        if (e.key === 'Escape') this.closeAnalytics();
        if (e.ctrlKey && e.key === 'Enter') {
            this.elements.form.dispatchEvent(new Event('submit'));
        }
    }

    // ============================================
    // NAVIGATION & UI
    // ============================================

    switchSection(sectionName) {
        this.elements.sections.forEach(section => {
            section.classList.remove('active');
        });

        const targetSection = document.getElementById(`${sectionName}-section`);
        if (targetSection) {
            targetSection.classList.add('active');
            if (sectionName === 'history') this.loadHistory();
        }
    }

    setLoadingState(loading) {
        if (loading) {
            this.elements.submitBtn.classList.add('loading');
            this.elements.submitBtn.disabled = true;
            this.elements.submitBtn.innerHTML = '<i class="fas fa-spinner"></i><span>Создаем...</span>';
            this.elements.loading.classList.remove('hidden');
        } else {
            this.elements.submitBtn.classList.remove('loading');
            this.elements.submitBtn.disabled = false;
            this.elements.submitBtn.innerHTML = '<i class="fas fa-magic"></i><span>Сократить</span>';
            this.elements.loading.classList.add('hidden');
        }
    }

    // ============================================
    // URL SHORTENING
    // ============================================

    async handleSubmit(e) {
        e.preventDefault();

        const url = this.elements.originalUrl.value.trim();
        const customShort = this.elements.customShort.value.trim();

        if (!this.isValidUrl(url)) {
            this.showToast('Пожалуйста, введите корректный URL', 'error');
            return;
        }

        this.setLoadingState(true);

        try {
            const result = await this.shortenUrl(url, customShort);
            this.showResult(result);
            this.addToHistory(result);
            this.showToast('Ссылка успешно создана!', 'success');
        } catch (error) {
            this.showToast(error.message || 'Произошла ошибка при создании ссылки', 'error');
        } finally {
            this.setLoadingState(false);
        }
    }

    isValidUrl(string) {
        try {
            new URL(string);
            return true;
        } catch (_) {
            return false;
        }
    }

    async shortenUrl(url, customShort = '') {
        const requestBody = { url };
        if (customShort) requestBody.custom = customShort;

        const response = await fetch('/shorten', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(requestBody)
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || 'Ошибка сервера');
        }

        return {
            short: data.short,
            original: url,
            shortUrl: `${window.location.origin}/s/${data.short}`,
            expires: data.expires,
            createdAt: Date.now()
        };
    }

    showResult(result) {
        this.currentResult = result;
        this.elements.shortUrlDisplay.value = result.shortUrl;
        this.elements.resultCard.classList.remove('hidden');
        this.elements.resultCard.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }

    resetForm() {
        this.elements.form.reset();
        this.elements.resultCard.classList.add('hidden');
        this.elements.originalUrl.focus();
    }

    copyToClipboard() {
        this.elements.shortUrlDisplay.select();
        this.elements.shortUrlDisplay.setSelectionRange(0, 99999);

        try {
            document.execCommand('copy');
            this.showToast('Ссылка скопирована в буфер обмена!', 'success');
            this.animateCopyBtn();
        } catch (err) {
            this.showToast('Не удалось скопировать ссылку', 'error');
        }
    }

    animateCopyBtn() {
        const copyBtn = document.querySelector('.copy-btn');
        const originalIcon = copyBtn.innerHTML;
        copyBtn.innerHTML = '<i class="fas fa-check"></i>';
        copyBtn.style.background = 'var(--success-color)';

        setTimeout(() => {
            copyBtn.innerHTML = originalIcon;
            copyBtn.style.background = 'var(--primary-color)';
        }, 1500);
    }
    // ============================================
    // ANALYTICS
    // ============================================

    async viewAnalytics() {
        if (!this.currentResult) {
            this.showToast('Нет данных для аналитики', 'error');
            return;
        }

        try {
            this.setLoadingState(true);
            const analyticsData = await this.getAnalytics(this.currentResult.short);
            this.showAnalyticsModal(analyticsData);
        } catch (error) {
            const errorMsg = error.message.includes('not found')
                ? 'Ссылка не найдена или истекла'
                : 'Ошибка при загрузке аналитики';
            this.showToast(errorMsg, 'error');
        } finally {
            this.setLoadingState(false);
        }
    }

    async getAnalytics(shortCode) {
        const response = await fetch(`/analytics/${shortCode}`);

        if (response.status === 404) {
            throw new Error('URL not found');
        }

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Server error: ${response.status}`);
        }

        return await response.json();
    }

    async getDetailedAnalytics(shortCode) {
        try {
            const endDate = new Date();
            const startDate = new Date();
            startDate.setDate(startDate.getDate() - 30);

            const response = await fetch(
                `/analytics/${shortCode}/detailed?from=${startDate.toISOString().split('T')[0]}&to=${endDate.toISOString().split('T')[0]}`
            );

            return response.ok ? await response.json() : null;
        } catch (error) {
            console.warn('Could not load detailed analytics:', error);
            return null;
        }
    }

    async showAnalyticsModal(data) {
        const modal = this.elements.analyticsModal;
        document.getElementById('analytics-title').textContent = `Аналитика: ${data.short}`;
        document.getElementById('total-clicks').textContent = data.visit_count || 0;

        const detailedData = await this.getDetailedAnalytics(data.short);
        if (detailedData) {
            const today = new Date().toISOString().split('T')[0];
            document.getElementById('today-clicks').textContent = detailedData.daily_clicks?.[today] || 0;
            document.getElementById('mobile-percent').textContent =
                detailedData.mobile_percentage ? `${detailedData.mobile_percentage}%` : '0%';
        } else {
            document.getElementById('today-clicks').textContent = '0';
            document.getElementById('mobile-percent').textContent = '0%';
        }

        modal.classList.add('active');
        document.body.style.overflow = 'hidden';

        await this.generateClicksChart(data);
    }

    closeAnalytics() {
        const modal = this.elements.analyticsModal;
        modal.style.animation = 'fadeOut 0.3s ease-out';

        setTimeout(() => {
            modal.classList.remove('active');
            document.body.style.overflow = 'auto';
        }, 300);

        if (this.currentChart) {
            this.currentChart.destroy();
            this.currentChart = null;
        }
    }

    async generateClicksChart(data) {
        const canvas = document.getElementById('clicks-chart');
        const ctx = canvas.getContext('2d');

        if (this.currentChart) this.currentChart.destroy();

        const detailedData = await this.getDetailedAnalytics(data.short);
        const { labels, clickData } = this.prepareChartData(detailedData);

        this.currentChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels,
                datasets: [{
                    label: 'Переходы',
                    data: clickData,
                    borderColor: '#6366f1',
                    backgroundColor: 'rgba(99, 102, 241, 0.1)',
                    tension: 0.4,
                    fill: true,
                    pointBackgroundColor: '#6366f1',
                    pointBorderColor: '#ffffff',
                    pointBorderWidth: 2,
                    pointRadius: 6,
                    pointHoverRadius: 8
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { display: false } },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: { precision: 0 },
                        grid: { color: 'rgba(0, 0, 0, 0.05)' }
                    },
                    x: { grid: { display: false } }
                }
            }
        });

        await this.loadRecentClicks(data);
    }

    prepareChartData(detailedData) {
        const labels = [];
        const clickData = [];
        const today = new Date();

        for (let i = 6; i >= 0; i--) {
            const date = new Date(today);
            date.setDate(date.getDate() - i);
            const dateStr = date.toISOString().split('T')[0];

            labels.push(date.toLocaleDateString('ru-RU', {
                month: 'short',
                day: 'numeric'
            }));

            clickData.push(
                detailedData?.daily_clicks?.[dateStr] || 0
            );
        }

        return { labels, clickData };
    }

    async loadRecentClicks(data) {
        const container = document.getElementById('recent-clicks');

        try {
            const response = await fetch(`/analytics/${data.short}/recent-clicks`);
            if (response.ok) {
                const clicksData = await response.json();
                if (clicksData.clicks?.length) {
                    container.innerHTML = clicksData.clicks.map(click => `
                        <div class="click-item">
                            <div class="click-info">
                                <div class="time">${this.formatTimeAgo(new Date(click.occurred_at))}</div>
                                <div class="details">${click.ip || 'Unknown IP'} • ${click.referrer || 'Direct'}</div>
                            </div>
                            <div class="click-device">
                                <i class="fas ${this.getDeviceIcon(click.device || 'desktop')}"></i>
                                ${click.device || 'Desktop'}
                            </div>
                        </div>
                    `).join('');
                    return;
                }
            }
        } catch (error) {
            console.warn('Could not load recent clicks:', error);
        }

        container.innerHTML = data.visit_count === 0
            ? '<div class="no-data">Переходов пока нет</div>'
            : '<div class="no-data">Детальная информация о переходах недоступна</div>';
    }

    // ============================================
    // HISTORY MANAGEMENT
    // ============================================

    async loadHistory() {
        try {
            const savedHistory = localStorage.getItem('urlHistory');
            if (savedHistory) {
                this.history = JSON.parse(savedHistory);
            }

            await Promise.all(this.history.map(async (item) => {
                try {
                    const res = await fetch(`/analytics/${item.short}`);
                    if (res.ok) {
                        const data = await res.json();
                        item.visits = data.visit_count || 0;
                    }
                } catch (e) {
                    console.warn('Failed to load visit count for', item.short);
                }
            }));
        } catch (e) {
            console.warn('Could not load history:', e);
            this.history = [];
        }

        this.renderHistory();
    }

    addToHistory(result) {
        this.history.unshift({
            ...result,
            id: Date.now(),
            visits: 0
        });

        this.history = this.history.slice(0, 20);

        try {
            if (typeof (Storage) !== "undefined") {
                localStorage.setItem('urlHistory', JSON.stringify(this.history));
            }
        } catch (e) {
            console.warn('Could not save history:', e);
        }
    }

    renderHistory() {
        const container = this.elements.historyList;

        if (this.history.length === 0) {
            container.innerHTML = `
                <div class="no-history">
                    <i class="fas fa-history"></i>
                    <h3>История пуста</h3>
                    <p>Создайте первую короткую ссылку, чтобы она появилась здесь</p>
                </div>
            `;
            return;
        }

        container.innerHTML = this.history.map((item, index) => `
            <div class="history-item" style="animation-delay: ${index * 0.1}s">
                <div class="history-header">
                    <div class="history-info">
                        <h4>${this.truncateUrl(item.original, 60)}</h4>
                        <p>Создано ${new Date(item.createdAt).toLocaleDateString('ru-RU')}</p>
                    </div>
                </div>
                <div class="history-stats">
                    <div class="history-stat">
                        <div class="value">${item.visits}</div>
                        <div class="label">Переходов</div>
                    </div>
                    <div class="history-stat">
                        <div class="value">${item.short}</div>
                        <div class="label">Код</div>
                    </div>
                </div>
                <div class="history-actions">
                    <a href="${item.shortUrl}" target="_blank" class="history-btn">
                        <i class="fas fa-external-link-alt"></i>
                        Открыть
                    </a>
                    <button class="history-btn" onclick="app.copyHistoryUrl('${item.shortUrl}')">
                        <i class="fas fa-copy"></i>
                        Копировать
                    </button>
                    <button class="history-btn" onclick="app.viewHistoryAnalytics('${item.short}')">
                        <i class="fas fa-chart-bar"></i>
                        Аналитика
                    </button>
                    <button class="history-btn" onclick="app.deleteHistoryItem(${item.id})" style="color: var(--error-color);">
                        <i class="fas fa-trash"></i>
                        Удалить
                    </button>
                </div>
            </div>
        `).join('');
    }

    async viewHistoryAnalytics(shortCode) {
        try {
            this.setLoadingState(true);
            const analyticsData = await this.getAnalytics(shortCode);
            this.showAnalyticsModal(analyticsData);
        } catch (error) {
            if (error.message.includes('not found')) {
                this.showToast('Ссылка не найдена или истекла', 'warning');
                this.removeFromHistory(shortCode);
            } else {
                this.showToast('Ошибка при загрузке аналитики', 'error');
            }
        } finally {
            this.setLoadingState(false);
        }
    }

    removeFromHistory(shortCode) {
        this.history = this.history.filter(item => item.short !== shortCode);
        try {
            localStorage.setItem('urlHistory', JSON.stringify(this.history));
        } catch (e) {
            console.warn('Could not update history:', e);
        }
        this.renderHistory();
    }

    deleteHistoryItem(id) {
        if (!confirm('Вы уверены, что хотите удалить эту ссылку из истории?')) return;

        this.history = this.history.filter(item => item.id !== id);
        try {
            localStorage.setItem('urlHistory', JSON.stringify(this.history));
        } catch (e) {
            console.warn('Could not save history:', e);
        }

        this.renderHistory();
        this.showToast('Ссылка удалена из истории', 'success');
    }

    copyHistoryUrl(url) {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(url).then(() => {
                this.showToast('Ссылка скопирована!', 'success');
            }).catch(() => this.fallbackCopyTextToClipboard(url));
        } else {
            this.fallbackCopyTextToClipboard(url);
        }
    }

    fallbackCopyTextToClipboard(text) {
        const textArea = document.createElement("textarea");
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
            this.showToast('Ссылка скопирована!', 'success');
        } catch (err) {
            this.showToast('Не удалось скопировать ссылку', 'error');
        }

        document.body.removeChild(textArea);
    }

    // ============================================
    // UTILITIES
    // ============================================

    truncateUrl(url, length) {
        return url.length <= length ? url : url.substring(0, length) + '...';
    }

    formatTimeAgo(date) {
        const now = new Date();
        const diffInSeconds = Math.floor((now - date) / 1000);

        if (diffInSeconds < 60) return `${diffInSeconds} сек назад`;
        if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} мин назад`;
        if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} ч назад`;
        return `${Math.floor(diffInSeconds / 86400)} дн назад`;
    }

    getDeviceIcon(device) {
        const icons = {
            'mobile': 'fa-mobile-alt',
            'tablet': 'fa-tablet-alt',
            'desktop': 'fa-desktop'
        };
        return icons[device.toLowerCase()] || icons['desktop'];
    }

    showToast(message, type = 'success') {
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;

        const iconMap = {
            success: 'fa-check-circle',
            error: 'fa-exclamation-circle',
            warning: 'fa-exclamation-triangle',
            info: 'fa-info-circle'
        };

        toast.innerHTML = `
            <i class="fas ${iconMap[type] || iconMap.success}"></i>
            <span class="toast-message">${message}</span>
            <button class="toast-close" onclick="this.parentElement.remove()">
                <i class="fas fa-times"></i>
            </button>
        `;

        this.elements.toastContainer.appendChild(toast);

        setTimeout(() => {
            if (toast.parentElement) {
                toast.style.animation = 'slideOutRight 0.3s ease-out';
                setTimeout(() => toast.remove(), 300);
            }
        }, 5000);
    }
}

// ============================================
// INITIALIZATION
// ============================================

let app;
document.addEventListener('DOMContentLoaded', function () {
    app = new URLShortenerApp();
});

// ============================================
// GLOBAL FUNCTION WRAPPERS (FOR HTML onclick)
// ============================================

function copyToClipboard() { app.copyToClipboard(); }
function resetForm() { app.resetForm(); }
function viewAnalytics() { app.viewAnalytics(); }
function closeAnalytics() { app.closeAnalytics(); }