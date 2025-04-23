import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Signup from './components/Signup';
import TaskList from './components/TaskList';
import './App.css';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token') || '');

  const handleLogout = () => {
    setToken('');
    localStorage.removeItem('token');
  };

  return (
    <Router>
      <div className="App">
        <h1>GoTasker</h1>
        <Routes>
          {token ? (
            <>
              <Route
                path="/tasks"
                element={
                  <div>
                    <button onClick={handleLogout}>Logout</button>
                    <TaskList token={token} />
                  </div>
                }
              />
              <Route path="*" element={<Navigate to="/tasks" />} />
            </>
          ) : (
            <>
              <Route path="/signup" element={<Signup setToken={setToken} />} />
              <Route path="/login" element={<Login setToken={setToken} />} />
              <Route path="*" element={<Navigate to="/signup" />} />
            </>
          )}
        </Routes>
      </div>
    </Router>
  );
}

export default App;