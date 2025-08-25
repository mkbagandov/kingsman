import React, { useState, useEffect } from 'react';
import './EditProfileModal.css';

const EditProfileModal = ({ user, isOpen, onClose, onSave }) => {
  const [editedUser, setEditedUser] = useState(user);

  useEffect(() => {
    setEditedUser(user);
  }, [user]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setEditedUser((prevUser) => ({
      ...prevUser,
      [name]: value,
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSave(editedUser);
  };

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal-content">
        <div className="modal-header">
          <h2>Редактировать профиль</h2>
          <button className="modal-close-button" onClick={onClose}>&times;</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">Полное имя:</label>
            <input
              type="text"
              id="username"
              name="username"
              value={editedUser.username || ''}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="email">Email:</label>
            <input
              type="email"
              id="email"
              name="email"
              value={editedUser.email || ''}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="phone_number">Телефон:</label>
            <input
              type="text"
              id="phone_number"
              name="phone_number"
              value={editedUser.phone_number || ''}
              onChange={handleChange}
            />
          </div>
          {/* Add other editable fields as needed */}
          <div className="modal-actions">
            <button type="submit" className="btn-save">Сохранить изменения</button>
            <button type="button" className="btn-cancel" onClick={onClose}>Отмена</button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default EditProfileModal;
