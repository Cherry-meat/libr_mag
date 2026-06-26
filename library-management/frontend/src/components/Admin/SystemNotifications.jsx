import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'

function SystemNotifications() {
  const [notifications, setNotifications] = useState([])
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    message: '',
    priority: 'normal'
  })

  useEffect(() => {
    fetchNotifications()
  }, [])

  const fetchNotifications = async () => {
    try {
      const response = await api.get('/admin/system-notifications')
      setNotifications(response.data)
    } catch (error) {
      toast.error('Failed to fetch system notifications')
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await api.post('/admin/system-notifications', formData)
      toast.success('System notification created')
      setFormData({ title: '', message: '', priority: 'normal' })
      setShowForm(false)
      fetchNotifications()
    } catch (error) {
      toast.error('Failed to create notification')
    }
  }

  const getPriorityColor = (priority) => {
    switch(priority) {
      case 'high': return '#ff4444'
      case 'medium': return '#ffa500'
      default: return '#4caf50'
    }
  }

  return (
    <div className="system-notifications">
      <div className="notifications-header">
        <h2>System Notifications</h2>
        <button onClick={() => setShowForm(true)} className="create-btn">
          New Notification
        </button>
      </div>

      {showForm && (
        <div className="modal-overlay">
          <div className="modal">
            <h3>Create System Notification</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Title</label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({...formData, title: e.target.value})}
                  required
                />
              </div>
              <div className="form-group">
                <label>Message</label>
                <textarea
                  value={formData.message}
                  onChange={(e) => setFormData({...formData, message: e.target.value})}
                  rows="4"
                  required
                />
              </div>
              <div className="form-group">
                <label>Priority</label>
                <select
                  value={formData.priority}
                  onChange={(e) => setFormData({...formData, priority: e.target.value})}
                >
                  <option value="normal">Normal</option>
                  <option value="medium">Medium</option>
                  <option value="high">High</option>
                </select>
              </div>
              <div className="modal-actions">
                <button type="submit">Create</button>
                <button type="button" onClick={() => setShowForm(false)}>
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <div className="notifications-list">
        {notifications.map(notification => (
          <div key={notification.id} className="notification-item">
            <div 
              className="priority-indicator"
              style={{ backgroundColor: getPriorityColor(notification.priority) }}
            />
            <div className="content">
              <h4>{notification.title}</h4>
              <p>{notification.message}</p>
              <span className="timestamp">
                {new Date(notification.created_at).toLocaleString()}
              </span>
            </div>
            <div className="priority-badge">
              {notification.priority.toUpperCase()}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default SystemNotifications