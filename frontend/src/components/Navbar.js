import React from 'react';
import { Link } from 'react-router-dom';
import './Navbar.css'; // Import the new CSS file

function Navbar() {
  return (
    <nav className="navbar">
      <div className="nav-left">
        <Link to="/" className="navbar-logo">MR.KINGSMAN</Link>
        <div className="navbar-links">
          <Link to="/products">Продукты</Link>
          <Link to="/stores">Магазины</Link>
          <Link to="/users/profile">Профиль</Link>
          <Link to="/notifications">Уведомления</Link>
        </div>
      </div>
      <div className="nav-right">
        <div className="auth-links">
          <Link to="/login">Войти</Link>
          <Link to="/register">Регистрация</Link>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
