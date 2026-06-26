import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'
import { FaBell } from 'react-icons/fa'

function Notifications() {
  const [notifications, setNotifications] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchNotifications()
  }, [])

  const fetchNotifications = async () => {
    try {
      const response = await api.get('/notifications')
      setNotifications(response.data)
    } catch (error) {
      toast.error('Failed to fetch notifications')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading notifications...</div>
  }

  return (
    <div className="notifications-container">
      <h2>Notifications</h2>
      {notifications.length === 0 ? (
        <div className="empty-state">
          <FaBell size={50} />
          <p>No notifications</p>
        </div>
      ) : (
        <div className="notifications-list">
          {notifications.map(notification => (
            <div key={notification.id} className="notification-item">
              <div className="notification-content">
                <h4>{notification.title}</h4>
                <p>{notification.message}</p>
                <span className="notification-time">
                  {new Date(notification.created_at).toLocaleString()}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default Notifications