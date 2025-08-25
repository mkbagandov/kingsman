import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux'; // Import useSelector
import { FaUserCircle, FaBell, FaShoppingCart } from 'react-icons/fa'; // Import FaBell and FaShoppingCart icons
import './Navbar.css'; 

function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const cartItemsCount = useSelector((state) => state.cart.cartItems.length); // Get cart items count from Redux store

  useEffect(() => {
    const checkAuth = () => {
      const token = localStorage.getItem('jwtToken');
      setIsAuthenticated(!!token);
    };

    checkAuth(); 

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
        </div>
      </div>
      <div className="nav-right">
        <div className="auth-links">
          {/* Cart Icon */}
          <Link to="/cart" className="cart-icon">
            <FaShoppingCart size={24} color="var(--text-light)" />
            {cartItemsCount > 0 && <span className="cart-badge">{cartItemsCount}</span>}
          </Link>

          {/* Notifications Icon */}
          <Link to="/notifications" className="notification-icon">
            <FaBell size={24} color="var(--text-light)" />
          </Link>

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
