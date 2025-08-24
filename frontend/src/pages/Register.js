import React, { useState } from 'react';
import { registerUser } from '../api/api';
import { useNavigate } from 'react-router-dom'; // Import useNavigate

function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate(); // Initialize useNavigate

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await registerUser({ username, email, password, phoneNumber });
      setMessage(`Registration successful: ${response.data.message}`);
      navigate('/login'); // Redirect to login page after successful registration
    } catch (error) {
      setMessage(`Registration failed: ${error.response?.data?.error || error.message}`);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-form-container">
        <h1>Регистрация</h1>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="register-username">Имя пользователя:</label>
            <input type="text" id="register-username" value={username} onChange={(e) => setUsername(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="register-email">Email:</label>
            <input type="email" id="register-email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="register-password">Пароль:</label>
            <input type="password" id="register-password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="register-phone">Номер телефона:</label>
            <input type="tel" id="register-phone" value={phoneNumber} onChange={(e) => setPhoneNumber(e.target.value)} />
          </div>
          <button type="submit">Зарегистрироваться</button>
        </form>
        {message && <p className="message">{message.includes('Registration successful') ? message.replace('Registration successful', 'Регистрация успешна') : message.replace('Registration failed', 'Ошибка регистрации')}</p>}
      </div>
    </div>
  );
}

export default Register;
