import React from 'react';
import { FaTag, FaDollarSign, FaBoxes } from 'react-icons/fa';
import './ProductCard.css'; // Assuming you'll create this CSS file

function ProductCard({ product, categories }) {
  const category = categories.find(cat => cat.id === product.category_id);
  const categoryName = category ? category.name : 'Неизвестно';

  return (
    <div className="product-card">
      {product.ImageURL && <img src={product.ImageURL} alt={product.name} className="product-image" />}
      <h2 className="product-name">{product.name}</h2>
      <p className="product-description">{product.description}</p>
      <div className="product-details">
        <p className="product-price"><FaDollarSign className="icon" /> {product.price.toFixed(2)}</p>
        <p className="product-category"><FaTag className="icon" /> {categoryName}</p>
        <p className="product-quantity"><FaBoxes className="icon" /> Количество: {product.quantity}</p>
      </div>
    </div>
  );
}

export default ProductCard;
