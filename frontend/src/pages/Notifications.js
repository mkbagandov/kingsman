import React, { useState, useEffect } from 'react';
import { getNotifications } from '../api/api';
// Removed useParams as userID will come from JWT on the backend

function Notifications() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  // const { userID } = useParams(); // No longer needed for notifications

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await getNotifications(); // No userID parameter
        setNotifications(response.data.notifications);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchNotifications();
  }, []); // Empty dependency array as userID is handled by backend JWT

  if (loading) return <div>Loading notifications...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Your Notifications</h1>
      {notifications.length > 0 ? (
        <ul>
          {notifications.map((notification) => (
            <li key={notification.id}>
              <strong>{notification.type}:</strong> {notification.message} - {new Date(notification.createdAt).toLocaleString()}
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
