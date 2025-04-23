import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const createTask = async (taskData, token) => {
  try {
    const response = await axios.post(`${API_URL}/tasks`, {
      title: taskData.title,
      description: taskData.description || '',
      status: taskData.status || 'Pending',
      due_date: taskData.due_date || null,
    }, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to create task';
    console.error('Error creating task:', error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const getTasks = async (params = {}, token) => {
  try {
    const response = await axios.get(`${API_URL}/tasks`, {
      params,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to fetch tasks';
    console.error('Error fetching tasks:', error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const getTaskById = async (id, token) => {
  try {
    const response = await axios.get(`${API_URL}/tasks/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to fetch task';
    console.error(`Error fetching task ${id}:`, error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const updateTask = async (id, taskData, token) => {
  try {
    const response = await axios.put(`${API_URL}/tasks/${id}`, {
      title: taskData.title,
      description: taskData.description || '',
      status: taskData.status || 'Pending',
      due_date: taskData.due_date || null,
    }, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to update task';
    console.error(`Error updating task ${id}:`, error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const deleteTask = async (id, token) => {
  try {
    const response = await axios.delete(`${API_URL}/tasks/${id}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to delete task';
    console.error(`Error deleting task ${id}:`, error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const signup = async (credentials) => {
  try {
    const response = await axios.post(`${API_URL}/register`, credentials, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to sign up';
    console.error('Error signing up:', error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};

export const login = async (credentials) => {
  try {
    const response = await axios.post(`${API_URL}/login`, credentials, {
      headers: {
        'Content-Type': 'application/json',
      },
    });
    return response.data;
  } catch (error) {
    const errorMessage = error.response?.data?.error || 'Failed to log in';
    console.error('Error logging in:', error.response?.data || error.message);
    throw new Error(errorMessage);
  }
};