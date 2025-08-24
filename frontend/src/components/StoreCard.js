import React from 'react';
import { Link } from 'react-router-dom';
import { FaMapMarkerAlt, FaPhone, FaStore, FaImage } from 'react-icons/fa'; // Import FaImage
import './StoreCard.css'; 

function StoreCard({ store }) {
  const imageUrl = store.ImageURL || "placeholder"; 

  return (
    <div className="store-card">
      <Link to={`/stores/${store.id}`} className="store-link">
        <div className="store-image-container">
          {imageUrl !== "placeholder" ? (
            <img src={imageUrl} alt={store.name} className="store-image" />
          ) : (
            <FaImage className="store-placeholder-icon" />
          )}
        </div>
        <div className="store-info">
          <h2 className="store-card-name">
            <FaStore className="store-icon" />
            {store.name}
          </h2>
          <p className="store-card-address">
            <FaMapMarkerAlt className="icon" /> Адрес: {store.address}
          </p>
          <p className="store-card-phone">
            <FaPhone className="icon" /> Телефон: {store.phone}
          </p>
          {/* Add more store details as needed */}
        </div>
      </Link>
    </div>
  );
}

export default StoreCard;
