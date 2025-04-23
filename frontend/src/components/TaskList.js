import React, { useState, useEffect, useCallback } from 'react';
import TaskForm from './TaskForm';
import { getTasks, deleteTask } from '../api/taskService';

function TaskList({ token }) {
  const [tasks, setTasks] = useState([]);
  const [editingTask, setEditingTask] = useState(null);
  const [error, setError] = useState('');

  const fetchTasks = useCallback(async () => {
    try {
      const res = await getTasks({}, token);
      setTasks(res.tasks);
      setError('');
    } catch (err) {
      setError(err.message || 'Failed to fetch tasks');
      console.error('Error fetching tasks:', err);
    }
  }, [token]);

  const handleDelete = async (id) => {
    try {
      await deleteTask(id, token);
      fetchTasks();
      setError('');
    } catch (err) {
      setError(err.message || 'Error deleting task');
      console.error('Task deletion error:', err);
    }
  };

  const handleEdit = (task) => {
    setEditingTask(task);
  };

  const handleCancelEdit = () => {
    setEditingTask(null);
  };

  useEffect(() => {
    if (token) {
      fetchTasks();
    }
  }, [token, fetchTasks]);

  return (
    <div>
      <h2>Tasks</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <TaskForm task={editingTask} onCancel={handleCancelEdit} onSuccess={fetchTasks} token={token} />
      <ul>
        {tasks.map((task) => (
          <li key={task.id}>
            {task.title} ({task.status}) - {task.due_date ? new Date(task.due_date).toLocaleDateString() : 'No due date'}
            <button onClick={() => handleEdit(task)}>Edit</button>
            <button onClick={() => handleDelete(task.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default TaskList;