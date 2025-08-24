import React, { useState, useEffect } from 'react';
// Removed useParams as userID will come from JWT on the backend
import { getUserProfile, getUserLoyaltyProfile, getUserDiscountCard, getUserQRCode } from '../api/api';
import { FaUserCircle, FaStar, FaCreditCard, FaQrcode, FaEnvelope, FaPhone, FaAward, FaHistory } from 'react-icons/fa';

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

  if (loading) return <div className="user-profile-page">Loading user profile...</div>;
  if (error) return <div className="user-profile-page">Error: {error}</div>;
  if (!userProfile) return <div className="user-profile-page">User profile not found.</div>;

  return (
    <div className="user-profile-page">
      <h1><FaUserCircle /> Профиль пользователя: {userProfile.username}</h1>

      <div className="profile-section">
        <h2>Персональная информация</h2>
        <p><FaEnvelope /> Email: {userProfile.email}</p>
        <p><FaPhone /> Телефон: {userProfile.phone_number}</p>
      </div>

      <div className="profile-section">
        <h2>Информация о лояльности <FaStar /></h2>
        {loyaltyProfile ? (
          <div>
            <p>Текущие баллы: {loyaltyProfile.current_points}</p>
            <p>Статус лояльности: {loyaltyProfile.loyalty_status}</p>
            {loyaltyProfile.current_tier && (
              <div className="loyalty-tier-details">
                <h3><FaAward /> Текущий уровень: {loyaltyProfile.current_tier.name}</h3>
                <p>Мин. баллов для этого уровня: {loyaltyProfile.current_tier.min_points}</p>
                <p>Описание: {loyaltyProfile.current_tier.description}</p>
                <p>Преимущества: {loyaltyProfile.current_tier.benefits}</p>
              </div>
            )}

            {loyaltyProfile.loyalty_activities && loyaltyProfile.loyalty_activities.length > 0 && (
              <div className="loyalty-activities">
                <h3><FaHistory /> Последние действия лояльности:</h3>
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
        ) : (
          <p>Информация о лояльности недоступна.</p>
        )}
      </div>

      <div className="profile-section">
        <h2>Дисконтная карта <FaCreditCard /></h2>
        {discountCard ? (
          <div className="discount-card-mock">
            <h3>MR.KINGSMAN</h3>
            <p>ДИСКОНТНАЯ КАРТА</p>
            <p>Уровень: {discountCard.discount_level}</p>
            <p>Прогресс: {discountCard.progress_to_next_level}% до следующего уровня</p>
            {/* Add more discount card details as needed */}
          </div>
        ) : (
          <p>Информация о дисконтной карте недоступна.</p>
        )}
      </div>

      {qrCode && (
        <div className="qr-code-container profile-section">
          <h2>Ваш QR-код <FaQrcode /></h2>
          <img src={qrCode} alt="QR Code" />
        </div>
      )}
    </div>
  );
}

export default UserProfile;
