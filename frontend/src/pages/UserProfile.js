import React, { useState, useEffect } from 'react';
// Removed useParams as userID will come from JWT on the backend
import { getUserProfile, getUserLoyaltyProfile, getUserDiscountCard, getUserQRCode } from '../api/api';
import { FaUserCircle, FaStar, FaCreditCard, FaQrcode, FaEnvelope, FaPhone, FaAward, FaHistory, FaGlobe, FaGithub, FaTwitter, FaInstagram, FaFacebook, FaTasks } from 'react-icons/fa';
import { Link } from 'react-router-dom'; // Added Link import

function UserProfile() {
  // const { userID } = useParams(); // No longer needed
  const [userProfile, setUserProfile] = useState(null);
  const [loyaltyProfile, setLoyaltyProfile] = useState(null);
  const [discountCard, setDiscountCard] = useState(null);
  const [qrCode, setQrCode] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Calls to API functions no longer need userID parameter
        const [userRes, loyaltyRes, discountRes, qrCodeRes] = await Promise.all([
          getUserProfile(),
          getUserLoyaltyProfile(),
          getUserDiscountCard(),
          getUserQRCode()
        ]);
        setUserProfile(userRes.data);
        setLoyaltyProfile(loyaltyRes.data);
        setDiscountCard(discountRes.data);

        // Convert arraybuffer to base64 for image display
        const base64Image = btoa(
          new Uint8Array(qrCodeRes.data).reduce(
            (data, byte) => data + String.fromCharCode(byte),
            ''
          )
        );
        setQrCode(`data:image/png;base64,${base64Image}`);

      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []); // Empty dependency array as userID is handled by backend JWT

  if (loading) return <div className="user-profile-page">Загрузка профиля пользователя...</div>;
  if (error) return <div className="user-profile-page">Ошибка: {error}</div>;
  if (!userProfile) return <div className="user-profile-page">Профиль пользователя не найден.</div>;

  return (
    <div className="user-profile-container">
      <div className="profile-header-breadcrumbs">
        <Link to="/">Главная</Link> / Пользователь / Профиль пользователя
      </div>

      <div className="profile-main-content">
        <div className="profile-sidebar">
          {/* Profile Card */}
          <div className="profile-card profile-info-card">
            <FaUserCircle className="profile-avatar-icon" /> {/* Replaced img with icon */}
            <h2>{userProfile.username}</h2>
            <div className="profile-actions">
              <button className="btn-follow">Редактировать профиль</button>
              <button className="btn-message">Выйти</button>
            </div>
          </div>

          {/* Social Links Card */}
          <div className="profile-card profile-social-card">
            <div className="social-link-item">
              <FaGlobe />
              <span>Веб-сайт</span>
              <a href="https://bootdey.com" target="_blank" rel="noopener noreferrer">bootdey.com</a>
            </div>
            <div className="social-link-item">
              <FaGithub />
              <span>GitHub</span>
              <span>bootdey</span>
            </div>
            <div className="social-link-item">
              <FaTwitter />
              <span>Twitter</span>
              <span>@bootdey</span>
            </div>
            <div className="social-link-item">
              <FaInstagram />
              <span>Instagram</span>
              <span>bootdey</span>
            </div>
            <div className="social-link-item">
              <FaFacebook />
              <span>Facebook</span>
              <span>bootdey</span>
            </div>
          </div>
        </div>

        <div className="profile-content-area">
          {/* Personal Information Card */}
          <div className="profile-card profile-details-card">
            <h2 className="card-title">Полное имя: {userProfile.username}</h2>
            <div className="detail-item">
              <span>Email</span>
              <span>{userProfile.email}</span>
            </div>
            <div className="detail-item">
              <span>Телефон</span>
              <span>{userProfile.phone_number}</span>
            </div>
            <div className="detail-item">
              <span>Адрес</span>
              <span>Москва, Россия</span> {/* Placeholder for address */}
            </div>
          </div>

          {/* Loyalty and Discount Card, QR Code - integrate as separate cards or within existing ones */}
          {/* Current Loyalty Profile (Moved or adapted) */}
          {loyaltyProfile && (
            <div className="profile-card profile-loyalty-card">
              <h3 className="card-title">Информация о лояльности <FaStar /></h3>
              <p>Текущие баллы: {loyaltyProfile.current_points}</p>
              <p>Статус лояльности: {loyaltyProfile.loyalty_status}</p>
              {loyaltyProfile.current_tier && (
                <div className="loyalty-tier-details">
                  <h4><FaAward /> Текущий уровень: {loyaltyProfile.current_tier.name}</h4>
                  <p>Мин. баллов для этого уровня: {loyaltyProfile.current_tier.min_points}</p>
                  <p>Описание: {loyaltyProfile.current_tier.description}</p>
                  <p>Преимущества: {loyaltyProfile.current_tier.benefits}</p>
                </div>
              )}
              {loyaltyProfile.loyalty_activities && loyaltyProfile.loyalty_activities.length > 0 && (
                <div className="loyalty-activities">
                  <h4><FaHistory /> Последние действия лояльности:</h4>
                  <ul>
                    {loyaltyProfile.loyalty_activities.map((activity) => (
                      <li key={activity.id}>
                        <strong>{activity.type}:</strong> {activity.description} ({new Date(activity.created_at).toLocaleString()})
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}

          <div className="discount-qr-container">
            {/* Discount Card (Moved or adapted) */}
            {discountCard && (
              <div className="profile-card profile-discount-card">
                <h3 className="card-title">Дисконтная карта <FaCreditCard /></h3>
                <div className="discount-card-mock">
                  <h3>MR.KINGSMAN</h3>
                  <p>ДИСКОНТНАЯ КАРТА</p>
                  <p>Уровень: {discountCard.discount_level}</p>
                  <p>Прогресс: {discountCard.progress_to_next_level}% до следующего уровня</p>
                </div>
              </div>
            )}

            {/* QR Code (Moved or adapted) */}
            {qrCode && (
              <div className="profile-qr-code-card">
                <img src={qrCode} alt="QR-код" className="profile-qr-code" />
              </div>
            )}
          </div>

        </div>
      </div>
    </div>
  );
}

export default UserProfile;
