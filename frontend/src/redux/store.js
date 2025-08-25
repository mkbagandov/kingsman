import { configureStore } from '@reduxjs/toolkit';
import cartReducer from './cartSlice';
import authReducer from './authSlice'; // Import the auth reducer
import alertReducer from './alertSlice'; // Import the alert reducer

export const store = configureStore({
  reducer: {
    cart: cartReducer,
    auth: authReducer, // Add the auth reducer
    alert: alertReducer, // Add the alert reducer
    // Add other reducers here if you have any
  },
});
