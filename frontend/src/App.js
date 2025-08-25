import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import ProductCatalog from './pages/ProductCatalog';
import About from './pages/About'; // Changed from Stores
import StoreDetail from './pages/StoreDetail';
import UserProfile from './pages/UserProfile';
import Notifications from './pages/Notifications';
import ProductDetail from './pages/ProductDetail';
import CartPage from './pages/CartPage'; // Import CartPage

function App() {
  return (
    <Router>
      <Navbar />
      <div className="container" style={{ padding: '20px' }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/products" element={<ProductCatalog />} />
          <Route path="/products/:productID" element={<ProductDetail />} />
          <Route path="/about" element={<About />} /> {/* New route for About page (stores list) */}
          <Route path="/stores/:storeID" element={<StoreDetail />} />
          <Route path="/profile" element={<Navigate to="/users/profile" replace />} />
          <Route path="/users/profile" element={<UserProfile />} />
          <Route path="/notifications" element={<Notifications />} />
          <Route path="/cart" element={<CartPage />} /> {/* New route for CartPage */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
