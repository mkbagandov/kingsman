import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { FaUserCircle, FaBell, FaShoppingCart, FaBars, FaTimes } from 'react-icons/fa';
import './Navbar.css';

function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const cartItemsCount = useSelector((state) => state.cart.cartItems.length);

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

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  // Handle clicks outside the mobile menu
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (isMobileMenuOpen && !event.target.closest('.mobile-drawer') && !event.target.closest('.burger-menu-btn')) {
        setIsMobileMenuOpen(false);
      }
    };

    if (isMobileMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside);
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isMobileMenuOpen]);

  return (
    <>
      <nav className="navbar">
        {/* Mobile Burger Menu Button */}
        <button
          className="burger-menu-btn"
          onClick={toggleMobileMenu}
          aria-label="Toggle Mobile Menu"
        >
          {isMobileMenuOpen ? <FaTimes /> : <FaBars />}
        </button>

        {/* Left side: Logo and Navigation Links */}
        <div className="nav-left">
          {/* Logo */}
          <Link to="/" className="navbar-logo">
            MR.KINGSMAN
          </Link>

          {/* Desktop Navigation Links */}
          <div className="navbar-links desktop-only">
            <Link to="/">Главная</Link>
            <Link to="/products">Каталог</Link>
            <Link to="/about">О нас</Link>
            <Link to="/contacts">Контакты</Link>
          </div>
        </div>

        {/* Right-aligned Icons/Auth */}
        <div className="nav-right">
          <Link to="/cart" className="cart-icon">
            <FaShoppingCart size={24} />
            {cartItemsCount > 0 && (
              <span className="cart-badge">
                {cartItemsCount}
              </span>
            )}
          </Link>
          <Link to="/notifications" className="notification-icon">
            <FaBell size={24} />
          </Link>
          {isAuthenticated ? (
            <Link to="/users/profile" className="profile-icon">
              <FaUserCircle size={30} />
            </Link>
          ) : (
            <Link to="/login" className="btn-login">
              Войти
            </Link>
          )}
        </div>
      </nav>

      {/* Mobile Sidebar Menu */}
      <div className={`mobile-drawer ${isMobileMenuOpen ? 'open' : ''}`}>
        <div className="drawer-content">
          <div className="drawer-header">
            <Link to="/" className="drawer-logo" onClick={toggleMobileMenu}>MR.KINGSMAN</Link>
            <button
              className="close-drawer-btn"
              onClick={toggleMobileMenu}
              aria-label="Close Mobile Menu"
            >
              <FaTimes />
            </button>
          </div>
          <div className="drawer-links">
            <Link to="/" onClick={toggleMobileMenu}>Главная</Link>
            <Link to="/products" onClick={toggleMobileMenu}>Каталог</Link>
            <Link to="/about" onClick={toggleMobileMenu}>О нас</Link>
            <Link to="/contacts" onClick={toggleMobileMenu}>Контакты</Link>
            {isAuthenticated ? (
              <Link to="/users/profile" className="drawer-link-with-icon" onClick={toggleMobileMenu}>
                <FaUserCircle size={20} />
                <span>Профиль</span>
              </Link>
            ) : (
              <Link to="/login" className="drawer-login-btn" onClick={toggleMobileMenu}>
                Войти
              </Link>
            )}
            <Link to="/cart" className="drawer-link-with-icon" onClick={toggleMobileMenu}>
              <FaShoppingCart size={20} />
              <span>Корзина ({cartItemsCount})</span>
            </Link>
            <Link to="/notifications" className="drawer-link-with-icon" onClick={toggleMobileMenu}>
              <FaBell size={20} />
              <span>Уведомления</span>
            </Link>
          </div>
        </div>
      </div>

      {/* Overlay for Mobile Menu */}
      {isMobileMenuOpen && (
        <div
          className="mobile-overlay"
          onClick={toggleMobileMenu}
        ></div>
      )}
    </>
  );
}

export default Navbar;
