import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getStoreByID } from '../api/api';

function StoreDetail() {
  const { storeID } = useParams();
  const [store, setStore] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStore = async () => {
      try {
        const response = await getStoreByID(storeID);
        setStore(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchStore();
  }, [storeID]);

  if (loading) return <div>Загрузка информации о магазине...</div>;
  if (error) return <div>Ошибка: {error}</div>;
  if (!store) return <div>Магазин не найден.</div>;

  return (
    <div>
      <h1>{store.name}</h1>
      <p>Адрес: {store.address}</p>
      <p>Телефон: {store.phone}</p>
      {/* Display other store details */}
    </div>
  );
}

export default StoreDetail;
