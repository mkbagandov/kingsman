import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchCart, updateCart, removeFromCart, clearUserCart } from '../redux/cartSlice';
import './Cart.css';

const Cart = () => {
  const dispatch = useDispatch();
  const { cartItems, status, error } = useSelector((state) => state.cart);

  useEffect(() => {
    dispatch(fetchCart());
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

  if (status === 'loading') {
    return <div className="cart-container">Загрузка корзины...</div>;
  }

  if (status === 'failed') {
    return <div className="cart-container">Ошибка: {error}</div>;
  }

  if (!cartItems || cartItems.length === 0) {
    return <div className="cart-container">Ваша корзина пуста.</div>;
  }

  const totalAmount = cartItems.reduce((total, item) => {
    const itemPrice = item.product ? item.product.price : 0;
    return total + (itemPrice * item.quantity);
  }, 0);

  const placeholderImage = 'https://fotobudka24.ru/wp-content/uploads/2015/03/icon-image1.png'; // Placeholder image URL with a gray background and white text, translated

  return (
    <div className="cart-container">
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
      <button onClick={handleClearCart}>Очистить корзину</button>
      <div className="cart-summary">
        <h3>Итого: ${totalAmount.toFixed(2)}</h3> 
      </div>
    </div>
  );
};

export default Cart;
