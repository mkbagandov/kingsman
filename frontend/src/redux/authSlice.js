import { createSlice } from '@reduxjs/toolkit';

const getInitialAuthState = () => {
  const token = localStorage.getItem('jwtToken');
  return {
    isAuthenticated: !!token,
    // You can add more auth-related state here, e.g., user info
  };
};

export const authSlice = createSlice({
  name: 'auth',
  initialState: getInitialAuthState(),
  reducers: {
    login: (state) => {
      state.isAuthenticated = true;
    },
    logout: (state) => {
      state.isAuthenticated = false;
    },
    // Add more reducers as needed for updating user info, etc.
  },
});

export const { login, logout } = authSlice.actions;

export default authSlice.reducer;
