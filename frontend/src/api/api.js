import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080'; // Assuming backend runs on 8080

// Helper functions for token management
export const setAuthToken = (token) => {
  if (token) {
    localStorage.setItem('jwtToken', token);
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  } else {
    localStorage.removeItem('jwtToken');
    delete api.defaults.headers.common['Authorization'];
  }
};

export const getAuthToken = () => {
  return localStorage.getItem('jwtToken');
};

export const removeAuthToken = () => {
  setAuthToken(null);
};

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Set the initial token if it exists in localStorage
const initialToken = getAuthToken();
if (initialToken) {
  setAuthToken(initialToken);
}

// User Endpoints
export const registerUser = (userData) => api.post('/users/register', userData);
export const loginUser = async (credentials) => {
  const response = await api.post('/users/login', credentials);
  if (response.data.token) {
    setAuthToken(response.data.token);
  }
  return response;
};
export const getUserProfile = () => api.get('/users/profile');
export const getUserLoyaltyProfile = () => api.get('/users/loyalty');
export const addLoyaltyPoints = (data) => api.post('/users/loyalty-points', data);
export const addLoyaltyActivity = (data) => api.post('/users/loyalty-activity', data);
export const getLoyaltyTiers = () => api.get('/loyalty-tiers');
export const getUserDiscountCard = () => api.get('/users/discount-card');
export const updateUserDiscountCard = (data) => api.put('/users/discount-card', data);
export const getUserQRCode = () => api.get('/users/qrcode', { responseType: 'arraybuffer' }); // For image data

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
export const getNotifications = () => api.get('/users/notifications');

// Cart Endpoints
export const getCart = () => api.get('/cart');
export const addItemToCart = (productID, quantity) => api.post('/cart/items', { product_id: productID, quantity });
export const updateCartItem = (productID, quantity) => api.put('/cart/items', { product_id: productID, quantity });
export const removeCartItem = (productID) => api.delete(`/cart/items/${productID}`);
export const clearCart = () => api.delete('/cart/clear');
export const placeOrder = () => api.post('/cart/checkout'); // New: Place an order from the cart

// Order Endpoints
export const getUserOrders = (paymentStatus = null) => {
  const params = {};
  if (paymentStatus) {
    params.paymentStatus = paymentStatus;
  }
  return api.get('/orders', { params });
}; // New: Get user's order history

// User Profile Update Endpoint
export const updateUserProfile = (userData) => api.put('/users/profile', userData); // New: Update user profile

export default api;
