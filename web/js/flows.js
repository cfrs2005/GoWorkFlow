// Flows Management Page Logic

let flowsData = {
    flows: [],
    flowTasks: {},
    showModal: false,
    currentFlow: null,
    isEditing: false,
};

async function loadFlows() {
    try {
        const response = await api.getFlows();
        flowsData.flows = response.data || [];
        renderFlows();
    } catch (error) {
        console.error('Failed to load flows:', error);
        throw error;
    }
}

function renderFlows() {
    const container = document.getElementById('flows-content');

    container.innerHTML = `
        <!-- Header Actions -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <p class="text-gray-600">管理工作流程定义和任务编排</p>
            </div>
            <button onclick="openCreateFlowModal()" class="btn btn-primary">
                <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
                创建流程
            </button>
        </div>

        <!-- Quick Templates -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
            <!-- YouTube Analysis Template (Featured) -->
            <div class="card border-2 border-red-500 cursor-pointer hover:shadow-lg" onclick="createYouTubeAnalysisFlow()">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 bg-red-500 rounded-lg flex items-center justify-center text-white mr-4">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg mb-1">🔥 YouTube 视频智能分析</h3>
                        <p class="text-sm text-gray-600">AI 深度分析：字幕提取、内容总结、思维导图、重点分析</p>
                        <p class="text-xs text-red-500 mt-2 font-semibold">✨ 真实可用！点击创建</p>
                    </div>
                </div>
            </div>

            <div class="card border-2 border-gray-300 cursor-pointer hover:shadow-lg opacity-60" onclick="alert('JIRA 流程暂未实现，请使用 YouTube 分析流程')">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 bg-gray-400 rounded-lg flex items-center justify-center text-white mr-4">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg mb-1 text-gray-500">JIRA 数据采集流程</h3>
                        <p class="text-sm text-gray-500">自动获取 JIRA 数据、解析下载链接、处理并推送</p>
                        <p class="text-xs text-gray-400 mt-2">敬请期待</p>
                    </div>
                </div>
            </div>

            <div class="card border-2 border-gray-300 cursor-pointer hover:shadow-lg opacity-60" onclick="alert('RobotSN 流程暂未实现，请使用 YouTube 分析流程')">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 bg-gray-400 rounded-lg flex items-center justify-center text-white mr-4">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg mb-1 text-gray-500">RobotSN 数据分析流程</h3>
                        <p class="text-sm text-gray-500">登录服务、获取机器人信息、生成分析报告</p>
                        <p class="text-xs text-gray-400 mt-2">敬请期待</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Flows List -->
        <div class="card">
            <div class="card-header">所有流程</div>
            <div class="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>名称</th>
                            <th>描述</th>
                            <th>版本</th>
                            <th>状态</th>
                            <th>创建时间</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${flowsData.flows.length > 0 ? flowsData.flows.map(flow => `
                            <tr>
                                <td><strong>#${flow.id}</strong></td>
                                <td>${flow.name}</td>
                                <td>${truncate(flow.description, 60)}</td>
                                <td><span class="px-2 py-1 bg-gray-100 rounded text-sm">${flow.version}</span></td>
                                <td>
                                    <span class="status-badge ${flow.is_active ? 'status-completed' : 'status-pending'}">
                                        ${flow.is_active ? '启用' : '禁用'}
                                    </span>
                                </td>
                                <td>${formatDate(flow.created_at)}</td>
                                <td class="space-x-2">
                                    <button onclick="viewFlowTasks(${flow.id})" class="btn btn-sm btn-secondary">
                                        查看任务
                                    </button>
                                    <button onclick="runFlow(${flow.id})" class="btn btn-sm btn-primary">
                                        运行
                                    </button>
                                    <button onclick="editFlow(${flow.id})" class="btn btn-sm btn-secondary">
                                        编辑
                                    </button>
                                    <button onclick="deleteFlow(${flow.id})" class="btn btn-sm btn-danger">
                                        删除
                                    </button>
                                </td>
                            </tr>
                        `).join('') : `
                            <tr>
                                <td colspan="7" class="text-center text-gray-500 py-8">
                                    暂无流程，请点击上方按钮创建
                                </td>
                            </tr>
                        `}
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Create/Edit Flow Modal -->
        <div id="flowModal" class="modal-overlay" style="display: none;">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title" id="modalTitle">创建流程</h3>
                </div>
                <div class="modal-body">
                    <form id="flowForm">
                        <div class="form-group">
                            <label class="form-label">流程名称</label>
                            <input type="text" id="flowName" class="form-input" required>
                        </div>
                        <div class="form-group">
                            <label class="form-label">描述</label>
                            <textarea id="flowDescription" class="form-textarea"></textarea>
                        </div>
                        <div class="form-group">
                            <label class="form-label">版本</label>
                            <input type="text" id="flowVersion" class="form-input" value="1.0.0">
                        </div>
                        <div class="form-group">
                            <label class="flex items-center">
                                <input type="checkbox" id="flowActive" class="mr-2" checked>
                                <span class="form-label mb-0">启用此流程</span>
                            </label>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button onclick="closeFlowModal()" class="btn btn-secondary">取消</button>
                    <button onclick="saveFlow()" class="btn btn-primary">保存</button>
                </div>
            </div>
        </div>
    `;
}

