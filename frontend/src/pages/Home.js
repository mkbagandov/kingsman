import React from 'react';

function Home() {
  return (
    <div className="home-page">
      {/* Hero Section */}
      <section className="hero-section">
        <div className="hero-content">
          <div className="hero-text">
            <h1>MR. KINGSMAN</h1>
            <h2>Премиальная мужская одежда и стильный образ жизни</h2>
            <div className="hero-buttons">
              <button className="btn-primary">Перейти в каталог</button>
              <button className="btn-secondary">Узнать больше</button>
            </div>
          </div>
          <div className="hero-icon-container">
            <div className="hero-icon">
              👔
            </div>
          </div>
        </div>
      </section>

      {/* Why MR. KINGSMAN? Section */}
      <section className="why-mrg-section">
        <h2>Почему MR. KINGSMAN?</h2>
        <p>Мы создаем не просто одежду, а стиль жизни для современного мужчины</p>
        <div className="why-mrg-cards">
          <div className="info-card">
            <div className="icon">✨</div>
            <h3>Качество</h3>
            <p>Безупречное исполнение и долговечность каждого изделия.</p>
          </div>
          <div className="info-card">
            <div className="icon">🎩</div>
            <h3>Стиль</h3>
            <p>Элегантность и соответствие последним тенденциям моды.</p>
          </div>
          <div className="info-card">
            <div className="icon">💎</div>
            <h3>Индивидуальность</h3>
            <p>Подчеркните свой уникальный образ с нашей коллекцией.</p>
          </div>
        </div>
      </section>
    </div>
  );
}

export default Home;
