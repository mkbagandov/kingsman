import React from 'react';
import { Link } from 'react-router-dom';

function Navbar() {
  return (
    <nav style={{ background: '#333', padding: '10px 20px', color: '#fff', display: 'flex', justifyContent: 'space-between' }}>
      <div className="nav-left">
        <Link to="/" style={{ color: '#fff', textDecoration: 'none', marginRight: '20px' }}>Home</Link>
        <Link to="/products" style={{ color: '#fff', textDecoration: 'none', marginRight: '20px' }}>Products</Link>
        <Link to="/stores" style={{ color: '#fff', textDecoration: 'none', marginRight: '20px' }}>Stores</Link>
        <Link to="/profile/1" style={{ color: '#fff', textDecoration: 'none', marginRight: '20px' }}>Profile (User 1)</Link>
        <Link to="/notifications/1" style={{ color: '#fff', textDecoration: 'none' }}>Notifications (User 1)</Link>
      </div>
      <div className="nav-right">
        <Link to="/login" style={{ color: '#fff', textDecoration: 'none', marginRight: '20px' }}>Login</Link>
        <Link to="/register" style={{ color: '#fff', textDecoration: 'none' }}>Register</Link>
      </div>
    </nav>
  );
}

export default Navbar;
