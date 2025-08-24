import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FaUserCircle } from 'react-icons/fa'; // Import a user icon
import './Navbar.css'; // Import the new CSS file

function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const checkAuth = () => {
      const token = localStorage.getItem('jwtToken');
      setIsAuthenticated(!!token);
    };

    checkAuth(); // Check auth status on component mount

    // Listen for changes in localStorage across tabs/windows
    window.addEventListener('storage', checkAuth);

    return () => {
      window.removeEventListener('storage', checkAuth);
    };
  }, []);

  return (
    <nav className="navbar">
      <div className="nav-left">
        <Link to="/" className="navbar-logo">MR.KINGSMAN</Link>
        <div className="navbar-links">
          <Link to="/">Главная</Link>
          <Link to="/products">Каталог</Link>
          <Link to="/about">О нас</Link>
          <Link to="/contacts">Контакты</Link>
          <Link to="/notifications">Уведомления</Link>
        </div>
      </div>
      <div className="nav-right">
        <div className="auth-links">
          {isAuthenticated ? (
            <Link to="/users/profile" className="profile-icon">
              <FaUserCircle size={30} color="var(--text-light)" />
            </Link>
          ) : (
            <Link to="/login" className="btn-login">Войти</Link>
          )}
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
