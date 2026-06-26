import React from 'react'

function AdminDashboard() {
  return (
    <div className="admin-dashboard">
      <h1>Admin Dashboard</h1>
      <div className="admin-stats">
        <div className="stat-card">
          <h3>Total Users</h3>
          <p>0</p>
        </div>
        <div className="stat-card">
          <h3>Total Books</h3>
          <p>0</p>
        </div>
      </div>
      <div className="admin-content">
        <p>Admin panel coming soon...</p>
      </div>
    </div>
  )
}

export default AdminDashboard