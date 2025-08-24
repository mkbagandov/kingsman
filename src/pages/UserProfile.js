import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getUserProfile, getUserLoyaltyProfile, getUserDiscountCard, getUserQRCode } from '../api/api';

function UserProfile() {
  const { userID } = useParams();
  const [userProfile, setUserProfile] = useState(null);
  const [loyaltyProfile, setLoyaltyProfile] = useState(null);
  const [discountCard, setDiscountCard] = useState(null);
  const [qrCode, setQrCode] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [userRes, loyaltyRes, discountRes, qrCodeRes] = await Promise.all([
          getUserProfile(userID),
          getUserLoyaltyProfile(userID),
          getUserDiscountCard(userID),
          getUserQRCode(userID)
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
  }, [userID]);

  if (loading) return <div>Loading user profile...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!userProfile) return <div>User profile not found.</div>;

  return (
    <div>
      <h1>User Profile: {userProfile.username}</h1>
      <p>Email: {userProfile.email}</p>

      <h2>Loyalty Information</h2>
      {loyaltyProfile ? (
        <div>
          <p>Points: {loyaltyProfile.points}</p>
          <p>Tier: {loyaltyProfile.tier}</p>
          {/* Display other loyalty details */}
        </div>
      ) : (
        <p>No loyalty information available.</p>
      )}

      <h2>Discount Card</h2>
      {discountCard ? (
        <div>
          <p>Discount Level: {discountCard.discount_level}</p>
          <p>Progress to Next Level: {discountCard.progress_to_next_level}%</p>
          {/* Display other discount card details */}
        </div>
      ) : (
        <p>No discount card information available.</p>
      )}

      <h2>QR Code</h2>
      {qrCode && <img src={qrCode} alt="QR Code" />}
    </div>
  );
}

export default UserProfile;
