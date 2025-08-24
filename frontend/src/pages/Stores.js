import React, { useState, useEffect } from 'react';
import { getStores } from '../api/api';
import StoreCard from '../components/StoreCard'; // Import the new StoreCard component

function Stores() {
  const [stores, setStores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStores = async () => {
      try {
        const response = await getStores();
        // The response JSON structure is { "stores": [...] }, so access response.data.stores
        setStores(response.data.stores);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchStores();
  }, []);

  if (loading) return <div>Загрузка магазинов...</div>;
  if (error) return <div>Ошибка: {error}</div>;

  return (
    <div className="stores-page">
      <h1>Наши Магазины</h1>
      <div className="stores-grid">
        {stores.map((store) => (
          <StoreCard key={store.id} store={store} />
        ))}
      </div>
    </div>
  );
}

export default Stores;
