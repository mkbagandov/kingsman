import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useDispatch } from 'react-redux'; // Import useDispatch
import { getProductByID, getCategories } from '../api/api';
import { FaTag, FaDollarSign, FaBoxes, FaShoppingCart } from 'react-icons/fa'; // Import FaShoppingCart icon
import { addToCart } from '../redux/cartSlice'; // Import addToCart action
import '../components/ProductCard.css'; // Reusing some styles from ProductCard

function ProductDetail() {
  const { productID } = useParams();
  const dispatch = useDispatch(); // Initialize useDispatch
  const [product, setProduct] = useState(null);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [quantity, setQuantity] = useState(1); // State for quantity

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [productRes, categoriesRes] = await Promise.all([
          getProductByID(productID),
          getCategories()
        ]);
        setProduct(productRes.data.product);
        setCategories(categoriesRes.data.categories);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [productID]);

  const handleAddToCart = () => {
    if (product) {
      dispatch(addToCart({ productID: product.id.toString(), quantity }));
    }
  };

  if (loading) return <div>Загрузка информации о продукте...</div>;
  if (error) return <div>Ошибка: {error}</div>;
  if (!product) return <div>Продукт не найден.</div>;

  const category = categories.find(cat => cat.id === product.category_id);
  const categoryName = category ? category.name : 'Неизвестно';

  return (
    <div className="product-detail-page">
      <h1>{product.name}</h1>
      <div className="product-detail-content">
        {product.ImageURL && <img src={product.ImageURL} alt={product.name} className="product-detail-image" />}
        <div className="product-detail-info">
          <p className="product-detail-description">{product.description}</p>
          <p className="product-detail-price"><FaDollarSign className="icon" /> Цена: {product.price.toFixed(2)}</p>
          <p className="product-detail-category"><FaTag className="icon" /> Категория: {categoryName}</p>
          <p className="product-detail-quantity"><FaBoxes className="icon" /> Количество: {product.quantity}</p>
          
          <div className="add-to-cart-controls">
            <input
              type="number"
              min="1"
              value={quantity}
              onChange={(e) => setQuantity(parseInt(e.target.value))}
              className="quantity-input"
            />
            <button className="add-to-cart-button" onClick={handleAddToCart}>
              <FaShoppingCart /> Добавить в корзину
            </button>
          </div>
          {/* Add more product details as needed */}
        </div>
      </div>
    </div>
  );
}

export default ProductDetail;
