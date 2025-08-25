import React, { useEffect, useState } from 'react'; // Import useState
import { useSelector, useDispatch } from 'react-redux';
import { fetchCart, updateCart, removeFromCart, clearUserCart, processPayment } from '../redux/cartSlice'; // Import processPayment
import { Link } from 'react-router-dom'; // Import Link
import { addAlert } from '../redux/alertSlice'; // Import addAlert action
import { v4 as uuidv4 } from 'uuid'; // Import uuid for unique alert IDs
import './Cart.css';
import { FaShoppingCart } from 'react-icons/fa'; // Import FaShoppingCart for the empty cart icon

const Cart = () => {
  const dispatch = useDispatch();
  const { cartItems, status, error } = useSelector((state) => state.cart);
  const isAuthenticated = useSelector((state) => state.auth.isAuthenticated); // Get authentication status
  const [paymentInfo, setPaymentInfo] = useState({
    cardNumber: '',
    expiryDate: '',
    cvv: '',
  });

  useEffect(() => {
    // Dispatch fetchCart and handle potential errors for displaying alerts
    const fetchCartData = async () => {
      try {
        const result = await dispatch(fetchCart()).unwrap();
        if (!result || !result.cartItems || result.cartItems.length === 0) {
          dispatch(addAlert({ id: uuidv4(), message: 'Пусто: Корзина пуста.', type: 'success' })); // Changed to success
        }
      } catch (err) {
        const errorMessage = err.message || 'Неизвестная ошибка при загрузке корзины.';
        dispatch(addAlert({ id: uuidv4(), message: `Ошибка: ${errorMessage}`, type: 'error' }));
      }
    };
    fetchCartData();
  }, [dispatch]);

  const handleUpdateQuantity = (productID, quantity) => {
    dispatch(updateCart({ productID, quantity }));
  };

  const handleRemoveItem = (productID) => {
    dispatch(removeFromCart(productID));
  };

  const handleClearCart = () => {
    dispatch(clearUserCart());
  };

  const handlePaymentInfoChange = (e) => {
    const { name, value } = e.target;
    setPaymentInfo((prevInfo) => ({
      ...prevInfo,
      [name]: value,
    }));
  };

  const handleProcessPayment = async () => {
    // In a real application, you would send paymentInfo to a payment gateway.
    // For this test implementation, we just dispatch the processPayment thunk.
    try {
      await dispatch(processPayment()).unwrap();
      dispatch(addAlert({ id: uuidv4(), message: 'Заказ успешно оплачен!', type: 'success' }));
    } catch (err) {
      const errorMessage = err.message || 'Неизвестная ошибка при оплате.';
      dispatch(addAlert({ id: uuidv4(), message: `Ошибка при оплате: ${errorMessage}`, type: 'error' }));
    }
  };

  if (status === 'loading') {
    return <div className="cart-container">Загрузка корзины...</div>;
  }

  if (status === 'failed') {
    return <div className="cart-container">Ошибка: {error}</div>;
  }

  const totalAmount = cartItems.reduce((total, item) => {
    const itemPrice = item.product ? item.product.price : 0;
    return total + (itemPrice * item.quantity);
  }, 0);

  const placeholderImage = 'https://fotobudka24.ru/wp-content/uploads/2015/03/icon-image1.png'; // Placeholder image URL with a gray background and white text, translated

  return (
    <div className="cart-container">
 {/* Added back the cart title */}
      {cartItems.length > 0 ? (
        <>
          <div className="cart-items">
            {cartItems.map((item) => (
              <div key={item.product_id} className="cart-item">
                <img 
                  src={placeholderImage} 
                  alt={item.product?.name || 'Изображение товара'}
                  className="cart-item-image" 
                />
                <div className="cart-item-details">
                  <p className="cart-item-name">{item.product?.name || 'Неизвестный товар'}</p>
                  <p className="cart-item-price">${(item.product?.price * item.quantity).toFixed(2)}</p>
                </div>
                <input
                  type="number"
                  min="1"
                  value={item.quantity}
                  onChange={(e) => handleUpdateQuantity(item.product_id, parseInt(e.target.value))}
                />
                <button onClick={() => handleRemoveItem(item.product_id)}>Удалить</button>
              </div>
            ))}
          </div>
          <div className="cart-summary">
            <h3>Итого: ₽{totalAmount.toFixed(2)}</h3> 
            <button onClick={handleClearCart} className="clear-cart-button">Очистить корзину</button>
          </div>
        </>
      ) : (
        <div className="empty-cart-message-container">
          <FaShoppingCart className="empty-cart-icon" />
          <p>Ваша корзина пуста.</p>
        </div>
      )}

      {isAuthenticated && (
        <Link to="/orders" className="view-orders-button-bottom">История заказов</Link>
      )}

      {cartItems.length > 0 && (
        <div className="payment-card">
          <h3>Оплата заказа</h3>
          <div className="form-group">
            <label htmlFor="cardNumber">Номер карты</label>
            <input
              type="text"
              id="cardNumber"
              name="cardNumber"
              value={paymentInfo.cardNumber}
              onChange={handlePaymentInfoChange}
              placeholder="XXXX XXXX XXXX XXXX"
              maxLength="19" // Max length for credit card numbers (16 digits + 3 spaces)
            />
          </div>
          <div className="form-group-row">
            <div className="form-group">
              <label htmlFor="expiryDate">Срок действия</label>
              <input
                type="text"
                id="expiryDate"
                name="expiryDate"
                value={paymentInfo.expiryDate}
                onChange={handlePaymentInfoChange}
                placeholder="ММ/ГГ"
                maxLength="5" // MM/YY format
              />
            </div>
            <div className="form-group">
              <label htmlFor="cvv">CVV</label>
              <input
                type="text"
                id="cvv"
                name="cvv"
                value={paymentInfo.cvv}
                onChange={handlePaymentInfoChange}
                placeholder="123"
                maxLength="3" // Standard CVV length
              />
            </div>
          </div>
          <button onClick={handleProcessPayment} className="pay-button">Оплатить</button>
        </div>
      )}
    </div>
  );
};

export default Cart;
