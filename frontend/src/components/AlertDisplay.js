import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { removeAlert } from '../redux/alertSlice';
import './AlertDisplay.css';

const AlertDisplay = () => {
  const dispatch = useDispatch();
  const alerts = useSelector((state) => state.alert.alerts);

  useEffect(() => {
    if (alerts.length > 0) {
      const timer = setTimeout(() => {
        dispatch(removeAlert(alerts[0].id));
      }, 3000); // Alerts disappear after 3 seconds
      return () => clearTimeout(timer);
    }
  }, [alerts, dispatch]);

  return (
    <div className="alert-container">
      {alerts.map((alert) => (
        <div key={alert.id} className={`alert-item ${alert.type || 'info'}`}>
          <span>{alert.message}</span>
          <button onClick={() => dispatch(removeAlert(alert.id))}>&times;</button>
        </div>
      ))}
    </div>
  );
};

export default AlertDisplay;
