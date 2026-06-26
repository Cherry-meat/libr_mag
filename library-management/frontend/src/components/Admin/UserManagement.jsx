import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'

function UserManagement() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await api.get('/admin/users')
      setUsers(response.data)
    } catch (error) {
      toast.error('Failed to fetch users')
    } finally {
      setLoading(false)
    }
  }

  const updateRole = async (userId, newRole) => {
    try {
      await api.put(`/admin/users/${userId}/role`, { role: newRole })
      setUsers(users.map(u => 
        u.id === userId ? { ...u, role: newRole } : u
      ))
      toast.success('User role updated')
    } catch (error) {
      toast.error('Failed to update role')
    }
  }

  const toggleBlock = async (userId, isBlocked) => {
    try {
      await api.post(`/admin/users/${userId}/block`, { 
        blocked: !isBlocked,
        reason: isBlocked ? 'Unblocked by admin' : 'Blocked by admin'
      })
      setUsers(users.map(u => 
        u.id === userId ? { ...u, is_blocked: !isBlocked } : u
      ))
      toast.success(`User ${isBlocked ? 'unblocked' : 'blocked'}`)
    } catch (error) {
      toast.error('Failed to update user status')
    }
  }

  if (loading) {
    return <div className="loading">Loading users...</div>
  }

  return (
    <div className="user-management">
      <h2>User Management</h2>
      <table className="users-table">
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Ticket</th>
            <th>Status</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map(user => (
            <tr key={user.id}>
              <td>{user.first_name} {user.last_name}</td>
              <td>{user.email}</td>
              <td>
                <select 
                  value={user.role}
                  onChange={(e) => updateRole(user.id, e.target.value)}
                >
                  <option value="user">User</option>
                  <option value="admin">Admin</option>
                </select>
              </td>
              <td>{user.ticket_linked ? '✅' : '❌'}</td>
              <td>
                <span className={user.is_blocked ? 'blocked' : 'active'}>
                  {user.is_blocked ? 'Blocked' : 'Active'}
                </span>
              </td>
              <td>
                <button 
                  className={user.is_blocked ? 'unblock-btn' : 'block-btn'}
                  onClick={() => toggleBlock(user.id, user.is_blocked)}
                >
                  {user.is_blocked ? 'Unblock' : 'Block'}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default UserManagement