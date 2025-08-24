import os

def create_file(filepath, content):
    os.makedirs(os.path.dirname(filepath), exist_ok=True)
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)

# Define file contents
api_js_content = '''import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080'; // Assuming backend runs on 8080

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// User Endpoints
export const registerUser = (userData) => api.post('/users/register', userData);
export const loginUser = (credentials) => api.post('/users/login', credentials);
export const getUserProfile = (userID) => api.get(`/users/${userID}`);
export const getUserLoyaltyProfile = (userID) => api.get(`/users/${userID}/loyalty`);
export const addLoyaltyPoints = (userID, data) => api.post(`/users/${userID}/loyalty/points`, data);
export const addLoyaltyActivity = (userID, data) => api.post(`/users/${userID}/loyalty/activity`, data);
export const getLoyaltyTiers = () => api.get('/loyalty/tiers');
export const getUserDiscountCard = (userID) => api.get(`/users/${userID}/discount-card`);
export const updateUserDiscountCard = (userID, data) => api.put(`/users/${userID}/discount-card`, data);
export const getUserQRCode = (userID) => api.get(`/users/${userID}/qrcode`, { responseType: 'arraybuffer' }); // For image data

// Store Endpoints
export const getStores = () => api.get('/stores');
export const getStoreByID = (storeID) => api.get(`/stores/${storeID}`);

// Product Endpoints
export const getProductCatalog = (params) => api.get('/products', { params });
export const getProductByID = (productID) => api.get(`/products/${productID}`);

// Category Endpoints
export const getCategories = () => api.get('/categories');

// Notification Endpoints
export const sendNotification = (notificationData) => api.post('/notifications', notificationData);
export const getNotifications = (userID) => api.get(`/notifications/${userID}`);

export default api;
'''

home_js_content = '''import React from 'react';

function Home() {
  return (
    <div>
      <h1>Welcome to Kingsman App</h1>
      <p>This is the homepage.</p>
    </div>
  );
}

export default Home;
'''

login_js_content = '''import React, { useState } from 'react';
import { loginUser } from '../api/api';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await loginUser({ email, password });
      setMessage(`Login successful: ${response.data.message}`);
      // Handle successful login, e.g., store token, redirect
    } catch (error) {
      setMessage(`Login failed: ${error.response?.data?.error || error.message}`);
    }
  };

  return (
    <div>
      <h1>Login</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        </div>
        <div>
          <label>Password:</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        </div>
        <button type="submit">Login</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default Login;
'''

register_js_content = '''import React, { useState } from 'react';
import { registerUser } from '../api/api';

function Register() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await registerUser({ username, email, password });
      setMessage(`Registration successful: ${response.data.message}`);
      // Optionally redirect to login or show success message
    } catch (error) {
      setMessage(`Registration failed: ${error.response?.data?.error || error.message}`);
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>Username:</label>
          <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} required />
        </div>
        <div>
          <label>Email:</label>
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        </div>
        <div>
          <label>Password:</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        </div>
        <button type="submit">Register</button>
      </form>
      {message && <p>{message}</p>}
    </div>
  );
}

export default Register;
'''

productcatalog_js_content = '''import React, { useState, useEffect } from 'react';
import { getProductCatalog, getCategories } from '../api/api';

function ProductCatalog() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchProducts();
    fetchCategories();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const params = {
        ...(selectedCategory && { category_id: selectedCategory }),
        ...(minPrice && { min_price: parseFloat(minPrice) }),
        ...(maxPrice && { max_price: parseFloat(maxPrice) }),
      };
      const response = await getProductCatalog(params);
      setProducts(response.data.products);
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await getCategories();
      setCategories(response.data);
    } catch (err) {
      console.error("Error fetching categories:", err);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchProducts();
  };

  if (loading) return <div>Loading products...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Product Catalog</h1>
      <form onSubmit={handleSearch}>
        <div>
          <label>Category:</label>
          <select value={selectedCategory} onChange={(e) => setSelectedCategory(e.target.value)}>
            <option value="">All Categories</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label>Min Price:</label>
          <input type="number" value={minPrice} onChange={(e) => setMinPrice(e.target.value)} />
        </div>
        <div>
          <label>Max Price:</label>
          <input type="number" value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)} />
        </div>
        <button type="submit">Filter Products</button>
      </form>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
        {products.map((product) => (
          <div key={product.id} style={{ border: '1px solid #ccc', padding: '15px', borderRadius: '8px' }}>
            <h2>{product.name}</h2>
            <p>{product.description}</p>
            <p>Price: ${product.price}</p>
            <p>Category ID: {product.category_id}</p>
            {/* Add more product details as needed */}
          </div>
        ))}
      </div>
    </div>
  );
}

export default ProductCatalog;
'''

stores_js_content = '''import React, { useState, useEffect } from 'react';
import { getStores } from '../api/api';
import { Link } from 'react-router-dom';

function Stores() {
  const [stores, setStores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStores = async () => {
      try {
        const response = await getStores();
        setStores(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchStores();
  }, []);

  if (loading) return <div>Loading stores...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Our Stores</h1>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
        {stores.map((store) => (
          <div key={store.id} style={{ border: '1px solid #ccc', padding: '15px', borderRadius: '8px' }}>
            <h2><Link to={`/stores/${store.id}`}>{store.name}</Link></h2>
            <p>{store.address}</p>
            <p>Phone: {store.phone}</p>
            {/* Add more store details as needed */}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Stores;
'''

