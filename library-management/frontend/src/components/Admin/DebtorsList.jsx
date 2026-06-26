import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'

function DebtorsList() {
  const [debtors, setDebtors] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchDebtors()
  }, [])

  const fetchDebtors = async () => {
    try {
      const response = await api.get('/admin/debtors')
      setDebtors(response.data)
    } catch (error) {
      toast.error('Failed to fetch debtors')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading debtors...</div>
  }

  return (
    <div className="debtors-list">
      <h2>Debtors List</h2>
      {debtors.length === 0 ? (
        <p>No debtors found. Great job!</p>
      ) : (
        <table className="debtors-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Overdue Books</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {debtors.map(debtor => (
              <tr key={debtor.id}>
                <td>{debtor.first_name} {debtor.last_name}</td>
                <td>{debtor.email}</td>
                <td className="overdue-count">{debtor.overdue_books}</td>
                <td>
                  <button className="notify-btn">Send Notification</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default DebtorsList