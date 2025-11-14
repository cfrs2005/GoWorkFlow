// Jobs Monitor Page Logic

let jobsData = {
    jobs: [],
    selectedJob: null,
    jobTasks: [],
    autoRefresh: true,
};

async function loadJobs() {
    try {
        const response = await api.getJobs();
        jobsData.jobs = response.data || [];
        jobsData.jobs.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
        renderJobs();
    } catch (error) {
        console.error('Failed to load jobs:', error);
        throw error;
    }
}

function renderJobs() {
    const container = document.getElementById('jobs-content');

    container.innerHTML = `
        <!-- Header Actions -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <p class="text-gray-600">实时监控作业执行状态</p>
            </div>
            <div class="flex space-x-2">
                <label class="flex items-center">
                    <input type="checkbox" ${jobsData.autoRefresh ? 'checked' : ''}
                           onchange="toggleAutoRefresh(this.checked)" class="mr-2">
                    <span class="text-sm text-gray-700">自动刷新</span>
                </label>
            </div>
        </div>

        <!-- Status Filter -->
        <div class="flex space-x-2 mb-6">
            <button onclick="filterJobs('all')" class="btn btn-sm btn-primary">全部</button>
            <button onclick="filterJobs('pending')" class="btn btn-sm btn-secondary">等待中</button>
            <button onclick="filterJobs('running')" class="btn btn-sm btn-secondary">运行中</button>
            <button onclick="filterJobs('completed')" class="btn btn-sm btn-secondary">已完成</button>
            <button onclick="filterJobs('failed')" class="btn btn-sm btn-secondary">失败</button>
        </div>

        <!-- Jobs List -->
        <div class="card">
            <div class="card-header">作业列表</div>
            <div class="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>作业 ID</th>
                            <th>流程 ID</th>
                            <th>状态</th>
                            <th>进度</th>
                            <th>创建时间</th>
                            <th>更新时间</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${jobsData.jobs.length > 0 ? jobsData.jobs.map(job => renderJobRow(job)).join('') : `
                            <tr>
                                <td colspan="7" class="text-center text-gray-500 py-8">暂无作业</td>
                            </tr>
                        `}
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Job Details Modal -->
        <div id="jobDetailsModal" class="modal-overlay" style="display: none;">
            <div class="modal-content" style="max-width: 900px;">
                <div class="modal-header">
                    <h3 class="modal-title">作业详情</h3>
                </div>
                <div class="modal-body">
                    <div id="jobDetailsContent"></div>
                </div>
                <div class="modal-footer">
                    <button onclick="closeJobDetailsModal()" class="btn btn-secondary">关闭</button>
                </div>
            </div>
        </div>
    `;
}

function renderJobRow(job) {
    // Calculate progress if possible
    let progressHTML = '<span class="text-gray-400">-</span>';
    if (job.status === 'running') {
        progressHTML = `
            <div class="w-full bg-gray-200 rounded-full h-2">
                <div class="bg-primary h-2 rounded-full" style="width: 50%"></div>
            </div>
        `;
    } else if (job.status === 'completed') {
        progressHTML = `
            <div class="w-full bg-green-200 rounded-full h-2">
                <div class="bg-green-500 h-2 rounded-full" style="width: 100%"></div>
            </div>
        `;
    }

    return `
        <tr>
            <td><strong>#${job.id}</strong></td>
            <td>#${job.flow_id}</td>
            <td><span class="${getStatusClass(job.status)}">${formatStatus(job.status)}</span></td>
            <td>${progressHTML}</td>
            <td>${formatDate(job.created_at)}</td>
            <td>${formatDate(job.updated_at)}</td>
            <td class="space-x-2">
                <button onclick="viewJobDetails(${job.id})" class="btn btn-sm btn-primary">
                    详情
                </button>
                ${job.status === 'pending' ? `
                    <button onclick="startJobNow(${job.id})" class="btn btn-sm btn-primary">
                        启动
                    </button>
                ` : ''}
                ${job.status === 'running' ? `
                    <button onclick="refreshJobStatus(${job.id})" class="btn btn-sm btn-secondary">
                        刷新
                    </button>
                ` : ''}
            </td>
        </tr>
    `;
}

