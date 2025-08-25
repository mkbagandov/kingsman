import React from 'react';
import { Link } from 'react-router-dom';
import { useDispatch } from 'react-redux'; // Import useDispatch
import { FaTag, FaDollarSign, FaBoxes, FaImage, FaHeart, FaShoppingCart } from 'react-icons/fa'; // Import FaShoppingCart icon
import { addToCart } from '../redux/cartSlice'; // Import addToCart action
import './ProductCard.css';

function ProductCard({ product, categories }) {
  const dispatch = useDispatch(); // Initialize useDispatch
  const category = categories.find(cat => cat.id === product.category_id);
  const categoryName = category ? category.name : 'Неизвестно';

  const imageUrl = product.ImageURL || 'https://via.placeholder.com/200/cccccc/ffffff?text=Нет+Изображения'; // Placeholder image URL for product card
  const discount = product.discount || null; // Assuming product might have a discount field

  const handleAddToCart = (e) => {
    e.preventDefault(); // Prevent navigating to product detail page
    dispatch(addToCart({ productID: product.id.toString(), quantity: 1 })); // Add 1 quantity by default
  };

  return (
    <div className="product-card">
      <Link to={`/products/${product.id}`} className="product-link">
        <div className="product-image-container">
          {imageUrl !== "https://via.placeholder.com/200/cccccc/ffffff?text=Нет+Изображения" ? (
            <img src={imageUrl} alt={product.name} className="product-image" />
          ) : (
            <FaImage className="product-placeholder-icon" />
          )}
          {discount && <div className="discount-badge">-{discount}%</div>}
          <div className="favorite-icon"><FaHeart /></div>
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
      <button className="add-to-cart-button" onClick={handleAddToCart}>
        <FaShoppingCart /> Добавить в корзину
      </button>
    </div>
  );
}

export default ProductCard;
