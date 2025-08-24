import React, { useState } from 'react';
import { loginUser } from '../api/api';
import { useNavigate } from 'react-router-dom';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await loginUser({ email, password });
      setMessage(`Login successful: ${response.data.message}`);
      navigate('/');
    } catch (error) {
      setMessage(`Login failed: ${error.response?.data?.error || error.message}`);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-form-container">
        <h1>Вход</h1>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="login-email">Email:</label>
            <input type="email" id="login-email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="login-password">Пароль:</label>
            <input type="password" id="login-password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <button type="submit">Войти</button>
        </form>
        {message && <p className="message">{message.includes('Login successful') ? message.replace('Login successful', 'Вход выполнен успешно') : message.replace('Login failed', 'Ошибка входа')}</p>}
        <p style={{ marginTop: '20px' }}>
          Ещё нет аккаунта? <a href="/register" style={{ color: 'var(--primary-dark-blue)', textDecoration: 'none', fontWeight: 'bold' }}>Зарегистрироваться</a>
        </p>
      </div>
    </div>
  );
}

export default Login;
