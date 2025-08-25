import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getCart, addItemToCart, updateCartItem, removeCartItem, clearCart, getProductByID, placeOrder } from '../api/api'; // Import placeOrder

export const fetchCart = createAsyncThunk(
  'cart/fetchCart',
  async (_, { rejectWithValue }) => {
    try {
      const response = await getCart();
      const cart = response.data.cart;
      let cartItems = response.data.cart_items || [];

      // Fetch product details for each item in the cart
      if (cartItems.length > 0) {
        const productDetailsPromises = cartItems.map(async (item) => {
          const productRes = await getProductByID(item.product_id);
          return { ...item, product: productRes.data.product }; // Combine cart item with product details
        });
        cartItems = await Promise.all(productDetailsPromises);
      }
      
      return { cart, cartItems };
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const addToCart = createAsyncThunk(
  'cart/addToCart',
  async ({ productID, quantity }, { rejectWithValue }) => {
    try {
      const response = await addItemToCart(productID, quantity);
      const newItem = response.data;
      const productRes = await getProductByID(productID);
      return { ...newItem, product: productRes.data.product }; // Return item with product details
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const updateCart = createAsyncThunk(
  'cart/updateCart',
  async ({ productID, quantity }, { rejectWithValue }) => {
    try {
      const response = await updateCartItem(productID, quantity);
      const updatedItem = response.data;
      const productRes = await getProductByID(productID);
      return { ...updatedItem, product: productRes.data.product }; // Return item with product details
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const removeFromCart = createAsyncThunk(
  'cart/removeFromCart',
  async (productID, { rejectWithValue }) => {
    try {
      await removeCartItem(productID);
      return productID;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const clearUserCart = createAsyncThunk(
  'cart/clearUserCart',
  async (_, { rejectWithValue }) => {
    try {
      await clearCart();
      return null;
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

export const processPayment = createAsyncThunk(
  'cart/processPayment',
  async (_, { rejectWithValue }) => {
    try {
      const response = await placeOrder();
      return response.data; // Should contain order ID and a message
    } catch (error) {
      return rejectWithValue(error.response.data);
    }
  }
);

const cartSlice = createSlice({
  name: 'cart',
  initialState: {
    cart: null,
    cartItems: [],
    status: 'idle',
    error: null,
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchCart.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(fetchCart.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.cart = action.payload.cart;
        state.cartItems = action.payload.cartItems; // Now contains product details
      })
      .addCase(fetchCart.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload;
      })
      .addCase(addToCart.fulfilled, (state, action) => {
        state.status = 'succeeded';
        const newItemWithProduct = action.payload;
        const existingItemIndex = state.cartItems.findIndex(item => item.product_id === newItemWithProduct.product_id);
        if (existingItemIndex !== -1) {
          state.cartItems[existingItemIndex].quantity = newItemWithProduct.quantity;
          state.cartItems[existingItemIndex].product = newItemWithProduct.product; // Update product details too if needed
        } else {
          state.cartItems.push(newItemWithProduct);
        }
      })
      .addCase(updateCart.fulfilled, (state, action) => {
        state.status = 'succeeded';
        const updatedItemWithProduct = action.payload;
        const existingItemIndex = state.cartItems.findIndex(item => item.product_id === updatedItemWithProduct.product_id);
        if (existingItemIndex !== -1) {
          state.cartItems[existingItemIndex].quantity = updatedItemWithProduct.quantity;
          state.cartItems[existingItemIndex].product = updatedItemWithProduct.product; // Update product details too if needed
        }
      })
      .addCase(removeFromCart.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.cartItems = state.cartItems.filter(item => item.product_id !== action.payload);
      })
      .addCase(clearUserCart.fulfilled, (state) => {
        state.status = 'succeeded';
        state.cart = null;
        state.cartItems = [];
      })
      .addCase(processPayment.pending, (state) => {
        state.status = 'loading';
      })
      .addCase(processPayment.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.cart = null; // Clear the cart after successful payment
        state.cartItems = [];
        // Optionally, you might want to store the order ID or message from action.payload
      })
      .addCase(processPayment.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload;
      });
  },
});

export default cartSlice.reducer;
