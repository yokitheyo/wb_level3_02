// Global variables
let currentChart = null;
let history = [];

// DOM elements
const elements = {
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

// Initialize app
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

function initializeApp() {
    setupEventListeners();
    loadHistory();

    // Add smooth animations
    document.body.style.opacity = '0';
    setTimeout(() => {
        document.body.style.transition = 'opacity 0.5s ease-in-out';
        document.body.style.opacity = '1';
    }, 100);
}

function setupEventListeners() {
    // Form submission
    elements.form.addEventListener('submit', handleSubmit);

    // Navigation
    elements.navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const section = link.dataset.section;
            switchSection(section);

            // Update active nav link
            elements.navLinks.forEach(l => l.classList.remove('active'));
            link.classList.add('active');
        });
    });

    // Close modal on outside click
    elements.analyticsModal.addEventListener('click', (e) => {
        if (e.target === elements.analyticsModal) {
            closeAnalytics();
        }
    });

    // Keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            closeAnalytics();
        }
        if (e.ctrlKey && e.key === 'Enter') {
            elements.form.dispatchEvent(new Event('submit'));
        }
    });

    // Auto-focus on URL input
    elements.originalUrl.focus();
}

function switchSection(sectionName) {
    elements.sections.forEach(section => {
        section.classList.remove('active');
    });

    const targetSection = document.getElementById(`${sectionName}-section`);
    if (targetSection) {
        targetSection.classList.add('active');

        if (sectionName === 'history') {
            loadHistory();
        }
    }
}

async function handleSubmit(e) {
    e.preventDefault();

    const url = elements.originalUrl.value.trim();
    const customShort = elements.customShort.value.trim();

    if (!isValidUrl(url)) {
        showToast('Пожалуйста, введите корректный URL', 'error');
        return;
    }

    setLoadingState(true);

    try {
        const result = await shortenUrl(url, customShort);
        showResult(result);
        addToHistory(result);
        showToast('Ссылка успешно создана!', 'success');
    } catch (error) {
        console.error('Error shortening URL:', error);
        showToast(error.message || 'Произошла ошибка при создании ссылки', 'error');
    } finally {
        setLoadingState(false);
    }
}

async function getDetailedAnalytics(shortCode) {
    try {
        // Получаем данные за последние 30 дней
        const endDate = new Date();
        const startDate = new Date();
        startDate.setDate(startDate.getDate() - 30);

        const response = await fetch(`/analytics/${shortCode}/detailed?from=${startDate.toISOString().split('T')[0]}&to=${endDate.toISOString().split('T')[0]}`);

        if (response.ok) {
            const data = await response.json();
            return data;
        }
    } catch (error) {
        console.warn('Could not load detailed analytics:', error);
    }

    return null;
}

