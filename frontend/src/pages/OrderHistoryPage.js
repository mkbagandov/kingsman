import React, { useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux'; // Import useDispatch
import { getUserOrders, getProductByID } from '../api/api'; // Assuming getProductByID is available
import { addAlert } from '../redux/alertSlice'; // Import addAlert action
import { v4 as uuidv4 } from 'uuid'; // Import uuid for unique alert IDs
import './OrderHistoryPage.css'; // We'll create this CSS file

const OrderHistoryPage = () => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const userStatus = useSelector((state) => state.auth.status); // Assuming auth slice for user status
  const dispatch = useDispatch(); // Initialize useDispatch

  useEffect(() => {
    const fetchOrders = async () => {
      if (userStatus !== 'succeeded') { // Only fetch if user is logged in
        setLoading(false);
        dispatch(addAlert({ id: uuidv4(), message: 'Пусто: Пользователь не авторизован.', type: 'info' }));
        return;
      }
      try {
        setLoading(true);
        const response = await getUserOrders('paid'); // Pass 'paid' to filter orders
        let fetchedOrders = response.data.orders; // Assuming the response has an 'orders' field

        if (fetchedOrders && fetchedOrders.length > 0) {
          // Fetch product details for each item in each order
          const ordersWithProductDetails = await Promise.all(
            fetchedOrders.map(async (order) => {
              const itemsWithDetails = await Promise.all(
                order.Items.map(async (item) => {
                  try {
                    const productRes = await getProductByID(item.ProductID);
                    return { ...item, product: productRes.data.product };
                  } catch (productError) {
                    console.error(`Error fetching product ${item.ProductID}:`, productError);
                    dispatch(addAlert({ id: uuidv4(), message: `Ошибка при загрузке данных продукта для заказа: ${item.ProductID}`, type: 'error' }));
                    return { ...item, product: { name: 'Неизвестный товар', ImageURL: 'https://via.placeholder.com/50?text=N/A' } };
                  }
                })
              );
              return { ...order, Items: itemsWithDetails };
            })
          );
          setOrders(ordersWithProductDetails);
        } else {
          setOrders([]);
          dispatch(addAlert({ id: uuidv4(), message: 'Пусто: Заказы не найдены.', type: 'info' }));
        }

      } catch (err) {
        const errorMessage = err.response?.data?.error || err.message || 'Неизвестная ошибка при загрузке заказов.';
        setError(errorMessage);
        dispatch(addAlert({ id: uuidv4(), message: `Ошибка: ${errorMessage}`, type: 'error' }));
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [userStatus, dispatch]); // Added dispatch to dependencies

  if (loading) return <div className="order-history-container">Загрузка истории заказов...</div>;
  if (error) return <div className="order-history-container">Ошибка: {error}</div>;
  // if (orders.length === 0) return <div className="order-history-container">У вас пока нет заказов.</div>; // Removed as alert handles this

  return (
    <div className="order-history-container">
      <h2>История ваших заказов</h2>
      {orders.map((order) => (
        <div key={order.ID} className="order-card">
          <div className="order-header">
            <p><strong>Заказ #:</strong> {order.ID}</p>
            <p><strong>Дата:</strong> {new Date(order.OrderDate).toLocaleDateString()}</p>
            <p><strong>Статус:</strong> {order.Status}</p>
            <p><strong>Общая сумма:</strong> ${order.TotalAmount.toFixed(2)}</p>
          </div>
          <div className="order-items">
            <h4>Товары:</h4>
            {order.Items.map((item) => (
              <div key={item.ID} className="order-item">
                <img src={item.product?.ImageURL || 'https://via.placeholder.com/50?text=N/A'} alt={item.product?.name || 'Изображение товара'} className="order-item-image" />
                <div className="order-item-details">
                  <p><strong>{item.product?.name || 'Неизвестный товар'}</strong></p>
                  <p>Количество: {item.Quantity}</p>
                  <p>Цена за ед.: ${item.Price.toFixed(2)}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default OrderHistoryPage;
