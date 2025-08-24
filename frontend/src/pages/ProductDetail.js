import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getProductByID, getCategories } from '../api/api';
import { FaTag, FaDollarSign, FaBoxes } from 'react-icons/fa';
import '../components/ProductCard.css'; // Reusing some styles from ProductCard

function ProductDetail() {
  const { productID } = useParams();
  const [product, setProduct] = useState(null);
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

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
          {/* Add more product details as needed */}
        </div>
      </div>
    </div>
  );
}

export default ProductDetail;
