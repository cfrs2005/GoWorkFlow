// Main Application Logic

function appData() {
    return {
        currentPage: 'dashboard',
        loading: false,
        error: null,

        // Page titles
        get pageTitle() {
            const titles = {
                dashboard: '仪表盘',
                flows: '流程管理',
                tasks: '任务库',
                jobs: '作业监控',
            };
            return titles[this.currentPage] || '工作流管理系统';
        },

        // Initialize
        async init() {
            // Set initial page from hash or default to dashboard
            const hash = window.location.hash.slice(1);
            if (hash && ['dashboard', 'flows', 'tasks', 'jobs'].includes(hash)) {
                this.currentPage = hash;
            }

            // Watch for page changes
            this.$watch('currentPage', (value) => {
                window.location.hash = value;
                this.loadPageContent(value);
            });

            // Load initial content
            await this.loadPageContent(this.currentPage);

            // Auto-refresh every 30 seconds for dashboard and jobs
            setInterval(() => {
                if (this.currentPage === 'dashboard' || this.currentPage === 'jobs') {
                    this.refreshData();
                }
            }, 30000);
        },

        // Load page content
        async loadPageContent(page) {
            this.loading = true;
            this.error = null;

            try {
                switch (page) {
                    case 'dashboard':
                        await loadDashboard();
                        break;
                    case 'flows':
                        await loadFlows();
                        break;
                    case 'tasks':
                        await loadTasks();
                        break;
                    case 'jobs':
                        await loadJobs();
                        break;
                }
            } catch (error) {
                console.error('Error loading page:', error);
                this.error = error.message;
                this.showError(error.message);
            } finally {
                this.loading = false;
            }
        },

        // Refresh current page data
        async refreshData() {
            await this.loadPageContent(this.currentPage);
        },

        // Show error message
        showError(message) {
            alert(`错误: ${message}`);
        },

        // Show success message
        showSuccess(message) {
            alert(`成功: ${message}`);
        },

        // Confirm dialog
        confirm(message) {
            return window.confirm(message);
        },
    };
}

// Utility functions

function formatDate(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
    });
}

function formatStatus(status) {
    const statusMap = {
        pending: '等待中',
        running: '运行中',
        completed: '已完成',
        failed: '失败',
        skipped: '已跳过',
        rolled_back: '已回滚',
    };
    return statusMap[status] || status;
}

function getStatusClass(status) {
    const classMap = {
        pending: 'status-pending',
        running: 'status-running',
        completed: 'status-completed',
        failed: 'status-failed',
        skipped: 'status-skipped',
        rolled_back: 'status-skipped',
    };
    return `status-badge ${classMap[status] || ''}`;
}

function formatTaskType(type) {
    const typeMap = {
        manual: '手动任务',
        automated: '自动化任务',
        approval: '审批任务',
    };
    return typeMap[type] || type;
}

function truncate(str, length = 50) {
    if (!str) return '-';
    return str.length > length ? str.substring(0, length) + '...' : str;
}
