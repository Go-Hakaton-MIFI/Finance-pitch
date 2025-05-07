import axios from 'axios';

const API_URL = '/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Добавляем токен к каждому запросу
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Обработка ошибок
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error.response?.data || error.message);
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    // Добавляем дополнительную информацию об ошибке
    if (error.response?.data?.error) {
      error.message = error.response.data.error;
    }
    return Promise.reject(error);
  }
);

// Auth API
const authAPI = {
  login: (credentials) => api.post('/login', credentials),
  register: (userData) => api.post('/registration', userData),
  getSubjectTypes: () => api.get('/subject_types')
};

// Transactions API
const transactionsAPI = {
  getAll: (filters = {}) => api.post('/transactions/filter', filters),
  getById: (id) => api.get(`/transactions/${id}`),
  create: (data) => api.post('/transactions', data),
  prepare: (data) => api.post('/transactions/prepared', data),
  delete: (id) => api.delete(`/transactions/${id}`),
  getCategories: () => api.get('/categories'),
  getStatuses: () => api.get('/trans_statuses'),
  getBanks: () => api.get('/banks')
};

// Categories API
const categoriesAPI = {
  getAll: () => api.get('/categories')
};

// Statuses API
const statusesAPI = {
  getAll: () => api.get('/trans_statuses')
};

// Banks API
const banksAPI = {
  getAll: () => api.get('/banks')
};

// Analytics API
const analyticsAPI = {
  getDynamicsByPeriod: (period, filter) => 
    api.post(`/analytics/dynamics/by-period?period=${period}`, {
      date: {
        from: filter.date.from,
        to: filter.date.to
      }
    }),
  getDynamicsByType: (type, filter) => 
    api.post(`/analytics/dynamics/by-type?trans_type=${type}`, {
      ...filter,
      user_login: localStorage.getItem('user_login')
    }),
  getIncomeExpenseComparison: (filter) =>
    api.post('/analytics/compare-income-expense', {
      date: {
        from: filter.date.from,
        to: filter.date.to
      }
    }).catch(error => {
      console.error('Error in getIncomeExpenseComparison:', error.response?.data || error);
      throw error;
    }),
  getStatusSummary: (filter) =>
    api.post('/analytics/status-summary', {
      ...filter,
      user_login: localStorage.getItem('user_login')
    }),
  getBanksSummary: (filter) =>
    api.post('/analytics/banks-summary', {
      ...filter,
      user_login: localStorage.getItem('user_login')
    }),
  getCategoriesSummary: (transType, filter) =>
    api.post(`/analytics/categories-summary?trans_type=${transType}`, {
      date: {
        from: filter.date.from,
        to: filter.date.to
      }
    }),
  generateBanksReport: (filter) =>
    api.post('/analytics/banks-summary/report', {
      ...filter,
      user_login: localStorage.getItem('user_login')
    }),
  downloadReport: (reportId) =>
    api.get(`/analytics/banks-summary/report/${reportId}`, { responseType: 'blob' })
};

// Reports API
const reportsAPI = {
  generateBanksReport: (filters) => 
    api.post('/reports/banks', filters),
  generateCategoriesReport: (filters) => 
    api.post('/reports/categories', filters),
  generatePeriodReport: (filters) => 
    api.post('/reports/period', filters),
  downloadReport: (reportId) => 
    api.get(`/reports/${reportId}`, { responseType: 'blob' })
};

// Settings API
const settingsAPI = {
  getSettings: () => api.get('/settings'),
  updateSettings: (settings) => api.put('/settings', settings)
};

// Profile API
const profileAPI = {
  getProfile: () => api.get('/profile'),
  updateProfile: (data) => api.put('/profile', data)
};

// Participants API
const participantsAPI = {
  getAll: () => api.get('/participants'),
  create: (data) => api.post('/participants', data)
};

export {
  authAPI,
  transactionsAPI,
  categoriesAPI,
  statusesAPI,
  banksAPI,
  analyticsAPI,
  reportsAPI,
  settingsAPI,
  profileAPI,
  participantsAPI
};

export default api;