async function viewJobDetails(jobId) {
    try {
        const response = await api.getJob(jobId);
        jobsData.selectedJob = response.data;

        const content = document.getElementById('jobDetailsContent');
        content.innerHTML = `
            <div class="space-y-4">
                <!-- Job Info -->
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <label class="text-sm text-gray-600">作业 ID</label>
                        <p class="font-semibold">#${jobsData.selectedJob.id}</p>
                    </div>
                    <div>
                        <label class="text-sm text-gray-600">流程 ID</label>
                        <p class="font-semibold">#${jobsData.selectedJob.flow_id}</p>
                    </div>
                    <div>
                        <label class="text-sm text-gray-600">状态</label>
                        <p><span class="${getStatusClass(jobsData.selectedJob.status)}">${formatStatus(jobsData.selectedJob.status)}</span></p>
                    </div>
                    <div>
                        <label class="text-sm text-gray-600">创建时间</label>
                        <p>${formatDate(jobsData.selectedJob.created_at)}</p>
                    </div>
                </div>

                <!-- Input Data -->
                <div>
                    <label class="text-sm text-gray-600 font-semibold">输入数据</label>
                    <pre class="bg-gray-50 p-3 rounded text-xs overflow-x-auto">${JSON.stringify(jobsData.selectedJob.input || {}, null, 2)}</pre>
                </div>

                <!-- Output Data -->
                ${jobsData.selectedJob.output ? `
                    <div>
                        <label class="text-sm text-gray-600 font-semibold">输出数据</label>
                        <pre class="bg-gray-50 p-3 rounded text-xs overflow-x-auto">${JSON.stringify(jobsData.selectedJob.output, null, 2)}</pre>
                    </div>
                ` : ''}

                <!-- Error -->
                ${jobsData.selectedJob.error ? `
                    <div>
                        <label class="text-sm text-red-600 font-semibold">错误信息</label>
                        <pre class="bg-red-50 p-3 rounded text-xs text-red-700">${jobsData.selectedJob.error}</pre>
                    </div>
                ` : ''}

                <!-- Job Tasks Timeline -->
                <div>
                    <label class="text-sm text-gray-600 font-semibold mb-2 block">任务执行时间线</label>
                    <div class="space-y-2">
                        <div class="flex items-center">
                            <div class="w-3 h-3 rounded-full bg-primary"></div>
                            <div class="ml-3 flex-1">
                                <p class="text-sm font-medium">作业创建</p>
                                <p class="text-xs text-gray-500">${formatDate(jobsData.selectedJob.created_at)}</p>
                            </div>
                        </div>
                        ${jobsData.selectedJob.started_at ? `
                            <div class="flex items-center">
                                <div class="w-3 h-3 rounded-full bg-blue-500"></div>
                                <div class="ml-3 flex-1">
                                    <p class="text-sm font-medium">作业启动</p>
                                    <p class="text-xs text-gray-500">${formatDate(jobsData.selectedJob.started_at)}</p>
                                </div>
                            </div>
                        ` : ''}
                        ${jobsData.selectedJob.completed_at ? `
                            <div class="flex items-center">
                                <div class="w-3 h-3 rounded-full bg-green-500"></div>
                                <div class="ml-3 flex-1">
                                    <p class="text-sm font-medium">作业完成</p>
                                    <p class="text-xs text-gray-500">${formatDate(jobsData.selectedJob.completed_at)}</p>
                                </div>
                            </div>
                        ` : ''}
                    </div>
                </div>
            </div>
        `;

        document.getElementById('jobDetailsModal').style.display = 'flex';
    } catch (error) {
        alert('获取作业详情失败: ' + error.message);
    }
}

function closeJobDetailsModal() {
    document.getElementById('jobDetailsModal').style.display = 'none';
}

async function startJobNow(jobId) {
    if (!confirm(`确定启动作业 #${jobId} 吗？`)) return;

    try {
        await api.startJob(jobId);
        alert('作业已启动！');
        await loadJobs();
    } catch (error) {
        alert('启动失败: ' + error.message);
    }
}

async function refreshJobStatus(jobId) {
    try {
        await loadJobs();
        alert('状态已刷新');
    } catch (error) {
        alert('刷新失败: ' + error.message);
    }
}

function toggleAutoRefresh(enabled) {
    jobsData.autoRefresh = enabled;

    if (enabled) {
        // Start auto-refresh interval
        if (!window.jobsRefreshInterval) {
            window.jobsRefreshInterval = setInterval(() => {
                if (jobsData.autoRefresh) {
                    loadJobs();
                }
            }, 5000); // Refresh every 5 seconds
        }
    } else {
        // Stop auto-refresh
        if (window.jobsRefreshInterval) {
            clearInterval(window.jobsRefreshInterval);
            window.jobsRefreshInterval = null;
        }
    }
}

function filterJobs(status) {
    // Simple filter - in real app, would re-render with filtered data
    alert(`筛选: ${status}\n\n此功能将在下一版本实现。`);
}
