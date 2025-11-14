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
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div class="card border-2 border-primary cursor-pointer hover:shadow-lg" onclick="createJiraFlow()">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 bg-primary-light rounded-lg flex items-center justify-center text-white mr-4">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg mb-1">JIRA 数据采集流程</h3>
                        <p class="text-sm text-gray-600">自动获取 JIRA 数据、解析下载链接、处理并推送</p>
                        <p class="text-xs text-primary mt-2">点击快速创建</p>
                    </div>
                </div>
            </div>

            <div class="card border-2 border-blue-500 cursor-pointer hover:shadow-lg" onclick="createRobotSNFlow()">
                <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 bg-blue-500 rounded-lg flex items-center justify-center text-white mr-4">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                        </svg>
                    </div>
                    <div>
                        <h3 class="font-semibold text-lg mb-1">RobotSN 数据分析流程</h3>
                        <p class="text-sm text-gray-600">登录服务、获取机器人信息、生成分析报告</p>
                        <p class="text-xs text-blue-500 mt-2">点击快速创建</p>
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

async function createJiraFlow() {
    if (!confirm('创建 JIRA 数据采集流程？')) return;

    try {
        const flow = await api.createFlow({
            name: 'JIRA 数据采集流程',
            description: '自动获取 JIRA 页面内容，解析下载链接，批量下载并推送数据',
            version: '1.0.0',
            is_active: true,
            created_by: 1,
        });

        alert('流程创建成功！流程ID: ' + flow.data.id + '\n\n请在任务库中配置具体任务。');
        await loadFlows();
    } catch (error) {
        alert('创建失败: ' + error.message);
    }
}

async function createRobotSNFlow() {
    if (!confirm('创建 RobotSN 数据分析流程？')) return;

    try {
        const flow = await api.createFlow({
            name: 'RobotSN 数据分析流程',
            description: '登录服务，获取 RobotSN 最新信息，生成数据分析报告',
            version: '1.0.0',
            is_active: true,
            created_by: 1,
        });

        alert('流程创建成功！流程ID: ' + flow.data.id + '\n\n请在任务库中配置具体任务。');
        await loadFlows();
    } catch (error) {
        alert('创建失败: ' + error.message);
    }
}
