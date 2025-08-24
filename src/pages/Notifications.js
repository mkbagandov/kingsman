import React, { useState, useEffect } from 'react';
import { getNotifications } from '../api/api';

function Notifications() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  // In a real app, you would get the userID from authentication context
  const userID = "1"; // Placeholder userID

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await getNotifications(userID);
        setNotifications(response.data);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchNotifications();
  }, [userID]);

  if (loading) return <div>Loading notifications...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Your Notifications</h1>
      {notifications.length > 0 ? (
        <ul>
          {notifications.map((notification) => (
            <li key={notification.id}>
              <strong>{notification.type}:</strong> {notification.message} - {new Date(notification.timestamp).toLocaleString()}
            </li>
          ))}
        </ul>
      ) : (
        <p>No notifications found.</p>
      )}
    </div>
  );
}

export default Notifications;
