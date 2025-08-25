import React from 'react';
import { Link } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { FaHeart, FaShoppingCart, FaImage } from 'react-icons/fa'; // Added FaImage icon
import { addToCart } from '../redux/cartSlice';
import './ProductCard.css';

function ProductCard({ product, categories }) {
  const dispatch = useDispatch();
  // No longer need category name explicitly displayed on the card, but it's useful if we want to filter by category later.
  // const category = categories.find(cat => cat.id === product.category_id);
  // const categoryName = category ? category.name : 'Неизвестно';

  const imageUrl = product.ImageURL; // Use actual image URL directly
  const discount = product.discount || null;

  const handleAddToCart = (e) => {
    e.preventDefault();
    dispatch(addToCart({ productID: product.id.toString(), quantity: 1 }));
  };

  return (
    <div className="product-card">
      <Link to={`/products/${product.id}`} className="product-link">
        <div className="product-image-wrapper">
          {imageUrl ? (
            <img src={imageUrl} alt={product.name} className="product-image" />
          ) : (
            <div className="product-placeholder-image">
              <FaImage className="product-placeholder-icon" />
              <span>Нет Изображения</span>
            </div>
          )}
          {discount && <div className="discount-badge">-{discount}%</div>}
          <div className="favorite-icon"><FaHeart /></div>
        </div>
        <div className="product-content">
          <h3 className="product-name">{product.name}</h3>
          {/* <p className="product-description">{product.description}</p> */}
          <div className="product-price-wrapper">
            <span className="product-price">${product.price.toFixed(2)}</span>
          </div>
        </div>
      </Link>
      <button className="add-to-cart-button" onClick={handleAddToCart}>
        <FaShoppingCart />
      </button>
    </div>
  );
}

export default ProductCard;
