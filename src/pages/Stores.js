import React, { useState, useEffect } from 'react';
import { getStores } from '../api/api';
import { Link } from 'react-router-dom';

function Stores() {
  const [stores, setStores] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStores = async () => {
      try {
        const response = await getStores();
        setStores(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchStores();
  }, []);

  if (loading) return <div>Loading stores...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Our Stores</h1>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
        {stores.map((store) => (
          <div key={store.id} style={{ border: '1px solid #ccc', padding: '15px', borderRadius: '8px' }}>
            <h2><Link to={`/stores/${store.id}`}>{store.name}</Link></h2>
            <p>{store.address}</p>
            <p>Phone: {store.phone}</p>
            {/* Add more store details as needed */}
          </div>
        ))}
      </div>
    </div>
  );
}

export default Stores;
