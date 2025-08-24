import React, { useState, useEffect } from 'react';
import { getProductCatalog, getCategories } from '../api/api';
import ProductCard from '../components/ProductCard'; // Import the new ProductCard component

function ProductCatalog() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState('');
  const [minPrice, setMinPrice] = useState('');
  const [maxPrice, setMaxPrice] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

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
      setProducts(response.data.products);
    } catch (err) {
      setError(err.response?.data?.error || err.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await getCategories();
      setCategories(response.data.categories); // Adjusting to the expected backend response
    } catch (err) {
      console.error("Error fetching categories:", err);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchProducts();
  };

  if (loading) return <div>Загрузка товаров...</div>;
  if (error) return <div>Ошибка: {error}</div>;

  return (
    <div className="product-catalog-page">
      <h1>Каталог Продуктов</h1>
      <form onSubmit={handleSearch} className="product-filter-form">
        <div className="form-group">
          <label htmlFor="category-select">Категория:</label>
          <select id="category-select" value={selectedCategory} onChange={(e) => setSelectedCategory(e.target.value)}>
            <option value="">Все Категории</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label htmlFor="min-price">Мин. Цена:</label>
          <input id="min-price" type="number" value={minPrice} onChange={(e) => setMinPrice(e.target.value)} />
        </div>
        <div className="form-group">
          <label htmlFor="max-price">Макс. Цена:</label>
          <input id="max-price" type="number" value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)} />
        </div>
        <button type="submit" className="filter-button">Фильтровать Продукты</button>
      </form>

      <div className="products-grid">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} categories={categories} />
        ))}
      </div>
    </div>
  );
}

export default ProductCatalog;