storedetail_js_content = '''import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getStoreByID } from '../api/api';

function StoreDetail() {
  const { storeID } = useParams();
  const [store, setStore] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStore = async () => {
      try {
        const response = await getStoreByID(storeID);
        setStore(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchStore();
  }, [storeID]);

  if (loading) return <div>Loading store details...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!store) return <div>Store not found.</div>;

  return (
    <div>
      <h1>{store.name}</h1>
      <p>Address: {store.address}</p>
      <p>Phone: {store.phone}</p>
      {/* Display other store details */}
    </div>
  );
}

export default StoreDetail;
'''

userprofile_js_content = '''import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getUserProfile, getUserLoyaltyProfile, getUserDiscountCard, getUserQRCode } from '../api/api';

function UserProfile() {
  const { userID } = useParams();
  const [userProfile, setUserProfile] = useState(null);
  const [loyaltyProfile, setLoyaltyProfile] = useState(null);
  const [discountCard, setDiscountCard] = useState(null);
  const [qrCode, setQrCode] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [userRes, loyaltyRes, discountRes, qrCodeRes] = await Promise.all([
          getUserProfile(userID),
          getUserLoyaltyProfile(userID),
          getUserDiscountCard(userID),
          getUserQRCode(userID)
        ]);
        setUserProfile(userRes.data);
        setLoyaltyProfile(loyaltyRes.data);
        setDiscountCard(discountRes.data);

        // Convert arraybuffer to base64 for image display
        const base64Image = btoa(
          new Uint8Array(qrCodeRes.data).reduce(
            (data, byte) => data + String.fromCharCode(byte),
            ''
          )
        );
        setQrCode(`data:image/png;base64,${base64Image}`);

      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [userID]);

  if (loading) return <div>Loading user profile...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!userProfile) return <div>User profile not found.</div>;

  return (
    <div>
      <h1>User Profile: {userProfile.username}</h1>
      <p>Email: {userProfile.email}</p>

      <h2>Loyalty Information</h2>
      {loyaltyProfile ? (
        <div>
          <p>Points: {loyaltyProfile.points}</p>
          <p>Tier: {loyaltyProfile.tier}</p>
          {/* Display other loyalty details */}
        </div>
      ) : (
        <p>No loyalty information available.</p>
      )}

      <h2>Discount Card</h2>
      {discountCard ? (
        <div>
          <p>Discount Level: {discountCard.discount_level}</p>
          <p>Progress to Next Level: {discountCard.progress_to_next_level}%</p>
          {/* Display other discount card details */}
        </div>
      ) : (
        <p>No discount card information available.</p>
      )}

      <h2>QR Code</h2>
      {qrCode && <img src={qrCode} alt="QR Code" />}
    </div>
  );
}

export default UserProfile;
'''

notifications_js_content = '''import React, { useState, useEffect } from 'react';
import { getNotifications } from '../api/api';

function Notifications() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  // In a real app, you would get the userID from authentication context
  const userID = "1"; // Placeholder userID

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await getNotifications(userID);
        setNotifications(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchNotifications();
  }, [userID]);

  if (loading) return <div>Loading notifications...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Your Notifications</h1>
      {notifications.length > 0 ? (
        <ul>
          {notifications.map((notification) => (
            <li key={notification.id}>
              <strong>{notification.type}:</strong> {notification.message} - {new Date(notification.timestamp).toLocaleString()}
            </li>
          ))}
        </ul>
      ) : (
        <p>No notifications found.</p>
      )}
    </div>
  );
}

export default Notifications;
'''

navbar_js_content = '''import React from 'react';
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
'''

app_js_content = '''import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import ProductCatalog from './pages/ProductCatalog';
import Stores from './pages/Stores';
import StoreDetail from './pages/StoreDetail';
import UserProfile from './pages/UserProfile';
import Notifications from './pages/Notifications';

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
          <Route path="/stores" element={<Stores />} />
          <Route path="/stores/:storeID" element={<StoreDetail />} />
          <Route path="/profile/:userID" element={<UserProfile />} />
          <Route path="/notifications/:userID" element={<Notifications />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
'''

index_js_content = '''import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
'''

index_css_content = '''body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New',
    monospace;
}
'''

# Create files
create_file('src/api/api.js', api_js_content)
create_file('src/pages/Home.js', home_js_content)
create_file('src/pages/Login.js', login_js_content)
create_file('src/pages/Register.js', register_js_content)
create_file('src/pages/ProductCatalog.js', productcatalog_js_content)
create_file('src/pages/Stores.js', stores_js_content)
create_file('src/pages/StoreDetail.js', storedetail_js_content)
create_file('src/pages/UserProfile.js', userprofile_js_content)
create_file('src/pages/Notifications.js', notifications_js_content)
create_file('src/components/Navbar.js', navbar_js_content)
create_file('src/App.js', app_js_content)
create_file('src/index.js', index_js_content)
create_file('src/index.css', index_css_content)

print("Frontend files and directories created successfully!")
