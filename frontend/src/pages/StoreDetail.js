import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getStoreByID } from '../api/api';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'; // Import Leaflet components
import { FaMapMarkerAlt, FaPhone, FaClock, FaInfoCircle } from 'react-icons/fa';
import 'leaflet/dist/leaflet.css'; // Import Leaflet CSS
import L from 'leaflet'; // Import Leaflet for custom icon
import './StoreDetail.css'; 

// Fix for default marker icon not showing
delete L.Icon.Default.prototype._getIconUrl;

L.Icon.Default.mergeOptions({
  iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: require('leaflet/dist/images/marker-icon.png'),
  shadowUrl: require('leaflet/dist/images/marker-shadow.png'),
});

function StoreDetail() {
  const { storeID } = useParams();
  const [store, setStore] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [position, setPosition] = useState(null); // State for map position [latitude, longitude]

  useEffect(() => {
    const fetchStore = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await getStoreByID(storeID);
        setStore(response.data.store);
        // Assuming your store data includes latitude and longitude
        if (response.data.store && response.data.store.latitude && response.data.store.longitude) {
          setPosition([response.data.store.latitude, response.data.store.longitude]);
        }
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };

    if (storeID) {
      fetchStore();
    }
  }, [storeID]);

  if (loading) return <div className="loading-message">Загрузка информации о магазине...</div>;
  if (error) return <div className="error-message">Ошибка: {error}</div>;
  if (!store) return <div className="no-store-message">Магазин не найден.</div>;

  return (
    <div className="store-detail-page">
      <h1 className="page-title">{store.name}</h1>
      <div className="store-detail-content">
        {store.ImageURL && <img src={store.ImageURL} alt={store.name} className="store-detail-image" />}
        <div className="store-detail-info">
          <p className="detail-item"><FaMapMarkerAlt className="icon" /> Адрес: {store.address}</p>
          <p className="detail-item"><FaPhone className="icon" /> Телефон: {store.phone}</p>
          {store.opening_hours && <p className="detail-item"><FaClock className="icon" /> Часы работы: {store.opening_hours}</p>}
          {store.description && <p className="detail-item"><FaInfoCircle className="icon" /> Описание: {store.description}</p>}
          {/* Add map if position is available */}
          {position && (
            <div className="map-container">
              <h2>Местоположение</h2>
              <MapContainer center={position} zoom={13} scrollWheelZoom={false} className="leaflet-map">
                <TileLayer
                  attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                <Marker position={position}>
                  <Popup>
                    {store.name} <br /> {store.address}
                  </Popup>
                </Marker>
              </MapContainer>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default StoreDetail;
