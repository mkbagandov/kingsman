import React from 'react';
import { Link } from 'react-router-dom';
import { FaMapMarkerAlt, FaPhone, FaStore } from 'react-icons/fa';
import './StoreCard.css'; // Assuming you'll create this CSS file

function StoreCard({ store }) {
  return (
    <div className="store-card">
      <h2 className="store-card-name">
        <FaStore className="store-icon" />
        <Link to={`/stores/${store.id}`}>{store.name}</Link>
      </h2>
      <p className="store-card-address">
        <FaMapMarkerAlt className="icon" /> Адрес: {store.address}
      </p>
      <p className="store-card-phone">
        <FaPhone className="icon" /> Телефон: {store.phone}
      </p>
      {/* Add more store details as needed */}
    </div>
  );
}

export default StoreCard;
