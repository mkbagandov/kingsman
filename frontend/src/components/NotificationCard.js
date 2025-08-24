import React from 'react';
import { FaInfoCircle, FaGift, FaBell } from 'react-icons/fa';
import './NotificationCard.css'; // Assuming you'll create this CSS file

function NotificationCard({ notification }) {
  const getIcon = (type) => {
    switch (type) {
      case 'promotion':
        return <FaGift className="icon promotion" />;
      case 'new_arrival':
        return <FaBell className="icon new-arrival" />;
      default:
        return <FaInfoCircle className="icon info" />;
    }
  };

  return (
    <div className="notification-card">
      <div className="notification-icon-container">
        {getIcon(notification.type)}
      </div>
      <div className="notification-content">
        <h3 className="notification-title">{notification.title}</h3>
        <p className="notification-message">{notification.message}</p>
        <span className="notification-timestamp">{new Date(notification.created_at).toLocaleString()}</span>
      </div>
    </div>
  );
}

export default NotificationCard;
