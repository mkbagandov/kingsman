import React, { useState, useEffect } from 'react';
import { getNotifications } from '../api/api';
import NotificationCard from '../components/NotificationCard'; // Import the new NotificationCard component

function Notifications() {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await getNotifications();
        setNotifications(response.data.notifications);
      } catch (err) {
        setError(err.response?.data?.error || err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchNotifications();
  }, []);

  if (loading) return <div>Loading notifications...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="notifications-page">
      <h1>Your Notifications</h1>
      {notifications.length > 0 ? (
        <div className="notifications-list">
          {notifications.map((notification) => (
            <NotificationCard key={notification.id} notification={notification} />
          ))}
        </div>
      ) : (
        <p>No notifications found.</p>
      )}
    </div>
  );
}

export default Notifications;
