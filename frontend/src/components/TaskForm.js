import React, { useState, useEffect } from 'react';
import { createTask, updateTask } from '../api/taskService';

function TaskForm({ task, onCancel, onSuccess, token }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('Pending');
  const [dueDate, setDueDate] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (task) {
      setTitle(task.title || '');
      setDescription(task.description || '');
      setStatus(task.status || 'Pending');
      setDueDate(task.due_date ? new Date(task.due_date).toISOString().slice(0, 16) : '');
    } else {
      setTitle('');
      setDescription('');
      setStatus('Pending');
      setDueDate('');
    }
  }, [task]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!title) {
      setError('Title is required');
      return;
    }
    if (!['Pending', 'In Progress', 'Completed'].includes(status)) {
      setError('Invalid status');
      return;
    }

    const taskData = {
      title,
      description: description || '',
      status,
      due_date: dueDate ? new Date(dueDate).toISOString() : null,
    };

    try {
      if (task) {
        await updateTask(task.id, taskData, token);
      } else {
        await createTask(taskData, token);
      }
      setError('');
      setTitle('');
      setDescription('');
      setStatus('Pending');
      setDueDate('');
      if (onSuccess) onSuccess();
      if (task && onCancel) onCancel();
    } catch (err) {
      setError(err.message || (task ? 'Error updating task' : 'Error creating task'));
      console.error(task ? 'Task update error:' : 'Task creation error:', err);
    }
  };

  return (
    <div>
      <h2>{task ? 'Edit Task' : 'Create Task'}</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        <textarea
          placeholder="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <select value={status} onChange={(e) => setStatus(e.target.value)} required>
          <option value="Pending">Pending</option>
          <option value="In Progress">In Progress</option>
          <option value="Completed">Completed</option>
        </select>
        <input
          type="datetime-local"
          value={dueDate}
          onChange={(e) => setDueDate(e.target.value)}
        />
        <button type="submit">{task ? 'Update Task' : 'Add Task'}</button>
        {task && (
          <button type="button" onClick={onCancel}>
            Cancel
          </button>
        )}
      </form>
    </div>
  );
}

export default TaskForm;