function openCreateFlowModal() {
    document.getElementById('flowModal').style.display = 'flex';
    document.getElementById('modalTitle').textContent = '创建流程';
    document.getElementById('flowForm').reset();
    flowsData.isEditing = false;
    flowsData.currentFlow = null;
}

function closeFlowModal() {
    document.getElementById('flowModal').style.display = 'none';
}

async function saveFlow() {
    const name = document.getElementById('flowName').value;
    const description = document.getElementById('flowDescription').value;
    const version = document.getElementById('flowVersion').value;
    const isActive = document.getElementById('flowActive').checked;

    if (!name) {
        alert('请输入流程名称');
        return;
    }

    try {
        const data = {
            name,
            description,
            version,
            is_active: isActive,
            created_by: 1, // Default user
        };

        if (flowsData.isEditing && flowsData.currentFlow) {
            data.id = flowsData.currentFlow.id;
            await api.updateFlow(data);
        } else {
            await api.createFlow(data);
        }

        closeFlowModal();
        await loadFlows();
        alert('流程保存成功！');
    } catch (error) {
        alert('保存失败: ' + error.message);
    }
}

async function editFlow(id) {
    try {
        const response = await api.getFlow(id);
        const flow = response.data;

        flowsData.currentFlow = flow;
        flowsData.isEditing = true;

        document.getElementById('flowName').value = flow.name;
        document.getElementById('flowDescription').value = flow.description;
        document.getElementById('flowVersion').value = flow.version;
        document.getElementById('flowActive').checked = flow.is_active;

        document.getElementById('flowModal').style.display = 'flex';
        document.getElementById('modalTitle').textContent = '编辑流程';
    } catch (error) {
        alert('获取流程信息失败: ' + error.message);
    }
}

async function deleteFlow(id) {
    if (!confirm('确定要删除此流程吗？')) return;

    try {
        await api.deleteFlow(id);
        await loadFlows();
        alert('删除成功！');
    } catch (error) {
        alert('删除失败: ' + error.message);
    }
}

async function runFlow(flowId) {
    try {
        const input = prompt('请输入流程输入参数（JSON格式，可选）:', '{}');
        let inputData = {};

        if (input && input.trim()) {
            try {
                inputData = JSON.parse(input);
            } catch (e) {
                alert('输入参数格式错误');
                return;
            }
        }

        const job = await api.createJob({
            flow_id: flowId,
            input: inputData,
        });

        await api.startJob(job.data.id);
        alert(`作业已启动！作业ID: ${job.data.id}`);

        // Navigate to jobs page
        if (window.appData) {
            window.appData.currentPage = 'jobs';
        }
    } catch (error) {
        alert('启动失败: ' + error.message);
    }
}

