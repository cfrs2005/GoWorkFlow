// Dashboard Page Logic

let dashboardData = {
    stats: {
        totalFlows: 0,
        totalJobs: 0,
        runningJobs: 0,
        completedToday: 0,
    },
    recentJobs: [],
    chart: null,
};

async function loadDashboard() {
    try {
        // Fetch data
        const [flows, jobs] = await Promise.all([
            api.getFlows(),
            api.getJobs(),
        ]);

        // Calculate stats
        dashboardData.stats.totalFlows = flows.data?.length || 0;
        dashboardData.stats.totalJobs = jobs.data?.length || 0;
        dashboardData.stats.runningJobs = jobs.data?.filter(j => j.status === 'running').length || 0;

        // Get today's completed jobs
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        dashboardData.stats.completedToday = jobs.data?.filter(j => {
            if (j.status !== 'completed') return false;
            const updatedAt = new Date(j.updated_at);
            return updatedAt >= today;
        }).length || 0;

        // Get recent jobs (last 10)
        dashboardData.recentJobs = (jobs.data || [])
            .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
            .slice(0, 10);

        // Render dashboard
        renderDashboard();
    } catch (error) {
        console.error('Failed to load dashboard:', error);
        throw error;
    }
}

function renderDashboard() {
    const container = document.getElementById('dashboard-content');

    container.innerHTML = `
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div class="stat-card">
                <div class="stat-value">${dashboardData.stats.totalFlows}</div>
                <div class="stat-label">总流程数</div>
            </div>
            <div class="stat-card" style="background: linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%);">
                <div class="stat-value">${dashboardData.stats.totalJobs}</div>
                <div class="stat-label">总作业数</div>
            </div>
            <div class="stat-card" style="background: linear-gradient(135deg, #10B981 0%, #34D399 100%);">
                <div class="stat-value">${dashboardData.stats.runningJobs}</div>
                <div class="stat-label">运行中作业</div>
            </div>
            <div class="stat-card" style="background: linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%);">
                <div class="stat-value">${dashboardData.stats.completedToday}</div>
                <div class="stat-label">今日完成</div>
            </div>
        </div>

        <!-- Charts Row -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
            <!-- Job Status Chart -->
            <div class="card">
                <div class="card-header">作业状态分布</div>
                <canvas id="jobStatusChart" height="200"></canvas>
            </div>

            <!-- Daily Trend Chart -->
            <div class="card">
                <div class="card-header">7日趋势</div>
                <canvas id="dailyTrendChart" height="200"></canvas>
            </div>
        </div>

        <!-- Recent Jobs Table -->
        <div class="card">
            <div class="card-header">最近作业</div>
            <div class="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>作业 ID</th>
                            <th>流程 ID</th>
                            <th>状态</th>
                            <th>创建时间</th>
                            <th>更新时间</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${dashboardData.recentJobs.length > 0 ? dashboardData.recentJobs.map(job => `
                            <tr>
                                <td><strong>#${job.id}</strong></td>
                                <td>${job.flow_id}</td>
                                <td><span class="${getStatusClass(job.status)}">${formatStatus(job.status)}</span></td>
                                <td>${formatDate(job.created_at)}</td>
                                <td>${formatDate(job.updated_at)}</td>
                                <td>
                                    <button onclick="viewJobDetails(${job.id})" class="btn btn-sm btn-primary">
                                        查看详情
                                    </button>
                                </td>
                            </tr>
                        `).join('') : `
                            <tr>
                                <td colspan="6" class="text-center text-gray-500 py-8">暂无数据</td>
                            </tr>
                        `}
                    </tbody>
                </table>
            </div>
        </div>
    `;

    // Render charts
    setTimeout(() => {
        renderJobStatusChart();
        renderDailyTrendChart();
    }, 100);
}

function renderJobStatusChart() {
    const ctx = document.getElementById('jobStatusChart');
    if (!ctx) return;

    // Calculate status counts
    const statusCounts = dashboardData.recentJobs.reduce((acc, job) => {
        acc[job.status] = (acc[job.status] || 0) + 1;
        return acc;
    }, {});

    if (dashboardData.chart) {
        dashboardData.chart.destroy();
    }

    dashboardData.chart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['等待中', '运行中', '已完成', '失败'],
            datasets: [{
                data: [
                    statusCounts.pending || 0,
                    statusCounts.running || 0,
                    statusCounts.completed || 0,
                    statusCounts.failed || 0,
                ],
                backgroundColor: [
                    '#FCD34D',
                    '#60A5FA',
                    '#34D399',
                    '#F87171',
                ],
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                },
            },
        },
    });
}

function renderDailyTrendChart() {
    const ctx = document.getElementById('dailyTrendChart');
    if (!ctx) return;

    // Generate last 7 days data
    const days = [];
    const data = [];
    for (let i = 6; i >= 0; i--) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        days.push(date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' }));

        const dayStart = new Date(date);
        dayStart.setHours(0, 0, 0, 0);
        const dayEnd = new Date(date);
        dayEnd.setHours(23, 59, 59, 999);

        const count = dashboardData.recentJobs.filter(j => {
            const created = new Date(j.created_at);
            return created >= dayStart && created <= dayEnd;
        }).length;

        data.push(count);
    }

    new Chart(ctx, {
        type: 'line',
        data: {
            labels: days,
            datasets: [{
                label: '作业数量',
                data: data,
                borderColor: '#7C3AED',
                backgroundColor: 'rgba(124, 58, 237, 0.1)',
                tension: 0.4,
                fill: true,
            }],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false,
                },
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1,
                    },
                },
            },
        },
    });
}

function viewJobDetails(jobId) {
    // Navigate to jobs page with this job
    window.appData = Alpine.store('app') || window.appData;
    if (window.appData && window.appData.currentPage) {
        window.appData.currentPage = 'jobs';
        setTimeout(() => {
            // Highlight the job or open details
            console.log('View job:', jobId);
        }, 100);
    }
}
