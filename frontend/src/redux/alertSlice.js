import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  alerts: [],
};

export const alertSlice = createSlice({
  name: 'alert',
  initialState,
  reducers: {
    addAlert: (state, action) => {
      // action.payload should be { id: string, message: string, type: 'success' | 'error' | 'info' }
      state.alerts.push(action.payload);
    },
    removeAlert: (state, action) => {
      state.alerts = state.alerts.filter((alert) => alert.id !== action.payload);
    },
  },
});

export const { addAlert, removeAlert } = alertSlice.actions;

export default alertSlice.reducer;