async function shortenUrl(url, customShort = '') {
    const requestBody = {
        url: url
    };

    // If custom short is provided, add it to request
    if (customShort) {
        requestBody.custom = customShort;
    }

    const response = await fetch('/shorten', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
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

function showResult(result) {
    elements.shortUrlDisplay.value = result.shortUrl;
    elements.resultCard.classList.remove('hidden');
    elements.resultCard.scrollIntoView({ behavior: 'smooth', block: 'center' });

    // Store current result for analytics
    window.currentResult = result;

    // Add pulse animation to result card
    elements.resultCard.style.animation = 'none';
    setTimeout(() => {
        elements.resultCard.style.animation = 'slideInUp 0.5s ease-out';
    }, 10);
}

function setLoadingState(loading) {
    if (loading) {
        elements.submitBtn.classList.add('loading');
        elements.submitBtn.disabled = true;
        elements.submitBtn.innerHTML = '<i class="fas fa-spinner"></i><span>Создаем...</span>';
        elements.loading.classList.remove('hidden');
    } else {
        elements.submitBtn.classList.remove('loading');
        elements.submitBtn.disabled = false;
        elements.submitBtn.innerHTML = '<i class="fas fa-magic"></i><span>Сократить</span>';
        elements.loading.classList.add('hidden');
    }
}

function isValidUrl(string) {
    try {
        const url = new URL(string);
        return url.protocol === 'http:' || url.protocol === 'https:';
    } catch (_) {
        return false;
    }
}

function copyToClipboard() {
    elements.shortUrlDisplay.select();
    elements.shortUrlDisplay.setSelectionRange(0, 99999); // For mobile devices

    try {
        document.execCommand('copy');
        showToast('Ссылка скопирована в буфер обмена!', 'success');

        // Visual feedback
        const copyBtn = document.querySelector('.copy-btn');
        const originalIcon = copyBtn.innerHTML;
        copyBtn.innerHTML = '<i class="fas fa-check"></i>';
        copyBtn.style.background = 'var(--success-color)';

        setTimeout(() => {
            copyBtn.innerHTML = originalIcon;
            copyBtn.style.background = 'var(--primary-color)';
        }, 1500);

    } catch (err) {
        showToast('Не удалось скопировать ссылку', 'error');
    }
}

function resetForm() {
    elements.form.reset();
    elements.resultCard.classList.add('hidden');
    elements.originalUrl.focus();

    // Add fade out animation
    elements.resultCard.style.animation = 'fadeOut 0.3s ease-out';
}

async function viewAnalytics() {
    if (!window.currentResult) {
        showToast('Нет данных для аналитики', 'error');
        return;
    }

    try {
        setLoadingState(true);
        const analyticsData = await getAnalytics(window.currentResult.short);
        showAnalyticsModal(analyticsData);
    } catch (error) {
        console.error('Error loading analytics:', error);
        if (error.message.includes('not found')) {
            showToast('Ссылка не найдена или истекла', 'warning');
        } else {
            showToast('Ошибка при загрузке аналитики', 'error');
        }
    } finally {
        setLoadingState(false);
    }
}

async function getAnalytics(shortCode) {
    try {
        const response = await fetch(`/analytics/${shortCode}`);

        if (response.status === 404) {
            throw new Error('URL not found - it may have expired or been deleted');
        }

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || `Server error: ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Analytics fetch error:', error);
        throw error;
    }
}

async function showAnalyticsModal(data) {
    const modal = elements.analyticsModal;
    const title = document.getElementById('analytics-title');
    const totalClicks = document.getElementById('total-clicks');
    const todayClicks = document.getElementById('today-clicks');
    const mobilePercent = document.getElementById('mobile-percent');

    title.textContent = `Аналитика: ${data.short}`;
    totalClicks.textContent = data.visit_count || 0;

    // Загружаем реальные данные для сегодняшних переходов и мобильных устройств
    const detailedData = await getDetailedAnalytics(data.short);

    if (detailedData) {
        const today = new Date().toISOString().split('T')[0];
        todayClicks.textContent = detailedData.daily_clicks?.[today] || 0;
        mobilePercent.textContent = detailedData.mobile_percentage ? `${detailedData.mobile_percentage}%` : '0%';
    } else {
        todayClicks.textContent = '0';
        mobilePercent.textContent = '0%';
    }

    // Show modal with animation
    modal.classList.add('active');
    modal.style.animation = 'fadeIn 0.3s ease-out';

    // Генерируем график с реальными данными
    await generateClicksChart(data);

    // Disable body scroll
    document.body.style.overflow = 'hidden';
}

function closeAnalytics() {
    const modal = elements.analyticsModal;
    modal.style.animation = 'fadeOut 0.3s ease-out';

    setTimeout(() => {
        modal.classList.remove('active');
        modal.style.animation = '';
        document.body.style.overflow = 'auto';
    }, 300);

    // Destroy existing chart
    if (currentChart) {
        currentChart.destroy();
        currentChart = null;
    }
}

async function generateClicksChart(data) {
    const canvas = document.getElementById('clicks-chart');
    const ctx = canvas.getContext('2d');

    // Destroy existing chart
    if (currentChart) {
        currentChart.destroy();
    }

    // Получаем детальную аналитику
    const detailedData = await getDetailedAnalytics(data.short);

    const labels = [];
    const clickData = [];
    const today = new Date();

    // Создаем данные за последние 7 дней
    for (let i = 6; i >= 0; i--) {
        const date = new Date(today);
        date.setDate(date.getDate() - i);
        const dateStr = date.toISOString().split('T')[0];

        labels.push(date.toLocaleDateString('ru-RU', {
            month: 'short',
            day: 'numeric'
        }));

        // Используем только реальные данные или 0
        if (detailedData && detailedData.daily_clicks && detailedData.daily_clicks[dateStr]) {
            clickData.push(detailedData.daily_clicks[dateStr]);
        } else {
            clickData.push(0);
        }
    }

    currentChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
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
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        precision: 0
                    },
                    grid: {
                        color: 'rgba(0, 0, 0, 0.05)'
                    }
                },
                x: {
                    grid: {
                        display: false
                    }
                }
            },
            elements: {
                point: {
                    hoverBackgroundColor: '#4f46e5'
                }
            }
        }
    });

    // Загружаем список последних кликов
    await loadRecentClicks(data);
}

async function loadRecentClicks(data) {
    const container = document.getElementById('recent-clicks');

    try {
        const response = await fetch(`/analytics/${data.short}/recent-clicks`);

        if (response.ok) {
            const clicksData = await response.json();

            if (clicksData.clicks && clicksData.clicks.length > 0) {
                container.innerHTML = clicksData.clicks.map(click => `
                    <div class="click-item">
                        <div class="click-info">
                            <div class="time">${formatTimeAgo(new Date(click.occurred_at))}</div>
                            <div class="details">${click.ip || 'Unknown IP'} • ${click.referrer || 'Direct'}</div>
                        </div>
                        <div class="click-device">
                            <i class="fas ${getDeviceIcon(click.device || 'desktop')}"></i>
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

    // Если данных нет
    if (data.visit_count === 0) {
        container.innerHTML = '<div class="no-data">Переходов пока нет</div>';
    } else {
        container.innerHTML = '<div class="no-data">Детальная информация о переходах недоступна</div>';
    }
}

function getDeviceIcon(device) {
    switch (device.toLowerCase()) {
        case 'mobile': return 'fa-mobile-alt';
        case 'tablet': return 'fa-tablet-alt';
        default: return 'fa-desktop';
    }
}

function formatTimeAgo(date) {
    const now = new Date();
    const diffInSeconds = Math.floor((now - date) / 1000);

    if (diffInSeconds < 60) {
        return `${diffInSeconds} сек назад`;
    } else if (diffInSeconds < 3600) {
        const minutes = Math.floor(diffInSeconds / 60);
        return `${minutes} мин назад`;
    } else if (diffInSeconds < 86400) {
        const hours = Math.floor(diffInSeconds / 3600);
        return `${hours} ч назад`;
    } else {
        const days = Math.floor(diffInSeconds / 86400);
        return `${days} дн назад`;
    }
}

function addToHistory(result) {
    history.unshift({
        ...result,
        id: Date.now(),
        visits: 0
    });

    // Keep only last 20 items
    history = history.slice(0, 20);

    // Save to localStorage (if available)
    try {
        if (typeof(Storage) !== "undefined") {
            localStorage.setItem('urlHistory', JSON.stringify(history));
        }
    } catch (e) {
        console.warn('Could not save history to localStorage:', e);
    }
}

async function loadHistory() {
    try {
        if (typeof(Storage) !== "undefined") {
            const savedHistory = localStorage.getItem('urlHistory');
            if (savedHistory) {
                history = JSON.parse(savedHistory);
            }
        }

        await Promise.all(history.map(async (item) => {
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
        console.warn('Could not load history from localStorage:', e);
        history = [];
    }

    renderHistory();
}

function renderHistory() {
    const container = elements.historyList;

    if (history.length === 0) {
        container.innerHTML = `
            <div class="no-history">
                <i class="fas fa-history"></i>
                <h3>История пуста</h3>
                <p>Создайте первую короткую ссылку, чтобы она появилась здесь</p>
            </div>
        `;
        return;
    }

    container.innerHTML = history.map((item, index) => `
        <div class="history-item" style="animation-delay: ${index * 0.1}s">
            <div class="history-header">
                <div class="history-info">
                    <h4>${truncateUrl(item.original, 60)}</h4>
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
                <button class="history-btn" onclick="copyHistoryUrl('${item.shortUrl}')">
                    <i class="fas fa-copy"></i>
                    Копировать
                </button>
                <button class="history-btn" onclick="viewHistoryAnalytics('${item.short}')">
                    <i class="fas fa-chart-bar"></i>
                    Аналитика
                </button>
                <button class="history-btn" onclick="deleteHistoryItem(${item.id})" style="color: var(--error-color);">
                    <i class="fas fa-trash"></i>
                    Удалить
                </button>
            </div>
        </div>
    `).join('');
}

function truncateUrl(url, length) {
    if (url.length <= length) return url;
    return url.substring(0, length) + '...';
}

function copyHistoryUrl(url) {
    if (navigator.clipboard) {
        navigator.clipboard.writeText(url).then(() => {
            showToast('Ссылка скопирована!', 'success');
        }).catch(() => {
            fallbackCopyTextToClipboard(url);
        });
    } else {
        fallbackCopyTextToClipboard(url);
    }
}

function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        document.execCommand('copy');
        showToast('Ссылка скопирована!', 'success');
    } catch (err) {
        showToast('Не удалось скопировать ссылку', 'error');
    }

    document.body.removeChild(textArea);
}

async function viewHistoryAnalytics(shortCode) {
    try {
        setLoadingState(true);
        const analyticsData = await getAnalytics(shortCode);
        showAnalyticsModal(analyticsData);
    } catch (error) {
        console.error('Error loading analytics:', error);

        // Improved error handling
        if (error.message.includes('not found') || error.message.includes('URL not found')) {
            showToast('Ссылка не найдена или истекла. Возможно, она была удалена.', 'warning');
            // Remove this item from history since it doesn't exist
            removeFromHistory(shortCode);
        } else {
            showToast('Ошибка при загрузке аналитики', 'error');
        }
    } finally {
        setLoadingState(false);
    }
}

function removeFromHistory(shortCode) {
    history = history.filter(item => item.short !== shortCode);
    try {
        if (typeof(Storage) !== "undefined") {
            localStorage.setItem('urlHistory', JSON.stringify(history));
        }
    } catch (e) {
        console.warn('Could not update history in localStorage:', e);
    }
    renderHistory();
}

function deleteHistoryItem(id) {
    if (confirm('Вы уверены, что хотите удалить эту ссылку из истории?')) {
        history = history.filter(item => item.id !== id);

        try {
            if (typeof(Storage) !== "undefined") {
                localStorage.setItem('urlHistory', JSON.stringify(history));
            }
        } catch (e) {
            console.warn('Could not save history to localStorage:', e);
        }

        renderHistory();
        showToast('Ссылка удалена из истории', 'success');
    }
}

function showToast(message, type = 'success') {
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

    elements.toastContainer.appendChild(toast);

    // Auto remove after 5 seconds
    setTimeout(() => {
        if (toast.parentElement) {
            toast.style.animation = 'slideOutRight 0.3s ease-out';
            setTimeout(() => toast.remove(), 300);
        }
    }, 5000);
}

// Add some utility animations
const style = document.createElement('style');
style.textContent = `
    @keyframes fadeOut {
        from { opacity: 1; }
        to { opacity: 0; }
    }
    
    @keyframes slideOutRight {
        from {
            opacity: 1;
            transform: translateX(0);
        }
        to {
            opacity: 0;
            transform: translateX(100%);
        }
    }
    
    .no-history {
        text-align: center;
        padding: 60px 20px;
        color: white;
    }
    
    .no-history i {
        font-size: 48px;
        opacity: 0.5;
        margin-bottom: 20px;
    }
    
    .no-history h3 {
        font-size: 24px;
        margin-bottom: 12px;
        opacity: 0.9;
    }
    
    .no-history p {
        opacity: 0.7;
        max-width: 400px;
        margin: 0 auto;
    }
    
    .no-data {
        text-align: center;
        padding: 40px;
        color: var(--text-secondary);
    }
    
    /* Loading state improvements */
    .submit-btn.loading {
        background: var(--text-secondary) !important;
        cursor: not-allowed;
    }
    
    /* Chart container sizing */
    .chart-container {
        height: 300px;
        position: relative;
    }
    
    /* Modal improvements */
    .modal-content {
        max-height: calc(100vh - 40px);
        overflow-y: auto;
    }
    
    /* Analytics stats hover effects */
    .stat-card:hover i {
        transform: scale(1.1);
        transition: transform 0.2s ease;
    }
    
    /* History item improvements */
    .history-item:hover {
        transform: translateY(-2px);
        box-shadow: var(--shadow-lg);
    }
    
    .history-btn:hover {
        background: var(--bg-secondary);
        transform: translateY(-1px);
    }
`;

document.head.appendChild(style);