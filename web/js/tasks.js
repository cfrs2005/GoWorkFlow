// Tasks Library Page Logic

let tasksData = {
    tasks: [],
    showModal: false,
    currentTask: null,
    isEditing: false,
};

async function loadTasks() {
    try {
        const response = await api.getTasks();
        tasksData.tasks = response.data || [];
        renderTasks();
    } catch (error) {
        console.error('Failed to load tasks:', error);
        throw error;
    }
}

function renderTasks() {
    const container = document.getElementById('tasks-content');

    container.innerHTML = `
        <!-- Header Actions -->
        <div class="flex justify-between items-center mb-6">
            <div>
                <p class="text-gray-600">管理可复用的任务定义</p>
            </div>
            <button onclick="openCreateTaskModal()" class="btn btn-primary">
                <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
                创建任务
            </button>
        </div>

        <!-- Task Type Filter -->
        <div class="flex space-x-2 mb-6">
            <button onclick="filterTasks('all')" class="btn btn-sm btn-primary">全部</button>
            <button onclick="filterTasks('manual')" class="btn btn-sm btn-secondary">手动任务</button>
            <button onclick="filterTasks('automated')" class="btn btn-sm btn-secondary">自动化任务</button>
            <button onclick="filterTasks('approval')" class="btn btn-sm btn-secondary">审批任务</button>
        </div>

        <!-- Tasks Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            ${tasksData.tasks.length > 0 ? tasksData.tasks.map(task => `
                <div class="card">
                    <div class="flex items-start justify-between mb-3">
                        <h3 class="text-lg font-semibold text-gray-800">${task.name}</h3>
                        <span class="status-badge ${task.is_active ? 'status-completed' : 'status-pending'}">
                            ${task.is_active ? '启用' : '禁用'}
                        </span>
                    </div>

                    <p class="text-sm text-gray-600 mb-3">${truncate(task.description, 100)}</p>

                    <div class="mb-4">
                        <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-primary-light text-white">
                            ${formatTaskType(task.task_type)}
                        </span>
                    </div>

                    <div class="text-xs text-gray-500 mb-4">
                        <div>创建时间: ${formatDate(task.created_at)}</div>
                        <div>更新时间: ${formatDate(task.updated_at)}</div>
                    </div>

                    <div class="flex space-x-2">
                        <button onclick="editTask(${task.id})" class="btn btn-sm btn-secondary flex-1">
                            编辑
                        </button>
                        <button onclick="deleteTask(${task.id})" class="btn btn-sm btn-danger flex-1">
                            删除
                        </button>
                    </div>
                </div>
            `).join('') : `
                <div class="col-span-full text-center py-12 text-gray-500">
                    暂无任务，请点击上方按钮创建
                </div>
            `}
        </div>

        <!-- Create/Edit Task Modal -->
        <div id="taskModal" class="modal-overlay" style="display: none;">
            <div class="modal-content">
                <div class="modal-header">
                    <h3 class="modal-title" id="taskModalTitle">创建任务</h3>
                </div>
                <div class="modal-body">
                    <form id="taskForm">
                        <div class="form-group">
                            <label class="form-label">任务名称</label>
                            <input type="text" id="taskName" class="form-input" required>
                        </div>
                        <div class="form-group">
                            <label class="form-label">描述</label>
                            <textarea id="taskDescription" class="form-textarea"></textarea>
                        </div>
                        <div class="form-group">
                            <label class="form-label">任务类型</label>
                            <select id="taskType" class="form-select">
                                <option value="manual">手动任务</option>
                                <option value="automated">自动化任务</option>
                                <option value="approval">审批任务</option>
                            </select>
                        </div>
                        <div class="form-group">
                            <label class="form-label">配置 (JSON格式)</label>
                            <textarea id="taskConfig" class="form-textarea" placeholder='{"key": "value"}'>{}</textarea>
                        </div>
                        <div class="form-group">
                            <label class="flex items-center">
                                <input type="checkbox" id="taskActive" class="mr-2" checked>
                                <span class="form-label mb-0">启用此任务</span>
                            </label>
                        </div>
                    </form>
                </div>
                <div class="modal-footer">
                    <button onclick="closeTaskModal()" class="btn btn-secondary">取消</button>
                    <button onclick="saveTask()" class="btn btn-primary">保存</button>
                </div>
            </div>
        </div>
    `;
}

function openCreateTaskModal() {
    document.getElementById('taskModal').style.display = 'flex';
    document.getElementById('taskModalTitle').textContent = '创建任务';
    document.getElementById('taskForm').reset();
    document.getElementById('taskConfig').value = '{}';
    tasksData.isEditing = false;
    tasksData.currentTask = null;
}

function closeTaskModal() {
    document.getElementById('taskModal').style.display = 'none';
}

async function saveTask() {
    const name = document.getElementById('taskName').value;
    const description = document.getElementById('taskDescription').value;
    const taskType = document.getElementById('taskType').value;
    const configStr = document.getElementById('taskConfig').value;
    const isActive = document.getElementById('taskActive').checked;

    if (!name) {
        alert('请输入任务名称');
        return;
    }

    let config = {};
    try {
        config = JSON.parse(configStr);
    } catch (e) {
        alert('配置格式错误，请输入有效的JSON');
        return;
    }

    try {
        const data = {
            name,
            description,
            task_type: taskType,
            config,
            is_active: isActive,
        };

        if (tasksData.isEditing && tasksData.currentTask) {
            data.id = tasksData.currentTask.id;
            await api.updateTask(data);
        } else {
            await api.createTask(data);
        }

        closeTaskModal();
        await loadTasks();
        alert('任务保存成功！');
    } catch (error) {
        alert('保存失败: ' + error.message);
    }
}

async function editTask(id) {
    try {
        const response = await api.getTask(id);
        const task = response.data;

        tasksData.currentTask = task;
        tasksData.isEditing = true;

        document.getElementById('taskName').value = task.name;
        document.getElementById('taskDescription').value = task.description;
        document.getElementById('taskType').value = task.task_type;
        document.getElementById('taskConfig').value = JSON.stringify(task.config || {}, null, 2);
        document.getElementById('taskActive').checked = task.is_active;

        document.getElementById('taskModal').style.display = 'flex';
        document.getElementById('taskModalTitle').textContent = '编辑任务';
    } catch (error) {
        alert('获取任务信息失败: ' + error.message);
    }
}

async function deleteTask(id) {
    if (!confirm('确定要删除此任务吗？')) return;

    try {
        await api.deleteTask(id);
        await loadTasks();
        alert('删除成功！');
    } catch (error) {
        alert('删除失败: ' + error.message);
    }
}

function filterTasks(type) {
    // Simple filter - in real app, would re-render with filtered data
    alert(`筛选: ${type}\n\n此功能将在下一版本实现。`);
}
