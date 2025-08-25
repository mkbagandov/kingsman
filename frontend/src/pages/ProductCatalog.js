import React, { useState, useEffect } from 'react';
import { useDispatch } from 'react-redux'; // Import useDispatch
import { getProductCatalog, getCategories } from '../api/api';
import ProductCard from '../components/ProductCard';
import './ProductCatalog.css';
import { FaFilter, FaSort, FaHeart } from 'react-icons/fa'; // Import icons
import { addAlert } from '../redux/alertSlice'; // Import addAlert action
import { v4 as uuidv4 } from 'uuid'; // Import uuid for unique alert IDs

function ProductCatalog() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showFilter, setShowFilter] = useState(false); // State to toggle filter form visibility
  const dispatch = useDispatch(); // Initialize useDispatch

  useEffect(() => {
    fetchProducts();
    fetchCategories();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    setError(null);
    try {
      const params = {
        ...(selectedCategory && { category_id: selectedCategory }),
        ...(minPrice && { min_price: parseFloat(minPrice) }),
        ...(maxPrice && { max_price: parseFloat(maxPrice) }),
      };
      const response = await getProductCatalog(params);
      if (response.data.products && response.data.products.length > 0) {
        setProducts(response.data.products);
      } else {
        setProducts([]);
        dispatch(addAlert({ id: uuidv4(), message: 'Пусто: Продукты не найдены.', type: 'info' }));
      }
    } catch (err) {
      const errorMessage = err.response?.data?.error || err.message || 'Неизвестная ошибка при загрузке продуктов.';
      setError(errorMessage);
      dispatch(addAlert({ id: uuidv4(), message: `Ошибка: ${errorMessage}`, type: 'error' }));
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await getCategories();
      if (response.data.categories && response.data.categories.length > 0) {
        setCategories(response.data.categories);
      } else {
        setCategories([]);
        dispatch(addAlert({ id: uuidv4(), message: 'Пусто: Категории не найдены.', type: 'info' }));
      }
    } catch (err) {
      console.error("Error fetching categories:", err);
      const errorMessage = err.response?.data?.error || err.message || 'Неизвестная ошибка при загрузке категорий.';
      dispatch(addAlert({ id: uuidv4(), message: `Ошибка: ${errorMessage}`, type: 'error' }));
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchProducts();
  };

  const handleCategoryTabClick = (categoryId) => {
    setSelectedCategory(categoryId);
    fetchProducts(); // Refetch products for the selected category
  };

  if (loading) return <div className="loading-message">Загрузка товаров...</div>;
  if (error) return <div className="error-message">Ошибка: {error}</div>;

  return (
    <div className="product-catalog-page">
      <h1 className="page-title">КАТАЛОГ</h1>

      <div className="category-tabs">
        <button 
          className={`category-tab ${selectedCategory === '' ? 'active' : ''}`}
          onClick={() => handleCategoryTabClick('')}
        >
          Все категории
        </button>
        {categories.map((cat) => (
          <button 
            key={cat.id} 
            className={`category-tab ${selectedCategory === cat.id ? 'active' : ''}`}
            onClick={() => handleCategoryTabClick(cat.id)}
          >
            {cat.name}
          </button>
        ))}
      </div>

      <div className="catalog-actions">
        <button className="filter-button" onClick={() => setShowFilter(!showFilter)} style={{backgroundColor: 'white', color: 'black'}}>
          <FaFilter /> Фильтр
        </button>
        <button className="sort-button">
          <FaSort /> Сортировка
        </button>
      </div>

      {showFilter && (
        <form onSubmit={handleSearch} className="filter-form-expanded">
          <div className="form-group">
            <label htmlFor="category-select">Категория:</label>
            <select id="category-select" value={selectedCategory} onChange={(e) => setSelectedCategory(e.target.value)} className="filter-select">
              <option value="">Все категории</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label htmlFor="min-price">Мин. Цена:</label>
            <input id="min-price" type="number" value={minPrice} onChange={(e) => setMinPrice(e.target.value)} className="filter-input" placeholder="0" />
          </div>
          <div className="form-group">
            <label htmlFor="max-price">Макс. Цена:</label>
            <input id="max-price" type="number" value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)} className="filter-input" placeholder="1000" />
          </div>
          <button type="submit" className="apply-filter-button">Найти</button>
        </form>
      )}

      {products.length === 0 && !loading && !error && (
        <div className="no-products-message">Продукты не найдены.</div>
      )}

      <div className="products-grid">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} categories={categories} />
        ))}
      </div>
    </div>
  );
}

export default ProductCatalog;
