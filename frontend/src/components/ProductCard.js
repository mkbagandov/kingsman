import React from 'react';
import { Link } from 'react-router-dom';
import { FaTag, FaDollarSign, FaBoxes, FaImage } from 'react-icons/fa'; // Import FaImage for placeholder
import './ProductCard.css';

function ProductCard({ product, categories }) {
  const category = categories.find(cat => cat.id === product.category_id);
  const categoryName = category ? category.name : 'Неизвестно';

  const imageUrl = product.ImageURL || "placeholder"; // Use a string 'placeholder' if no image URL

  return (
    <div className="product-card">
      <Link to={`/products/${product.id}`} className="product-link">
        <div className="product-image-container">
          {imageUrl !== "placeholder" ? (
            <img src={imageUrl} alt={product.name} className="product-image" />
          ) : (
            <FaImage className="product-placeholder-icon" />
          )}
        </div>
        <div className="product-info">
          <h2 className="product-name">{product.name}</h2>
          <p className="product-description">{product.description}</p>
          <div className="product-details">
            <p className="product-price"><FaDollarSign className="icon" /> {product.price.toFixed(2)}</p>
            <p className="product-category"><FaTag className="icon" /> {categoryName}</p>
            <p className="product-quantity"><FaBoxes className="icon" /> Количество: {product.quantity}</p>
          </div>
        </div>
      </Link>
    </div>
  );
}

export default ProductCard;
