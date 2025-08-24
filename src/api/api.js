import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080'; // Assuming backend runs on 8080

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// User Endpoints
export const registerUser = (userData) => api.post('/register', userData);
export const loginUser = (credentials) => api.post('/login', credentials);
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
