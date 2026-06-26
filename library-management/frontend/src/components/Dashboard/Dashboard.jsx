import React from 'react'
import { useAuth } from '../../contexts/AuthContext'
import { Link } from 'react-router-dom'
import { FaBook, FaUser, FaBell, FaTicketAlt } from 'react-icons/fa'

function Dashboard() {
  const { user } = useAuth()

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>Welcome, {user?.first_name || 'User'}!</h1>
        <div className="header-actions">
          <Link to="/notifications" className="icon-btn">
            <FaBell />
          </Link>
          <Link to="/profile" className="icon-btn">
            <FaUser />
          </Link>
        </div>
      </header>

      <div className="stats-grid">
        <div className="stat-card">
          <FaBook className="stat-icon" />
          <div className="stat-info">
            <h3>Active Loans</h3>
            <p>0</p>
          </div>
        </div>
        <div className="stat-card">
          <FaTicketAlt className="stat-icon" />
          <div className="stat-info">
            <h3>Ticket Status</h3>
            <p>{user?.ticket_linked ? 'Linked' : 'Not Linked'}</p>
          </div>
        </div>
      </div>

      <div className="dashboard-grid">
        <div className="quick-actions">
          <h2>Quick Actions</h2>
          <div className="action-grid">
            <Link to="/books" className="action-btn">View Book History</Link>
            <Link to="/ticket" className="action-btn">View Ticket</Link>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Dashboard