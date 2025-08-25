import { configureStore } from '@reduxjs/toolkit';
import cartReducer from './cartSlice';
import authReducer from './authSlice'; // Import the auth reducer

export const store = configureStore({
  reducer: {
    cart: cartReducer,
    auth: authReducer, // Add the auth reducer
    // Add other reducers here if you have any
  },
});
