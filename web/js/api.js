// API Client for GoWorkFlow

const API_BASE = '/api';

class APIClient {
    async request(method, url, data = null) {
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
        };

        if (data) {
            options.body = JSON.stringify(data);
        }

        try {
            const response = await fetch(`${API_BASE}${url}`, options);

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'API request failed');
            }

            return await response.json();
        } catch (error) {
            console.error('API Error:', error);
            throw error;
        }
    }

    // Tasks
    async getTasks() {
        return this.request('GET', '/tasks');
    }

    async getTask(id) {
        return this.request('GET', `/tasks?id=${id}`);
    }

    async createTask(data) {
        return this.request('POST', '/tasks', data);
    }

    async updateTask(data) {
        return this.request('PUT', '/tasks', data);
    }

    async deleteTask(id) {
        return this.request('DELETE', `/tasks?id=${id}`);
    }

    // Flows
    async getFlows() {
        return this.request('GET', '/flows');
    }

    async getFlow(id) {
        return this.request('GET', `/flows?id=${id}`);
    }

    async createFlow(data) {
        return this.request('POST', '/flows', data);
    }

    async updateFlow(data) {
        return this.request('PUT', '/flows', data);
    }

    async deleteFlow(id) {
        return this.request('DELETE', `/flows?id=${id}`);
    }

    // Jobs
    async getJobs() {
        return this.request('GET', '/jobs');
    }

    async getJob(id) {
        return this.request('GET', `/jobs?id=${id}`);
    }

    async createJob(data) {
        return this.request('POST', '/jobs', data);
    }

    async startJob(id) {
        return this.request('POST', '/jobs/start', { job_id: id });
    }

    async getNextTask(jobId) {
        return this.request('GET', `/jobs/next-task?job_id=${jobId}`);
    }

    // Job Tasks
    async startTask(data) {
        return this.request('POST', '/tasks/start', data);
    }

    async completeTask(data) {
        return this.request('POST', '/tasks/complete', data);
    }

    async failTask(data) {
        return this.request('POST', '/tasks/fail', data);
    }

    async skipTask(data) {
        return this.request('POST', '/tasks/skip', data);
    }

    async rollbackTask(data) {
        return this.request('POST', '/tasks/rollback', data);
    }

    // Job Context (新增)
    async getJobContext(jobId) {
        return this.request('GET', `/jobs/${jobId}/context`);
    }

    async updateJobContext(jobId, context) {
        return this.request('PUT', `/jobs/${jobId}/context`, context);
    }
}

// Global API instance
const api = new APIClient();
