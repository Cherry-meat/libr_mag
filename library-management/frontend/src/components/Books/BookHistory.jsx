import React, { useState, useEffect } from 'react'
import api from '../../services/api'
import { toast } from 'react-toastify'
import { FaBook } from 'react-icons/fa'

function BookHistory() {
  const [loans, setLoans] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchHistory()
  }, [])

  const fetchHistory = async () => {
    try {
      const response = await api.get('/books/history')
      setLoans(response.data)
    } catch (error) {
      toast.error('Failed to fetch book history')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  return (
    <div className="book-history">
      <h2>Book History</h2>
      {loans.length === 0 ? (
        <div className="empty-state">
          <FaBook size={50} />
          <p>No books borrowed yet</p>
        </div>
      ) : (
        <div className="loans-list">
          {loans.map(loan => (
            <div key={loan.id} className="loan-item">
              <div className="book-details">
                <h3>{loan.title}</h3>
                <p>by {loan.author}</p>
                <div className="loan-dates">
                  <span>Borrowed: {new Date(loan.loan_date).toLocaleDateString()}</span>
                  <span>Due: {new Date(loan.due_date).toLocaleDateString()}</span>
                </div>
              </div>
              <div className="loan-status">
                <span className={`status ${loan.is_returned ? 'returned' : 'active'}`}>
                  {loan.is_returned ? 'Returned' : 'Active'}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default BookHistory