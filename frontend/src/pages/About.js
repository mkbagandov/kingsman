import React, { useState, useEffect } from 'react';
import { getStores } from '../api/api';
import StoreCard from '../components/StoreCard';
import './About.css'; // Assuming you'll create this CSS file

function About() {
  const [stores, setStores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStores = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await getStores();
        setStores(response.data.stores);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchStores();
  }, []);

  if (loading) return <div className="loading-message">Загрузка магазинов...</div>;
  if (error) return <div className="error-message">Ошибка: {error}</div>;

  return (
    <div className="about-page-container">
      <h1 className="page-title">Наши Магазины</h1>
      <div className="stores-grid">
        {stores.length > 0 ? (
          stores.map((store) => (
            <StoreCard key={store.id} store={store} />
          ))
        ) : (
          <p className="no-stores-message">Магазины не найдены.</p>
        )}
      </div>
    </div>
  );
}

export default About;