function viewFlowTasks(flowId) {
    alert(`查看流程 #${flowId} 的任务列表\n\n此功能将在下一版本实现。`);
}

// Quick template creation functions

async function createYouTubeAnalysisFlow() {
    // 检查是否已存在
    const existing = flowsData.flows.find(f => f.name === 'YouTube 视频智能分析');
    if (existing) {
        if (confirm('YouTube 分析流程已存在（ID: ' + existing.id + '）\n\n是否直接运行该流程？')) {
            runYouTubeAnalysis(existing.id);
        }
        return;
    }

    if (!confirm('🎬 创建 YouTube 视频智能分析流程？\n\n包含以下功能：\n• YouTube 字幕提取\n• AI 深度分析（BigModel GLM-4-Air）\n• 生成精美 HTML 报告\n\n是否继续？')) return;

    alert('⚠️ 提示：\n\n1. 首次运行将使用模拟数据（演示效果）\n2. 如需真实分析，请设置环境变量：\n   BIGMODEL_API_KEY=你的密钥\n3. 报告将保存在 ./reports 目录\n\n点击确定继续创建...');

    try {
        const flow = await api.createFlow({
            name: 'YouTube 视频智能分析',
            description: 'AI 驱动的 YouTube 视频深度分析：自动提取字幕、生成阅读摘要、思维导图、重点分析和个人认知，最终生成精美的 HTML 分析报告。',
            version: '1.0.0',
            is_active: true,
            created_by: 1,
        });

        alert('✅ 流程创建成功！\n\n流程 ID: ' + flow.data.id + '\n\n💡 提示：数据库迁移可能需要手动执行：\nmysql -u root -p workflow < migrations/004_youtube_analysis_workflow.sql');

        await loadFlows();

        // 询问是否立即运行
        if (confirm('是否立即运行 YouTube 分析流程？')) {
            runYouTubeAnalysis(flow.data.id);
        }
    } catch (error) {
        alert('创建失败: ' + error.message + '\n\n💡 可能需要先运行数据库迁移：\nmysql -u root -p workflow < migrations/004_youtube_analysis_workflow.sql');
    }
}

async function runYouTubeAnalysis(flowId) {
    // 弹出自定义输入对话框
    const videoURL = prompt('🎬 请输入 YouTube 视频地址：\n\n支持格式：\n• https://www.youtube.com/watch?v=VIDEO_ID\n• https://youtu.be/VIDEO_ID\n• VIDEO_ID（11位字符）\n\n示例：\nhttps://www.youtube.com/watch?v=dQw4w9WgXcQ', 'https://www.youtube.com/watch?v=dQw4w9WgXcQ');

    if (!videoURL || videoURL.trim() === '') {
        alert('❌ 未输入视频地址，操作已取消');
        return;
    }

    try {
        // 创建作业
        const job = await api.createJob({
            flow_id: flowId,
            input: {
                video_url: videoURL.trim(),
                language: 'en' // 可以改为 'zh' 获取中文字幕
            },
        });

        // 自动执行作业
        const response = await fetch('/api/jobs/auto-execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ job_id: job.data.id }),
        });

        if (!response.ok) {
            throw new Error('启动自动执行失败');
        }

        alert(`✅ YouTube 分析已启动！\n\n作业 ID: ${job.data.id}\n视频: ${videoURL}\n\n⏱️ 预计耗时：1-3分钟\n📊 请在"作业监控"页面查看进度\n📄 完成后可在作业详情查看报告链接`);

        // 跳转到作业监控页面
        if (window.appData) {
            window.appData.currentPage = 'jobs';
        }
    } catch (error) {
        alert('❌ 启动失败: ' + error.message);
    }
}

async function createJiraFlow() {
    alert('JIRA 流程暂未实现\n\n请使用 YouTube 视频智能分析流程体验完整功能！');
}

async function createRobotSNFlow() {
    alert('RobotSN 流程暂未实现\n\n请使用 YouTube 视频智能分析流程体验完整功能！');
}
